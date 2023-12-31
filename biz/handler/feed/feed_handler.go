// Code generated by hertz generator.

package feed

import (
	"context"
	"simpleTiktok/biz/service"
	"simpleTiktok/pkg/errno"
	"simpleTiktok/pkg/utils"

	feed "simpleTiktok/biz/model/basic/feed"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Feed .
// @router /douyin/feed/ [GET]
func Feed(ctx context.Context, c *app.RequestContext) {
	var err error
	var req feed.DouyinFeedRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.JSON(consts.StatusOK, feed.DouyinFeedResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
		return
	}

	resp, err := service.NewFeedService(ctx, c).Feed(&req)
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.JSON(consts.StatusOK, feed.DouyinFeedResponse{
			StatusCode: resp.StatusCode,
			StatusMsg:  resp.StatusMsg,
		})
	}
	resp.StatusCode = errno.SuccessCode
	resp.StatusMsg = errno.SuccessMsg

	c.JSON(consts.StatusOK, resp)
}
