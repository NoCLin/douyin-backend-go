package controller

import (
	"errors"
	"fmt"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/NoCLin/douyin-backend-go/utils/json_response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
)

func Register(c *gin.Context) {

	username := c.Query("username")
	password := c.Query("password")

	if len(username) > 32 || len(password) > 32 {
		json_response.Error(c, -1, "the username or password is longer than 32 characters")
		return
	}
	if len(username) <= 0 || len(password) < 5 {
		json_response.Error(c, -1, "the username or password is too short")
		return
	}

	var user model.User
	err := G.DB.Table("users").Where("name = ?", username).Take(&user).Error
	if err == nil {
		json_response.Error(c, -1, "the username already exists.")
		return
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println(err)
			json_response.Error(c, -1, "unknown error")
			return
		}
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		json_response.Error(c, -1, "register failed")
	}

	user = model.User{
		Name:           username,
		PasswordHashed: hashedPassword,
	}

	result := G.DB.Create(&user)
	if result.Error != nil {
		log.Println("register insert failed.", result.Error)
		json_response.Error(c, -1, "register failed.")
		return
	}

	token, err := utils.GenerateToken(user.Name, strconv.FormatInt(int64(user.ID), 10))
	if err != nil {
		json_response.Error(c, -1, "unknown error")
		return
	}

	json_response.OK(c, "ok", model.UserLoginResponse{
		UserId: int64(user.ID),
		Token:  token,
	})

}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	var user model.User
	err := G.DB.Table("users").Where("name = ?", username).Take(&user).Error

	if err != nil {
		json_response.Error(c, -1, "the username doesn't exist")
		return
	}
	//if user.Password != password {
	//	json_response.Error(c, -1, "login failed")
	//	return
	//}
	// 使用密码的哈希值来验证
	if !utils.CheckPasswordHash(password, user.PasswordHashed) {
		json_response.Error(c, -1, "login failed")
		return
	}

	token, _ := utils.GenerateToken(username, strconv.FormatInt(int64(user.ID), 10))
	//_, err := CheckToken(token)
	//fmt.Println("检验token是否创建成功: ", err)
	//time.Sleep(time.Millisecond * 100)
	//_, err = CheckToken(token)
	//fmt.Println("检验过期功能: ", err)
	//_, err = CheckToken("adjacent")
	//fmt.Println("测试token: ", err)

	json_response.OK(c, "ok", model.UserLoginResponse{
		UserId: int64(user.ID),
		Token:  token,
	})
}

func UserInfo(c *gin.Context) {

	var user model.User
	cur_userId := c.GetString("userID")

	err := G.DB.Table("users").Where("id = ?", c.MustGet("userID")).Take(&user).Error
	if err != nil {
		json_response.Error(c, -1, "user not exists")
		return
	}
	user.IsFollow = isFollow(c, cur_userId, strconv.Itoa(int(user.ID)))
	curRelationkey := utils.GetUserRelationKey(strconv.Itoa(int(user.ID)))
	curFollowerkey := utils.GetUserFollowerKey(strconv.Itoa(int(user.ID)))
	followCount, _ := G.RedisDB.SCard(c, curRelationkey).Result()
	followerCount, _ := G.RedisDB.SCard(c, curFollowerkey).Result()
	json_response.OK(c, "ok", model.UserInfo{
		User:          user,
		FollowCount:   followCount,
		FollowerCount: followerCount,
	})

}

func Test(c *gin.Context) {
	json_response.OK(c, "OK", model.User{
		Model: model.Model{ID: 1},
		Name:  "123",
	})
	//json_response.Error(c, -1, "不OK")
}
