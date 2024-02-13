package routers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"xhyovo.cn/community/pkg/oss"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"

	"hash"
	"io"

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

	group.GET("/policy", getPolicy)
	group.GET("/singUrl/:fileKey", getUrl)
	group.POST("/upload", upload)
}

func get_gmt_iso8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

func getPolicy(ctx *gin.Context) {
	now := time.Now().Unix()
	expire_end := now + expire_time
	var tokenExpire = get_gmt_iso8601(expire_end)

	var userId int = middleware.GetUserId(ctx)

	var prefix string = strconv.Itoa(userId) + "/"

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
	debyte := base64.StdEncoding.EncodeToString(r)
	ossConfig := config.GetInstance().OssConfig
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(ossConfig.SecretKey))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	body := fmt.Sprintf("filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}&userId=%d", userId)
	var callbackParam CallbackParam
	callbackParam.CallbackUrl = ossConfig.Callback
	callbackParam.CallbackBody = body
	callbackParam.CallbackBodyType = "application/x-www-form-urlencoded"
	callback_str, err := json.Marshal(callbackParam)
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	callbackBase64 := base64.StdEncoding.EncodeToString(callback_str)

	var policyToken PolicyToken
	policyToken.AccessKeyId = ossConfig.AccessKey
	policyToken.Host = "http://" + ossConfig.Bucket + "." + ossConfig.Endpoint
	policyToken.Expire = expire_end
	policyToken.Signature = string(signedStr)
	policyToken.Directory = prefix
	policyToken.Policy = string(debyte)
	policyToken.Callback = string(callbackBase64)
	response, err := json.Marshal(policyToken)
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}

	result.Ok(response, "成功").Json(ctx)
}

func getUrl(ctx *gin.Context) {
	fileKey := ctx.Param("fileKey")
	if fileKey == "" {
		result.Err("fileKey 为空").Json(ctx)
		return
	}
	singUrl := oss.SingUrl(fileKey)
	result.Ok(singUrl, "").Json(ctx)
}

func upload(ctx *gin.Context) {

	file := &model.Files{}
	if err := ctx.ShouldBindJSON(&file); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}

	// check fileKey not empty
	bucket := oss.GetInstance()
	exist, err := bucket.IsObjectExist(file.FileKey)
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		log.Fatal("判断文件不存在出现异常:", err.Error())
		return
	}
	if !exist {
		result.Err("文件不存在").Json(ctx)
		log.Fatal("文件不存在", err.Error())
		return
	}

	file.UserId = middleware.GetUserId(ctx)

	var fileS services.FileService
	fileS.Save(file)
	result.Ok(nil, "").Json(ctx)
}
