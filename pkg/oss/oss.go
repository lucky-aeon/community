package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
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
		log.Fatalln(err.Error())
		return
	}
	bucket, err = client.Bucket(bucketN)

	if err != nil {
		log.Fatal(err.Error())
		return
	}
	en = endpoint
}

func SingUrl(fileKey string) string {
	singUrl, err := bucket.SignURL(fileKey, oss.HTTPGet, 12000)
	if err != nil {
		log.Println(err)
	}
	return singUrl
}
