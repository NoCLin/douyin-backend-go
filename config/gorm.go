package config

import (
	"fmt"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"

	//"gorm.io/driver/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initGorm() error {

	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	//db, err := gorm.Open(mysql.Open("test"), &gorm.Config{})

	args := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		G.Config.Database.User,
		G.Config.Database.Password,
		G.Config.Database.Path,
		G.Config.Database.Database,
		G.Config.Database.Config)
	fmt.Printf("args: ", args)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
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
	return nil
}
