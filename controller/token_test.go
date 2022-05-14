package controller

import (
	"log"
	"testing"
	"time"
)

func TestCheckToken(t *testing.T) {
	//进行该测试需要先把Token的有效期设置在1s以内
	token, _ := GenerateToken("xxx", "2")
	_, err := CheckToken(token)
	if err != nil {
		log.Println("检查校验是否顺利通过: ", err)
	}
	time.Sleep(time.Second)
	_, err = CheckToken(token)
	if err != nil {
		log.Println("检查是否过期: ", err)
	}
}
