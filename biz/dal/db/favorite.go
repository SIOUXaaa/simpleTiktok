package db

import (
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

func QueryFavoriteByUserIdAndVedioId(userId int64, videoId int64) (*Favorites, error) {
	var favorite = Favorites{}
	if err := DB.Where("user_id = ? AND vedio_id = ?", userId, videoId).Find(&favorite).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}
	return &favorite, nil
}

func CreateFavoriteAndIncreaseVedioLikes(favorite *Favorites) (int64, error) {
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

func DeleteFavoriteAndDecreaseVedioLikes(favorite *Favorites) (int64, error) {
	id := favorite.ID
	if err := DB.Delete(favorite).Error; err != nil {
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
