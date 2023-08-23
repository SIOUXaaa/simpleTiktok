package minio

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func MakeBucket(ctx context.Context, bucketName string) {
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Println(err)
	}
	if !exists {
		err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Println("make bucket失败: " + err.Error())
			return
		}
		fmt.Println("新建bucket成功")
	}
	fmt.Println("bucket: " + bucketName + "已存在")
}

func GetObjectURL(ctx context.Context, bucketName, objectName string) (*url.URL, error) {
	exp := time.Hour * 24
	url, err := minioClient.PresignedGetObject(ctx, bucketName, objectName, exp, make(url.Values))
	if err != nil {
		fmt.Println("获取object: " + objectName + "失败")
		return nil, err
	}

	return url, err
}

func PutObjectByBuf(ctx context.Context, bucketName, objectName string, buf *bytes.Buffer) (minio.UploadInfo, error) {
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Println("查询bucket失败" + err.Error())
		return minio.UploadInfo{}, err
	}
	if !exists {
		fmt.Println("bucket不存在")
		return minio.UploadInfo{}, nil
	}
	info, err := minioClient.PutObject(ctx, bucketName, objectName, buf, int64(buf.Len()), minio.PutObjectOptions{})
	if err != nil {
		fmt.Println("上传object失败")
		return minio.UploadInfo{}, nil
	}
	return info, nil
}

func PutObject(ctx context.Context, bucketName string, file *multipart.FileHeader) (minio.UploadInfo, error) {
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Println("查询bucket失败" + err.Error())
		return minio.UploadInfo{}, err
	}
	if !exists {
		fmt.Println("bucket不存在")
		return minio.UploadInfo{}, nil
	}
	f, err := file.Open()
	if err != nil {
		fmt.Println("打开文件失败: " + err.Error())
		return minio.UploadInfo{}, err
	}
	info, err := minioClient.PutObject(ctx, bucketName, file.Filename, f, file.Size, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println("上传object失败")
		return minio.UploadInfo{}, err
	}
	return info, nil
}

func Init() {
	var err error
	endpoint := "127.0.0.1:9000"
	accessKeyID := "SztX8sYOGBl9JnCcB7Im"
	secretAccessKey := "nYRmoFcVXvBheSfLXq74VNxCuzuQOrDBh2wiRSio"

	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal("初始化minio错误: " + err.Error())
	}
	MakeBucket(context.Background(), "video")
	MakeBucket(context.Background(), "snapshot")
}
