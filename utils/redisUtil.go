package utils

import (
	"errors"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/go-redis/redis"
	"sync"
)

const (
	SPLIT                 = ":"              //分隔符
	PREFIX_TOKEN          = "token"          //登录token
	PREFIX_USER           = "user"           //用户信息
	PREFIX_VIDEO_INFO     = "video:info"     //视频信息
	PREFIX_VIDEO_FAVORITE = "video:favorite" //视频点赞
	PREFIX_VIDEO_POST     = "video:post"     //视频评论
	PREFIX_USER_RELATION  = "user:relation"  //关注
	PREFIX_USER_FOLLWER   = "user:follower"  //被关注（粉丝）
	PREFIX_USER_VIDEO     = "user:video"     //用户发布的视频
)

var lock sync.Mutex

// GetRedisDB 这里可以通过函数获得redis的连接 虽然RedisDB 已经注册到了global全局变量中，但是可以在在单独的测试redis时，
//通过此函数初始化，获得redis连接
func GetRedisDB() (*redis.Client, error) {
	if G.RedisDB != nil {
		return G.RedisDB, nil
	}
	lock.Lock() //如果是在在项目中也使用GetRedisDB获得redis连接，则考虑并发的创建redis连接，得加锁
	defer lock.Unlock()
	if G.RedisDB == nil {
		redisDB := redis.NewClient(&redis.Options{})
		_, err := redisDB.Ping().Result()

		if err != nil {
			return nil, errors.New("failed to initialize redis")
		}
		G.RedisDB = redisDB
		return redisDB, nil
	}

	return G.RedisDB, nil

}

// GetTokenKey 存放token的key  key:token
func GetTokenKey(token string) string {
	return PREFIX_TOKEN + SPLIT + token
}

// GetUserKey 把用户信息放在redis的key  key:user
func GetUserKey(id int64) string {
	return PREFIX_USER + SPLIT + string(id)
}

// GetVideoInfoKey 需要把单独一条视频信息记录时的key ket:video
func GetVideoInfoKey(id int64) string {
	return PREFIX_VIDEO_INFO + SPLIT + string(id)
}

// GetVideoFavoriteKey 需要记录一条视频被点赞时，记录的key ket:set
func GetVideoFavoriteKey(id int64) string {
	return PREFIX_VIDEO_FAVORITE + SPLIT + string(id)
}

// GetVideoPostKey 需要记录一条视频被评论时，记录的key  key:list
func GetVideoPostKey(id int64) string {
	return PREFIX_VIDEO_POST + SPLIT + string(id)
}

// GetUserRelationKey 用户关注了那些人  key:set
func GetUserRelationKey(id int64) string {
	return PREFIX_USER_RELATION + SPLIT + string(id)
}

// GetUserFollowerKey 用户被那些人关注,被关注者的key，一个用户用有的粉丝  key:set
func GetUserFollowerKey(id int64) string {
	return PREFIX_USER_FOLLWER + SPLIT + string(id)
}

// GetUserVideoKey 存一个用户发送的视频的列表的  key:list
func GetUserVideoKey(id int64) string {
	return PREFIX_USER_RELATION + SPLIT + string(id)
}
