package kodo

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/sms/bytes"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
	"xhyovo.cn/community/pkg/config"
)

var authInstance *auth.Credentials

var bucketInstance string

var domainInstance string

func GetDomain() string {
	return domainInstance
}

func Init(kodo *config.KodoConfig) {
	authInstance = auth.New(kodo.AccessKey, kodo.SecretKey)
	bucketInstance = kodo.Bucket
	domainInstance = kodo.Domain
}

func GetAuth() *auth.Credentials {
	return authInstance
}

func getBucketManager() *storage.BucketManager {
	mac := authInstance
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: true,
	}
	return storage.NewBucketManager(mac, &cfg)
}

// 获取文件token
func GetToken() string {
	putPolicy := storage.PutPolicy{
		Scope: bucketInstance,
	}
	return putPolicy.UploadToken(GetAuth())
}

// 删除
func Delete(fileKey string) {
	getBucketManager().Delete(bucketInstance, fileKey)

}

// 获取文件信息
func GetFileInfo(fileKey string) (*storage.FileInfo, error) {
	fileInfo, err := getBucketManager().Stat(bucketInstance, fileKey)
	if err != nil {
		return nil, err
	}
	return &fileInfo, nil

}

func Upload(data []byte, fileKey string) (string, error) {

	upToken := GetToken()
	cfg := storage.Config{}
	// 空间对应的机房

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	dataLen := int64(len(data))
	err := formUploader.Put(context.Background(), &ret, upToken, fileKey, bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		log.Fatalf(err.Error())
		return "", err
	}

	return ret.Key, nil
}
