package config

import (
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/go-redis/redis"
	"log"
)

//获得Redis客户端连接
func initRedis() error {

	redisDB := redis.NewClient(&redis.Options{})
	_, err := redisDB.Ping().Result() //ping通了才代表客户端连接成功

	if err != nil {
		panic("failed to migrate redis")
	}

	G.RedisDB = redisDB
	log.Println("The redis was connect successfully")
	return nil

}
