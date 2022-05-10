package config

import (
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/go-redis/redis"
	"log"
)

func initRedis() error {

	redisDB := redis.NewClient(&redis.Options{})
	_, err := redisDB.Ping().Result()

	if err != nil {
		panic("failed to migrate redis")
	}

	G.RedisDB = redisDB
	log.Println("The redis was connect successfully")
	return nil

}
