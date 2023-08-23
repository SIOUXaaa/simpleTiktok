package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"simpleTiktok/biz/model/basic/feed"
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

func (s *FeedService) Feed(req *feed.DouyinFeedRequest) (*feed.DouyinFeedRequest, error) {

	return req, nil
}
