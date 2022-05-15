package config

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
)

func initGorm(d G.Database) *gorm.DB {
	var db *gorm.DB
	var err error

	switch strings.ToLower(d.Type) {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			d.User, d.Password, d.Host, d.Port, d.Database,
		)
		mysqlConfig := mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}
		db, err = gorm.Open(mysql.New(mysqlConfig), &gorm.Config{})

	case "postgres":
		// not tested
		// TODO: sslmode TimeZone in config
		dsn := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			d.Host, d.Port, d.User, d.Password, d.Database)

		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})

	case "sqlite":
		db, err = gorm.Open(sqlite.Open(d.Database), &gorm.Config{})

	default:
		panic("Unsupported DB type")
	}

	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database")
	}

	sqlDB.SetMaxIdleConns(d.MaxIdleConns)
	sqlDB.SetMaxOpenConns(d.MaxOpenConns)

	if G.Config.Database.AutoMigrate {
		err = db.AutoMigrate(
			&model.User{},
			&model.Video{},
			&model.Comment{},
			&model.Follow{},
		)
		if err != nil {
			panic("failed to migrate database")
		}
	}

	return db

}

func initTestGorm() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	// no regex
	//db, mock, _ = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	// https://github.com/DATA-DOG/go-sqlmock/issues/118
	mockDB, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})

	return mockDB, mock
}
