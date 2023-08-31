package service

import (
	"context"
	"fmt"
	"path"
	"simpleTiktok/biz/model/basic/publish"
	"simpleTiktok/biz/model/common"
	"simpleTiktok/pkg/utils"
	"time"

	"simpleTiktok/biz/dal/db"
	"simpleTiktok/biz/dal/minio"

	"github.com/cloudwego/hertz/pkg/app"
)

type PublishService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewPublishService(ctx context.Context, c *app.RequestContext) *PublishService {
	return &PublishService{ctx, c}
}

func (s *PublishService) PublishAction(req *publish.DouyinPublishActionRequest) error {
	v, _ := s.c.Get("current_user_id")
	title := req.Title
	user_id := v.(int64)
	t := time.Now()
	fileName := utils.NewFileName(user_id, t.Unix())
	req.Data.Filename = fileName + path.Ext(req.Data.Filename)
	_, err := minio.PutObject(s.ctx, "video", req.Data)
	if err != nil {
		fmt.Println("publish action: " + err.Error())
		return err
	}

	buf, err := utils.CreateSnapshot(utils.URLconvert("video/" + req.Data.Filename))
	if err != nil {
		fmt.Println("snapshot failed: " + err.Error())
		return err
	}
	_, err = minio.PutObjectByBuf(s.ctx, "snapshot", fileName+".png", buf)
	if err != nil {
		fmt.Println("snapshot upload failed: " + err.Error())
		return err
	}

	_, err = db.CreateVideo(&db.Video{
		AuthorID:    user_id,
		PlayURL:     "video/" + fileName + ".mp4",
		CoverURL:    "snapshot/" + fileName + ".png",
		PublishTime: t,
		Title:       title,
	})
	if err != nil {
		fmt.Println("create video failed: " + err.Error())
	}
	return nil
}

func (s *PublishService) PublishList(req *publish.DouyinPublishListRequest) (resp *publish.DouyinPublishListResponse, err error) {
	resp = &publish.DouyinPublishListResponse{}
	// 获取本次请求的用户id
	query_user_id := req.UserId
	// 判断用户是否登录
	current_user_id, exist := s.c.Get("current_user_id")

	if !exist {
		current_user_id = int64(0)
	}

	dbVideos, err := db.GetVideoByUserID(query_user_id)
	if err != nil {
		return resp, nil
	}

	var videos []*common.Video
	feedService := NewFeedService(s.ctx, s.c)
	err = feedService.CopyVideos(&videos, &dbVideos, current_user_id.(int64))
	if err != nil {
		return resp, nil
	}

	for _, item := range videos {
		video := &common.Video{
			Id: item.Id,
			Author: &common.User{
				Id:              item.Author.Id,
				Name:            item.Author.Name,
				FollowCount:     item.Author.FollowCount,
				FollowerCount:   item.Author.FollowerCount,
				Avatar:          item.Author.Avatar,
				BackgroundImage: item.Author.BackgroundImage,
				Signature:       item.Author.Signature,
				TotalFavorited:  item.Author.TotalFavorited,
				WorkCount:       item.Author.WorkCount,
			},
			PlayUrl:       item.PlayUrl,
			CoverUrl:      item.CoverUrl,
			FavoriteCount: item.FavoriteCount,
			CommentCount:  item.CommentCount,
			IsFavorite:    item.IsFavorite,
			Title:         item.Title,
		}

		resp.VideoList = append(resp.VideoList, video)
	}

	return resp, nil
}
