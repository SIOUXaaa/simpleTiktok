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

//  获取获赞总数
func GetTotalFavorited(user_id int64) (int64, error) {
	var favorite_count int64
	videos, err := GetVideoByUserID(user_id)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	for _, video := range videos  {
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
