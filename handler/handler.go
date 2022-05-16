package handler

import (
	"fmt"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/NoCLin/douyin-backend-go/utils/json_response"
	"github.com/gin-gonic/gin"
	"log"
	"runtime/debug"
)

func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			//打印错误堆栈信息
			log.Printf("panic: %v\n", err)
			debug.PrintStack()
			json_response.Error(c, -1, "Panic occurs during the api call, for details, please check the server terminal")
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}

func TokenHandler(c *gin.Context) {
	tokenString := c.Query("token")
	if tokenString == "" {
		return
	}

	_, err := utils.CheckToken(tokenString)
	if err != nil {
		log.Printf("token has error: %v\n", err)
		json_response.Error(c, -1, fmt.Sprintf("%v", err))
		//终止后续接口调用
		c.Abort()
	}
}
