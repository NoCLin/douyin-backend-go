package middleware

import (
	"fmt"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/NoCLin/douyin-backend-go/utils/json_response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		// TODO: whitelist array, regex

		if c.Request.URL.Path == "/douyin/feed" ||
			c.Request.URL.Path == "/douyin/feed/" ||
			c.Request.URL.Path == "/douyin/user/login" || // for client
			c.Request.URL.Path == "/douyin/user/login/" ||
			c.Request.URL.Path == "/douyin/user/register" ||
			c.Request.URL.Path == "/douyin/user/register/" {
			c.Next()
			return
		}

		token := c.Query("token")
		if token == ""{
			token = c.PostForm("token")
		}
		//fmt.Println("token: ",token)
		userClaim, err := utils.CheckToken(token)
		if err != nil {
			fmt.Println("验证不通过！")
			json_response.Error(c, -1, "forbidden")
			// 若验证不通过，不再调用后续的函数处理
			c.Abort()
			return
		}

		c.Set("userID", userClaim.UserID)
		c.Set("username", userClaim.Username)
		c.Next()

	}
}
