package controller

import (
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/NoCLin/douyin-backend-go/utils/json_response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

//var userIdSequence = int64(1)
var mutex sync.Mutex

func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if len(username) > 32 || len(password) > 32 {
		json_response.Error(c, -1, "the username or password is longer than 32 characters")
		return
	}

	var user model.User
	G.DB.Table("users").Where("name = ?", username).Find(&user)
	if user.Name != "" {
		json_response.Error(c, -1, "the user exists")
		return
	}

	var lastuser model.User

	mutex.Lock()
	G.DB.Table("users").Last(&lastuser)
	// FIXME: 使用自增ID
	user = model.User{
		Id:       lastuser.Id + 1,
		Name:     username,
		Password: password,
	}
	G.DB.Table("users").Create(&user)

	token, err := utils.GenerateToken(user.Name, strconv.FormatInt(user.Id, 10))
	if err != nil {
		json_response.Error(c, -1, "unknown error")
		return
	}

	c.JSON(http.StatusOK, model.UserLoginResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "register successfully",
		},
		UserId: user.Id,
		Token:  token,
	})
	mutex.Unlock()
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user model.User
	G.DB.Table("users").Where("name = ?", username).First(&user)
	if user.Name == "" {
		c.JSON(http.StatusOK, model.UserLoginResponse{
			Response: model.Response{
				StatusCode: -1,
				StatusMsg:  "the username doesn't exist",
			},
			UserId: -1,
		})
		return
	}

	if user.Password != password {
		c.JSON(http.StatusOK, model.UserLoginResponse{
			Response: model.Response{
				StatusCode: -1,
				StatusMsg:  "the password is wrong",
			},
			UserId: -1,
		})
		return
	}

	token, _ := utils.GenerateToken(username, strconv.FormatInt(user.Id, 10))
	//_, err := CheckToken(token)
	//fmt.Println("检验token是否创建成功: ", err)
	//time.Sleep(time.Millisecond * 100)
	//_, err = CheckToken(token)
	//fmt.Println("检验过期功能: ", err)
	//_, err = CheckToken("adjacent")
	//fmt.Println("测试token: ", err)

	c.JSON(http.StatusOK, model.UserLoginResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "login successfully",
		},
		UserId: user.Id,
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

	G.DB.Table("users").Where("id = ?", userClaim.UserID).Find(&user)

	if user.Name == "" {
		c.JSON(http.StatusOK, model.UserResponse{
			Response: model.Response{
				StatusCode: -1,
				StatusMsg:  "the user doesn't exist",
			},
			User: model.UserInfo{},
		})
		return
	}

	c.JSON(http.StatusOK, model.UserResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "get user info successfully",
		},
		User: model.UserInfo{
			User: model.User{
				Id:   user.Id,
				Name: user.Name,
			},
			FollowCount:   -1,
			FollowerCount: -1,
		},
	})
}

func Test(c *gin.Context) {
	json_response.OK(c, "OK", model.User{
		Id:   1,
		Name: "123",
	})
	//json_response.Error(c, -1, "不OK")
}
