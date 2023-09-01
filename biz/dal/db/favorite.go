package db

import (
	"fmt"
	"time"

	"simpleTiktok/pkg/constants"

	"gorm.io/gorm"
)

type Favorites struct {
	ID        int64          `json:"id"`
	UserId    int64          `json:"user_id"`
	VideoId   int64          `json:"video_id"`
	CreatedAt time.Time      `json:"create_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

func (Favorites) TableName() string {
	return constants.FavoritesTableName
}

func QueryFavoriteByUserIdAndVideoId(userId int64, videoId int64) (*Favorites, error) {
	var favorite = Favorites{}
	if err := DB.Where("user_id = ? AND video_id = ?", userId, videoId).Find(&favorite).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}
	return &favorite, nil
}

func CreateFavoriteAndIncreaseVideoLikes(favorite *Favorites) (int64, error) {
	if err := DB.Select("UserId", "VideoId", "CreateAt").Create(&favorite).Error; err != nil {
		return 0, err
	}
	var video = Video{}
	if err := DB.First(&video, favorite.VideoId).Error; err != nil {
		return 0, err
	}
	video.FavoriteCount++
	if err := DB.Save(&video).Error; err != nil {
		return 0, err
	}
	return favorite.ID, nil
}

func DeleteFavoriteAndDecreaseVideoLikes(favorite *Favorites) (int64, error) {
	id := favorite.ID
	if err := DB.Delete(&favorite).Error; err != nil {
		return 0, err
	}
	var video = Video{}
	if err := DB.First(&video, favorite.VideoId).Error; err != nil {
		return 0, err
	}
	video.FavoriteCount--
	if err := DB.Save(&video).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func GetFavoriteListByUserId(userId int64) ([]*Favorites, error) {
	var favorites []*Favorites
	if err := DB.Where("user_id = ?", userId).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

// 获取获赞总数
func GetTotalFavorited(user_id int64) (int64, error) {
	var favorite_count int64
	videos, err := GetVideoByUserID(user_id)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	for _, video := range videos {
		favorite_count += video.FavoriteCount
	}
	return favorite_count, nil
}

// 获取喜欢数
func GetFavoriteCount(user_id int64) (int64, error) {
	var favorite_count int64
	if err := DB.Model(&Favorites{}).Where("user_id = ?", user_id).Count(&favorite_count).Error; err != nil {
		return -1, err
	}
	return favorite_count, nil
}

func QueryIsFavorite(userId int64, videoId int64) (bool, error) {
	var favorite Favorites
	err := DB.Model(&Favorites{}).Where("user_id = ? AND video_id = ?", userId, videoId).First(&favorite).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetFavoriteCountOfVideo(videoId int64) (int64, error) {
	var count int64
	if err := DB.Model(&Favorites{}).Where("video_id = ?", videoId).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}
