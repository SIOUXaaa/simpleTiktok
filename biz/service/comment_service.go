package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"simpleTiktok/biz/dal/db"
	"simpleTiktok/biz/model/common"
	"simpleTiktok/biz/model/interact/comment"
	"simpleTiktok/pkg/errno"
)

type CommentService struct {
	ctx context.Context
	c   *app.RequestContext
}

// NewCommentService create comment service
func NewCommentService(ctx context.Context, c *app.RequestContext) *CommentService {
	return &CommentService{ctx: ctx, c: c}
}

func (c *CommentService) AddNewComment(req *comment.DouyinCommentActionRequest) (*comment.Comment, error) {
	currentUserID, _ := c.c.Get("current_user_id")
	videoID := req.VideoId
	actionType := req.ActionType
	commentText := req.CommentText
	commentID := req.CommentId

	resComment := &comment.Comment{}

	//发表评论
	if actionType == 1 {
		dbComment := &db.Comment{
			UserId:      currentUserID.(int64),
			VideoId:     videoID,
			CommentText: commentText,
		}
		err := db.AddNewComment(dbComment)
		if err != nil {
			return resComment, err
		}
		resComment.Id = dbComment.ID
		resComment.CreateDate = dbComment.CreatedAt.Format("01-02")
		resComment.Content = dbComment.CommentText
		resComment.User, err = c.getUserInfoById(currentUserID.(int64))
		if err != nil {
			return resComment, err
		}
		return resComment, nil
	} else { //删除评论
		err := db.DeleteCommentById(commentID)
		if err != nil {
			return resComment, err
		}
		return resComment, nil
	}
}

func (c *CommentService) getUserInfoById(userId int64) (*common.User, error) {
	u, err := UserInfoByUserId(userId)
	var comment_user *common.User
	if err != nil {
		return comment_user, err
	}
	comment_user = &common.User{
		Id:              u.Id,
		Name:            u.Name,
		FollowCount:     u.FollowCount,
		FollowerCount:   u.FollowerCount,
		IsFollow:        u.IsFollow,
		Avatar:          u.Avatar,
		BackgroundImage: u.BackgroundImage,
		Signature:       u.Signature,
		TotalFavorited:  u.TotalFavorited,
		WorkCount:       u.WorkCount,
		FavoriteCount:   u.FavoriteCount,
	}
	return comment_user, nil
}

func (c *CommentService) CommentList(req *comment.DouyinCommentListRequest) (*comment.DouyinCommentListResponse, error) {
	resp := &comment.DouyinCommentListResponse{}
	videoId := req.VideoId

	currentUserId, _ := c.c.Get("current_user_id")

	dbComments, err := db.GetCommentListByVideoID(videoId)
	if err != nil {
		return resp, err
	}
	var comments []*comment.Comment
	err = c.copyComment(&comments, &dbComments, currentUserId.(int64))
	if err != nil {
		return resp, err
	}
	resp.CommentList = comments
	resp.StatusMsg = errno.SuccessMsg
	resp.StatusCode = errno.SuccessCode
	return resp, nil
}

func (c *CommentService) copyComment(result *[]*comment.Comment, data *[]*db.Comment, currentUserId int64) error {
	for _, item := range *data {
		comment := c.createComment(item, currentUserId)
		*result = append(*result, comment)
	}
	return nil
}

func (c *CommentService) createComment(data *db.Comment, userId int64) *comment.Comment {
	resComment := &comment.Comment{
		Id:         data.ID,
		Content:    data.CommentText,
		CreateDate: data.CreatedAt.Format("01-02"),
	}

	userInfo, err := c.getUserInfoById(userId)
	if err != nil {
		log.Printf("func error")
	}
	resComment.User = userInfo
	return resComment
}
