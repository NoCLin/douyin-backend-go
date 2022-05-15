package utils

import (
	"context"
	"errors"
	"fmt"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/go-redis/redis/v8"
	"sync"
)

const (
	SPLIT                     = ":"                        //分隔符
	PREFIX_REFRESH_TOKEN      = "douyin:refresh_token"     //refresh token
	PREFIX_USER               = "douyin:user"              //用户信息
	PREFIX_USER_VIDEO         = "douyin:user:video"        //用户发布的视频
	PREFIX_VIDEO_INFO         = "douyin:video:info"        //视频信息
	PREFIX_VIDEO_FAVORITE     = "douyin:video:favorite"    //视频点赞
	PREFIX_VIDEO_FAVORITE_NUM = "douyin:video:favoriteNum" //视频被点赞数目
	PREFIX_VIDEO_POST         = "douyin:video:post"        //视频评论
	PREFIX_USER_RELATION      = "douyin:user:relation"     //关注
	PREFIX_USER_FOLLOWER      = "douyin:user:follower"     //被关注（粉丝）
)

var lock sync.Mutex

// GetRedisDB 这里可以通过函数获得redis的连接 虽然RedisDB 已经注册到了global全局变量中，但是可以在在单独的测试redis时，
//通过此函数初始化，获得redis连接
func GetRedisDB() (*redis.Client, error) {
	if G.RedisDB != nil {
		return G.RedisDB, nil
	}
	lock.Lock()
	defer lock.Unlock()
	if G.RedisDB == nil { //避免并发的建立客户端连接，客户端连接可以只有一个，单例的
		redisDB := redis.NewClient(&redis.Options{})
		_, err := redisDB.Ping(context.Background()).Result()

		if err != nil {
			return nil, errors.New("failed to initialize redis")
		}
		G.RedisDB = redisDB
		return redisDB, nil
	}

	return G.RedisDB, nil

}

//主键生成

// GetTokenKey 存放token的key  key:token
func GetTokenKey(token string) string {
	return PREFIX_REFRESH_TOKEN + SPLIT + token
}

// GetUserKey 把用户信息放在redis的key  key:user
func GetUserKey(id int64) string {
	return PREFIX_USER + SPLIT + fmt.Sprint(id)
}

// GetVideoInfoKey 需要把单独一条视频信息记录时的key key:video
func GetVideoInfoKey(id int64) string {
	return PREFIX_VIDEO_INFO + SPLIT + fmt.Sprint(id)
}

// GetVideoFavoriteKey 需要记录一条视频被点赞时，记录的key key:set
func GetVideoFavoriteKey(id int64) string {
	return PREFIX_VIDEO_FAVORITE + SPLIT + fmt.Sprint(id)
}

// GetVideoFavoriteNumKey 记录一条视频被点赞数量
func GetVideoFavoriteNumKey(id int64) string {
	return PREFIX_VIDEO_FAVORITE_NUM + SPLIT + fmt.Sprint(id)
}

// GetVideoPostKey 需要记录一条视频被评论时，记录的key  key:list
func GetVideoPostKey(id int64) string {
	return PREFIX_VIDEO_POST + SPLIT + fmt.Sprint(id)
}

// GetUserRelationKey 用户关注了那些人  key:set
func GetUserRelationKey(id string) string {
	return PREFIX_USER_RELATION + SPLIT + id
}

// GetUserFollowerKey 用户被那些人关注,被关注者的key，一个用户用有的粉丝  key:set
func GetUserFollowerKey(id string) string {
	return PREFIX_USER_FOLLOWER + SPLIT + id
}

// GetUserVideoKey 存一个用户发送的视频的列表的  key:list
func GetUserVideoKey(id int64) string {
	return PREFIX_USER_VIDEO + SPLIT + fmt.Sprint(id)
}
