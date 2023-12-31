package dal

import (
	"simpleTiktok/biz/dal/db"
	"simpleTiktok/biz/dal/minio"
)

func Init() {
	db.Init()
	db.DB.AutoMigrate(&db.User{}, &db.Video{}, &db.Comment{}, &db.Favorites{})
	minio.Init()
}
