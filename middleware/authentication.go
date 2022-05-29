package middleware

import (
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/NoCLin/douyin-backend-go/utils/json_response"
	"github.com/gin-gonic/gin"
	"strings"
)

var UnCheckList = []string{"/douyin/feed", "/douyin/user/login", "/douyin/user/register"}

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		for _, s := range UnCheckList {
			if strings.Contains(path, s) {
				c.Next()
				return
			}
		}

		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
		}
		//fmt.Println("token: ",token)
		userClaim, err := utils.CheckToken(token)
		if err != nil {
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
