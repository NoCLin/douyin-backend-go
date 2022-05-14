package controller

import (
	"errors"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/NoCLin/douyin-backend-go/utils/json_response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
)

//var userIdSequence = int64(1)
var mutex sync.Mutex

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if len(username) > 32 || len(password) > 32 {
		json_response.Error(c, -1, "the username or password is longer than 32 characters")
		return
	}

	var user model.User
	err := G.DB.Table("users").Where("name = ?", username).Take(&user).Error
	if err == nil {
		json_response.Error(c, -1, "the username already exists.")
		return
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			json_response.Error(c, -1, "unknown error")
		}
	}

	// TODO: hashedPassword
	user = model.User{
		Name:     username,
		Password: password,
	}

	if result := G.DB.Create(&user); result.Error != nil {

		json_response.Error(c, -1, "register failed.")
		log.Fatalf(result.Error.Error())
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
		log.Fatalf("login error", err)
		return
	}

	if user.Password != password {
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
	//userId := c.Query("user_id")
	token := c.Query("token")

	userClaim, err := utils.CheckToken(token)

	if err != nil {
		// TODO: global check
		json_response.Error(c, -1, "forbidden")
		return
	}

	var user model.User

	err = G.DB.Table("users").Where("id = ?", userClaim.UserID).Take(&user).Error
	if err != nil {
		json_response.Error(c, -1, "user not exists")
		return
	}
	json_response.OK(c, "ok", model.UserInfo{
		User:          user,
		FollowCount:   -1,
		FollowerCount: -1,
	})

}

func Test(c *gin.Context) {
	json_response.OK(c, "OK", model.User{
		Model: gorm.Model{ID: 1},
		Name:  "123",
	})
	//json_response.Error(c, -1, "不OK")
}
