package dal

import (
	"simpleTiktok/biz/dal/db"
)

func Init() {
	db.Init()
	db.DB.AutoMigrate(&db.User{}, &db.Video{}, &db.Comment{}, &db.Favorites{})
}
