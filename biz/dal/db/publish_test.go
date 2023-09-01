package db

import (
	"fmt"
	"testing"
	"time"
)

func TestGetVideosByLastTime(t *testing.T) {
	Init()
	videos, err := GetVideosByLastTime(time.Now())
	if err != nil {

	}
	fmt.Println(len(videos))
	for _, v := range videos {
		fmt.Println(v)
	}
}
