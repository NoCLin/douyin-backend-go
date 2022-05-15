package config

import (
	"context"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
)

//获得Redis客户端连接
func initRedis() *redis.Client {

	redisDB := redis.NewClient(&redis.Options{
		Password: G.Config.Redis.Password,
		Addr:     G.Config.Redis.Addr,
		DB:       G.Config.Redis.DB,
	})
	_, err := redisDB.Ping(context.Background()).Result() //ping通了才代表客户端连接成功

	if err != nil {
		panic("failed to connect redis")
	}

	return redisDB

}

func initTestRedis() (*redis.Client, redismock.ClientMock) {

	redisDB, redisMock := redismock.NewClientMock()
	return redisDB, redisMock

}
