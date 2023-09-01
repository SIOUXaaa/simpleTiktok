package minio

import (
	"bytes"
	"context"
	"os"
	"testing"
)

func TestPutObject(t *testing.T) {
	Init()
	ctx := context.Background()
	fileName := "test"

	file, err := os.Open("/Users/sunzy/Pictures/11111.png")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	// 读取图片数据到 buffer
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)

	uploadInfo, err := PutObjectByBuf(ctx, "snapshot", fileName+".png", buf)
	if err != nil {

	}
	print(uploadInfo.Location)

}
