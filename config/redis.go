package config

import (
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/go-redis/redis"
)

//获得Redis客户端连接
func initRedis() error {

	redisDB := redis.NewClient(&redis.Options{
		Password: G.Config.Redis.Password,
		Addr:     G.Config.Redis.Addr,
		DB:       G.Config.Redis.DB,
	})
	_, err := redisDB.Ping().Result() //ping通了才代表客户端连接成功

	if err != nil {
		panic("failed to connect redis")
	}

	G.RedisDB = redisDB
	return nil

}
