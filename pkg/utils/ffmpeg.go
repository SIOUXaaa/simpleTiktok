package utils

import (
	"bytes"
	"os"
	"fmt"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func CreateSnapshot(videoPath string) (buf *bytes.Buffer, err error) {
	buf = bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		fmt.Println("ffmpeg failed: " + err.Error())
		return nil, err
	}
	return buf, nil
}
