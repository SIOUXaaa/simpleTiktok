package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"simpleTiktok/biz/dal/db"
	"simpleTiktok/biz/model/basic/feed"
	"simpleTiktok/biz/model/common"
	"simpleTiktok/pkg/constants"
	"simpleTiktok/pkg/utils"
	"time"
)

type FeedService struct {
	ctx context.Context     // 传递上下文信息
	c   *app.RequestContext // 请求信息
}

func NewFeedService(ctx context.Context, c *app.RequestContext) *FeedService {
	return &FeedService{
		ctx: ctx,
		c:   c,
	}
}

func (s *FeedService) Feed(req *feed.DouyinFeedRequest) (*feed.DouyinFeedResponse, error) {
	resp := &feed.DouyinFeedResponse{}
	var lastTime time.Time
	if req.LatestTime == 0 {
		lastTime = time.Now()
	} else {
		lastTime = time.Unix(req.LatestTime/1000, 0)
	}
	fmt.Printf("LastTime: %v\n", lastTime)
	// 获取访问的用户id
	current_user_id, exists := s.c.Get("current_user_id")
	if !exists {
		// 如果用户没有登录 则将用户id设置为0
		current_user_id = int64(0)
	}
	dbVideos, err := db.GetVideosByLastTime(lastTime)
	if err != nil {
		return resp, err
	}

	videos := make([]*common.Video, 0, constants.VideoFeedCount)
	err = s.CopyVideos(&videos, &dbVideos, current_user_id.(int64))
	if err != nil {
		return resp, nil
	}
	resp.VideoList = videos
	if len(dbVideos) != 0 {
		resp.NextTime = dbVideos[len(dbVideos)-1].PublishTime.Unix()
	}
	return resp, nil
}

func (s *FeedService) CopyVideos(result *[]*common.Video, data *[]*db.Video, userId int64) error {
	for _, item := range *data {
		video := s.createVideos(item, userId)
		*result = append(*result, video)
	}
	return nil
}

func (s *FeedService) createVideos(data *db.Video, userId int64) *common.Video {
	video := &common.Video{
		Id:       data.ID,
		PlayUrl:  utils.URLconvert(s.ctx, s.c, data.PlayURL),
		CoverUrl: utils.URLconvert(s.ctx, s.c, data.CoverURL),
		Title:    data.Title,
	}

	// TODO
	//var wg sync.WaitGroup
	//wg.Add(4)

	//go func() {
	//defer wg.Done()

	//author, err := NewUserService(s.ctx, s.c).GetUserInfo(data.AuthorID, userId)
	//if err != nil {
	//	log.Printf("GetUserInfo func error:" + err.Error())
	//}
	//video.Author = &common.User{
	//	Id:              author.Id,
	//	Name:            author.Name,
	//	FollowCount:     author.FollowCount,
	//	FollowerCount:   author.FollowerCount,
	//	IsFollow:        author.IsFollow,
	//	Avatar:          author.Avatar,
	//	BackgroundImage: author.BackgroundImage,
	//	Signature:       author.Signature,
	//	TotalFavorited:  author.TotalFavorited,
	//	WorkCount:       author.WorkCount,
	//	FavoriteCount:   author.FavoriteCount,
	//}
	//}()

	video.Author = &common.User{
		Id:              1000,
		Name:            "test",
		FollowCount:     0,
		FollowerCount:   0,
		IsFollow:        false,
		Avatar:          "",
		BackgroundImage: "",
		Signature:       "",
		TotalFavorited:  0,
		WorkCount:       0,
		FavoriteCount:   0,
	}

	video.CommentCount = 0
	video.FavoriteCount = 0
	video.IsFavorite = false

	return video
}
