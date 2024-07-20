package frontend

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"hash"
	"io"
	"net/http"
	"strconv"
	"strings"
	"xhyovo.cn/community/pkg/cache"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/oss"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/request"
	services "xhyovo.cn/community/server/service"

	"github.com/gin-gonic/gin"
	"time"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/config"
	"xhyovo.cn/community/pkg/result"
)

var expire_time int64 = 5000

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}

type CallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

func InitFileRouters(ctx *gin.Engine) {
	group := ctx.Group("/community/file")
	group.POST("/upload", uploadCallback)
	group.Use(middleware.Auth)
	group.GET("", listFiles)
	group.GET("/byKey", getFileByKey)
	group.Use(middleware.OperLogger())
	group.GET("/policy", getPolicy)
	group.GET("/singUrl", getUrl)
}

func get_gmt_iso8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

func getPolicy(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Methods", "POST")
	ctx.Header("Access-Control-Allow-Origin", "*")
	now := time.Now().Unix()
	expire_end := now + expire_time
	var tokenExpire = get_gmt_iso8601(expire_end)

	var userId = middleware.GetUserId(ctx)

	var prefix = strconv.Itoa(userId) + "/"

	//create post policy json
	var cf ConfigStruct
	cf.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, prefix)
	cf.Conditions = append(cf.Conditions, condition)

	//calucate signature
	r, err := json.Marshal(cf)
	if err != nil {
		log.Warnln(err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	debyte := base64.StdEncoding.EncodeToString(r)
	ossConfig := config.GetInstance().OssConfig
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(ossConfig.SecretKey))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	uId := uuid.NewString()
	cache.GetInstance().Set(uId, 1, 60*time.Second)
	body := fmt.Sprintf("{\"fileKey\":${object},\"size\":${size},\"mimeType\":${mimeType},\"x:userId\":%d,\"x:uuid\":\"%s\"}", userId, uId)

	var callbackParam CallbackParam
	callbackParam.CallbackUrl = ossConfig.Callback
	callbackParam.CallbackBody = body
	callbackParam.CallbackBodyType = "application/json"
	callback_str, err := json.Marshal(callbackParam)
	if err != nil {
		log.Warnln(err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callback_str)

	var policyToken PolicyToken
	policyToken.AccessKeyId = ossConfig.AccessKey
	policyToken.Host = ossConfig.Endpoint
	policyToken.Expire = expire_end
	policyToken.Signature = signedStr
	policyToken.Directory = prefix
	policyToken.Policy = debyte
	policyToken.Callback = callbackBase64
	if err != nil {
		log.Warnln(err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}

	result.Ok(policyToken, "成功").Json(ctx)
}

func getUrl(ctx *gin.Context) {
	fileKey := ctx.Query("fileKey")
	if fileKey == "" {
		return
	}
	singUrl := oss.SingUrl(fileKey)

	replace := strings.Replace(singUrl, fmt.Sprintf("%s.%s", oss.GetInstance().BucketName, oss.GetEndpoint()),
		config.GetInstance().OssConfig.Cdn, 1)

	ctx.Redirect(http.StatusFound, replace)
}

func uploadCallback(ctx *gin.Context) {

	callback := &request.OssCallback{}

	if err := ctx.ShouldBindJSON(&callback); err != nil {
		log.Warnf("上传文件 callback 解析参数失败,err: %s", err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	v, b := cache.GetInstance().Get(callback.Uuid)
	if !b || v == nil {
		log.Warnf("用户id: %d 上传文件 callback 解析uuid失败")
		result.Err("文件上传 callback 失败").Json(ctx)
		return
	}

	file := &model.Files{
		FileKey: callback.FileKey,
		Size:    callback.Size,
		Format:  callback.MimeType,
		UserId:  callback.UserId,
	}

	// check fileKey not empty
	bucket := oss.GetInstance()
	exist, err := bucket.IsObjectExist(file.FileKey)
	if err != nil {
		log.Warnf("用户id: %d 判断文件为空失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	if !exist {
		log.Warnf("用户id: %d 判断文件为空", middleware.GetUserId(ctx))
		result.Err("文件不存在").Json(ctx)
		return
	}

	var fileS services.FileService
	fileS.Save(file)
	result.Ok(nil, "").Json(ctx)
}

// 查询用户的资源
func listFiles(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	userId := middleware.GetUserId(ctx)
	var fileS services.FileService
	files, count := fileS.PageFiles(p, limit, userId)
	result.Page(files, count, nil).Json(ctx)
}

func getFileByKey(ctx *gin.Context) {
	fileKey := ctx.Query("fileKey")

	// 先从 db 拿
	var fileS services.FileService
	if !fileS.ExistFile(fileKey) {
		// 从 oss 拿
		bucket := oss.GetInstance()
		exist, err := bucket.IsObjectExist(fileKey)
		if err != nil {
			result.Err("查询资源出错:" + err.Error()).Json(ctx)
			return
		}
		// 如果存在则放入 db
		if exist {
			meta, err := bucket.GetObjectDetailedMeta(fileKey)
			if err == nil {
				size := meta.Get("Content-Length")
				atoi, _ := strconv.Atoi(size)
				file := &model.Files{
					FileKey: fileKey,
					Size:    int64(atoi),
					Format:  meta.Get("Content-Type"),
					UserId:  middleware.GetUserId(ctx),
				}
				fileS.Save(file)
				result.OkWithMsg(true, "上传资源成功,如未能显示,则从资源库中复制获取").Json(ctx)
				return
			}
		}
		return
	}
	result.OkWithMsg(true, "上传资源成功,如未能显示,则从资源库中复制获取").Json(ctx)
	return
}
