package config

import (
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func initGorm() error {

	//db, err := gorm.Open(sqlite.Open("testdb.db"), &gorm.Config{})
	dsn := "root:root@tcp(127.0.0.1:3306)/db2?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Video{},
		&model.Comment{},
		&model.Follow{},
	)
	if err != nil {
		panic("failed to migrate database")
	}

	G.DB = db
	log.Println("The database was initialized successfully")
	return nil

}
