package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"xhyovo.cn/community/pkg/log"
)

var bucket *oss.Bucket

var en string

func GetEndpoint() string {
	return en
}

func GetInstance() *oss.Bucket {
	return bucket
}

func Init(endpoint, accessKey, accessSec, bucketN string) {
	client, err := oss.New(endpoint, accessKey, accessSec)
	if err != nil {
		log.Errorf("初始化 oss 失败,err: %s", err.Error())
		panic(err.Error())
		return
	}
	bucket, err = client.Bucket(bucketN)

	if err != nil {
		log.Errorf("获取 oss bucket 失败,err: %s", err.Error())
		panic(err.Error())
		return
	}
	en = endpoint
}

func SingUrl(fileKey string) string {
	singUrl, err := bucket.SignURL(fileKey, oss.HTTPGet, 3600)
	if err != nil {
		log.Errorf("获取 oss bucket 失败,err: %s", err.Error())
		return ""
	}
	return singUrl
}
