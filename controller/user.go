package controller

import (
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin

//var userIdSequence = int64(1)
var mutex sync.Mutex

func Register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"password": password,
	})
	tokenString, _ := token.SigningString()

	if len(username) > 32 || len(password) > 32 {
		c.JSON(http.StatusOK, model.UserLoginResponse{
			Response: model.Response{
				StatusCode: -1,
				StatusMsg:  "the username or password is longer than 32 characters",
			},
			UserId: -1,
			Token:  tokenString,
		})
		return
	}

	var user model.User
	G.DB.Table("users").Where("name = ?", username).Find(&user)
	if user.Name != "" {
		c.JSON(http.StatusOK, model.UserLoginResponse{
			Response: model.Response{
				StatusCode: -1,
				StatusMsg:  "the user has been existed",
			},
			UserId: -1,
			Token:  tokenString,
		})
		return
	}

	var lastuser model.User

	mutex.Lock()
	G.DB.Table("users").Last(&lastuser)
	user = model.User{
		Id:       lastuser.Id + 1,
		Name:     username,
		Password: password,
	}
	G.DB.Table("users").Create(&user)
	c.JSON(http.StatusOK, model.UserLoginResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "register successfully",
		},
		UserId: user.Id,
		Token:  tokenString,
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

	token, _ := GenerateToken(username, strconv.FormatInt(user.Id, 10))
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
	userId := c.Query("user_id")

	var user model.Follow
	G.DB.Table("users").Where("id = ?", userId).Find(&user)

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
			FollowCount:   user.FolloweeId,
			FollowerCount: user.FollowerId,
			IsFollow:      user.IsFollow,
		},
	})
}

func Test(c *gin.Context) {
	user := model.User{Name: "Jinzhu"}

	result := G.DB.Create(&user)
	print(result)
	c.JSON(http.StatusOK, model.UserResponse{
		Response: model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	})
}
