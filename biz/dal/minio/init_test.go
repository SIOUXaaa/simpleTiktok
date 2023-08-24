package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"testing"
)

func TestGetObjectURL(t *testing.T) {
	Init()
	ctx := context.Background()
	bucketName := "video"
	objectName := "1.mp4"
	opts := minio.GetObjectOptions{}
	url, err := GetObjectURL(ctx, bucketName, objectName, opts)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(url.String())
}
