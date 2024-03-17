package routers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strconv"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/oss"
	xt "xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"

	"github.com/gin-gonic/gin"
	"time"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/config"
	"xhyovo.cn/community/pkg/result"
)

var expire_time int64 = 30

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
	group.Use(middleware.OperLogger())
	group.GET("/policy", getPolicy)
	group.GET("/singUrl", getUrl)
	group.POST("/upload", uploadCallback)
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

	body := fmt.Sprintf("{\"fileKey\":${object},\"size\":${size},\"mimeType\":${mimeType},\"x:userId\":%d}", userId)

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
		log.Warnf("用户id: %d 获取 %s 失败,因为 %s 为空", middleware.GetUserId(ctx), fileKey, fileKey)
		result.Err("fileKey 为空").Json(ctx)
		return
	}
	singUrl := oss.SingUrl(fileKey)
	ctx.Redirect(http.StatusFound, singUrl)
}

func uploadCallback(ctx *gin.Context) {

	file := &model.Files{}

	if err := ctx.ShouldBindJSON(&file); err != nil {
		log.Warnf("用户id: %d 上传文件 callback 解析参数失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
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

	file.UserId = middleware.GetUserId(ctx)
	file.CreatedAt = xt.Now()
	file.UpdatedAt = xt.Now()
	var fileS services.FileService
	fileS.Save(file)
	result.Ok(nil, "").Json(ctx)
}
