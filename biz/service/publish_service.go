package service

import (
	"context"
	"fmt"
	"path"
	"simpleTiktok/biz/model/basic/publish"
	"simpleTiktok/pkg/utils"
	"time"

	"simpleTiktok/biz/dal/db"
	"simpleTiktok/biz/dal/minio"

	"github.com/cloudwego/hertz/pkg/app"
)

type PublishService struct {
	ctx context.Context
	c *app.RequestContext
}

func NewPublishService(ctx context.Context, c *app.RequestContext) *PublishService {
	return &PublishService{ctx, c}
}

func (s *PublishService) PublishAction(req *publish.DouyinPublishActionRequest) error {
	// v, _ := s.c.Get("current_user_id")
	title := req.Title
	// user_id := v.(int64)
	user_id := int64(2)
	t := time.Now()
	fileName := utils.NewFileName(user_id, t.Unix())
	req.Data.Filename = fileName + path.Ext(req.Data.Filename)
	_, err := minio.PutObject(s.ctx, "video", req.Data)
	if err != nil {
		fmt.Println("publish action: " + err.Error())
		return err
	}

	buf, err := utils.CreateSnapshot(utils.URLconvert(s.ctx, s.c, "video/" + fileName))
	if err != nil  {
		fmt.Println("snapshot failed: " + err.Error())
		return err
	}
	_, err = minio.PutObjectByBuf(s.ctx, "snapshot", fileName+".png", buf)
	if err != nil {
		fmt.Println("snapshot upload failed: " + err.Error())
		return err
	}

	_, err = db.CreateVideo(&db.Video{
		AuthorID: user_id,
		PlayURL: "video/" + fileName,
		CoverURL: "snapshot/" + fileName + ".png",
		PublishTime: t,
		Title: title,
	})
	return nil
}