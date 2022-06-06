package config

import (
	"fmt"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func loadConfig(configKey string) {

	viper.SetConfigName(configKey)
	viper.AddConfigPath("./config")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error while reading config file: %v", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Configuration file changed")
	})

	if err := viper.Unmarshal(&G.Config); err != nil {
		panic(fmt.Errorf("fatal error while decode config file: %v", err))
	}

}

func InitConfig() {
	env := os.Getenv("DOUYIN_ENV")
	if env == "" {
		env = "myConfig"
	}

	loadConfig(env)

	db := initGorm(G.Config.Database)
	G.DB = db
	G.DB = db.Debug()

	redisDB := initRedis()

	G.RedisDB = redisDB

	minioClient := initMinIO()
	G.MinioClient = minioClient

	G.WordFilter = initSensitiveTree()
}

func InitTestConfig() {

	mockDB, mock := initTestGorm()
	G.DB = mockDB.Debug()
	G.DBMock = mock

	redisDB, redisMock := initTestRedis()

	G.RedisDB = redisDB
	G.RedisMock = redisMock

	minioClient := initTestMinio()
	G.MinioClient = minioClient
}
