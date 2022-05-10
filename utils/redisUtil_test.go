package utils

import (
	"fmt"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"log"
	"testing"
	"time"
)

func TestTestRedis(t *testing.T) {
	r, err := GetRedisDB()
	if err != nil {
		log.Println("err!!!!!!!")
	}

	if r == nil {
		fmt.Println("errrr")
	}
	log.Print("sucessss")

	// 存普通string类型，10分钟过期
	r.Set("test:name", "科科儿子", time.Minute*10)
	// 存hash数据
	r.HSet("test:class", "521", 42)
	// 存list数据
	r.RPush("test:list", 1) // 向右边添加元素
	r.LPush("test:list", 2) // 向左边添加元素
	// 存set数据
	r.SAdd("test:set", "apple")
	G.RedisDB.SAdd("test:set", "pear")
	//g.Close()
}
