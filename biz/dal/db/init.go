package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
	"simpleTiktok/pkg/constants"
)

var DB *gorm.DB

func Init(){
	var err error
	DB, err = gorm.Open(mysql.Open(constants.MySQLDefaultDSN),
			&gorm.Config{
				PrepareStmt: true,
				SkipDefaultTransaction: true,
			},
	)
	if err != nil {
		panic(err)
	}
	if err = DB.Use(gormopentracing.New()); err != nil {
		panic(err)
	}
}