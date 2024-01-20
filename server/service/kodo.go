package services

import (
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Kodo struct{}

var kodo *kodoObject

type kodoObject struct {
	AccessKey string
	SecretKey string
	Bucket    string
}

// 初始化Kodo
func InitKodo(accessKey, secretKey, bucket string) {
	kodo = &kodoObject{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
	}
}

func getAuth() *auth.Credentials {
	return auth.New(kodo.AccessKey, kodo.SecretKey)
}

func getBucketManager() *storage.BucketManager {
	mac := getAuth()
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: true,
	}
	return storage.NewBucketManager(mac, &cfg)
}

// 获取文件token
func (*Kodo) GetToken() string {
	bucket := kodo.Bucket
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	return putPolicy.UploadToken(getAuth())
}

// 删除
func (*Kodo) Delete(fileKey string) {
	getBucketManager().Delete(kodo.Bucket, fileKey)

}

// 获取文件信息
func (*Kodo) GetFileInfo(fileKey string) (*storage.FileInfo, error) {
	fileInfo, err := getBucketManager().Stat(kodo.Bucket, fileKey)
	if err != nil {
		return nil, err
	}
	return &fileInfo, nil

}
