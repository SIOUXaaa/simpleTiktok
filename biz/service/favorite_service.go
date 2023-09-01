package service

import (
	"context"
	"simpleTiktok/biz/dal/db"
	"simpleTiktok/biz/model/common"
	"simpleTiktok/biz/model/interact/favorite"
	"simpleTiktok/pkg/constants"
	"simpleTiktok/pkg/errno"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"gorm.io/gorm"
)

type FavortieService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewFavoriteService(ctx context.Context, c *app.RequestContext) *FavortieService {
	return &FavortieService{ctx, c}
}

func (f *FavortieService) Action(req *favorite.DouyinFavoriteActionRequest) error {
	userId, _ := f.c.Get("current_user_id")
	videoId := req.GetVideoId()
	actionType := req.GetActionType()

	favorite, err := db.QueryFavoriteByUserIdAndVideoId(userId.(int64), videoId)

	if actionType == constants.FavoriteActionType {
		if err == gorm.ErrRecordNotFound { //之前未给这条视频点赞，需要进行点赞操作
			_, err := db.CreateFavoriteAndIncreaseVideoLikes(&db.Favorites{
				UserId:    userId.(int64),
				VideoId:   videoId,
				CreatedAt: time.Now(),
			})
			if err != nil {
				return err
			}
		} else if err != nil { //其他情况异常
			return err
		}
		//如果已经给这条视频点赞了，就无需处理
		return nil
	} else if actionType == constants.UnFavoriteActionType {
		if err == gorm.ErrRecordNotFound { //之前未给这条视频点赞，报错
			return errno.FavoriteRelationNotExistErr
		} else if err != nil { //其他情况异常
			return err
		}
		//撤销点赞
		_, err := db.DeleteFavoriteAndDecreaseVideoLikes(favorite)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetFavoriteList(userId int64) ([]*common.Video, error) {
	favorites, err := db.GetFavoriteListByUserId(userId)

	if err != nil {
		return nil, err
	}

	videoList := make([]*common.Video, len(favorites))

	for index, value := range favorites {
		id := value.VideoId
		//TODO 获取作者完整信息
		author, err := UserInfoByUserId(value.UserId)
		if err != nil {
			return nil, err
		}
		videoInfo, err := db.GetVideoById(id)
		if err != nil {
			return nil, err
		}
		playUrl := videoInfo.PlayURL
		coverUrl := videoInfo.CoverURL
		favoriteCount := videoInfo.FavoriteCount
		commentCount := videoInfo.CommentCount
		isFavorite := true
		title := videoInfo.Title
		video := common.Video{
			Id:            id,
			Author:        author,
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         title,
		}
		videoList[index] = &video
	}

	return videoList, nil
}
