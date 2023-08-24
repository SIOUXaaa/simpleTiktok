package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/url"
	"simpleTiktok/pkg/constants"
	"time"
)

var (
	Client *minio.Client
	err    error
)

func MakeBucket(ctx context.Context, bucketName string) {
	exists, err := Client.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !exists {
		err = Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Successfully created mybucket %v\n", bucketName)
	}
}

func GetObjURL(ctx context.Context, bucketName, filename string) (u *url.URL, err error) {
	exp := time.Hour * 24
	reqParams := make(url.Values)
	u, err = Client.PresignedGetObject(ctx, bucketName, filename, exp, reqParams)
	return u, err
}

func Init() {
	ctx := context.Background()
	Client, err = minio.New(constants.MinioEndPoint, &minio.Options{
		Creds: credentials.NewStaticV4(constants.MinioAccessKeyID, constants.MinioSecretAccessKey, ""),
		//Secure: constants.MiniouseSSL
	})
	if err != nil {
		log.Fatalln("minio连接错误: ", err)
	}
	log.Printf("%#v\n", Client)

	MakeBucket(ctx, "test")
}
