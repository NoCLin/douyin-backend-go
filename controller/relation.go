package controller

import (
	"fmt"
	"github.com/NoCLin/douyin-backend-go/config/global"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	tokenString := c.Query("token")
	claim := &utils.Claims{}
	token, _ := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return G.TokenSecret, nil
	})
	_ = token
	userId := claim.UserID
	toUserId := c.Query("to_user_id")
	actionType := c.Query("action_type")
	relationKey := utils.GetUserRelationKey(userId)
	followerKey := utils.GetUserFollowerKey(toUserId)

	if actionType == "1" {
		pipe := global.RedisDB.TxPipeline()
		//关注者的关注列表
		global.RedisDB.SAdd(c, relationKey, toUserId)
		//被关注者粉丝列表
		global.RedisDB.SAdd(c, followerKey, userId)
		_, err := pipe.Exec(c)
		if err != nil {
		}
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 0,
			StatusMsg:  "关注成功",
		})
	} else if actionType == "2" {
		pipe := global.RedisDB.TxPipeline()
		//关注者的关注列表
		global.RedisDB.SRem(c, relationKey, toUserId)
		//被关注者粉丝列表
		global.RedisDB.SRem(c, followerKey, userId)
		_, err := pipe.Exec(c)
		if err != nil {
		}
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 0,
			StatusMsg:  "取关成功",
		})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "无该操作类型"})
	}

}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	user_id := c.Query("user_id")
	relationKey := utils.GetUserRelationKey(user_id)
	es, _ := global.RedisDB.SMembers(c, relationKey).Result()
	len := len(es)
	var user_list = make([]model.UserInfo, len, len)
	for i := 0; i < len; i++ {
		curRelationkey := utils.GetUserRelationKey(es[i])
		curFollowerkey := utils.GetUserFollowerKey(es[i])
		followCount, _ := global.RedisDB.SCard(c, curRelationkey).Result()
		followerCount, _ := global.RedisDB.SCard(c, curFollowerkey).Result()
		var user model.User

		//id, _ := strconv.ParseInt(es[i], 10, 64)

		global.DB.Where("id = ?", es[i]).Find(&user)
		fmt.Println(user.Name)
		user_list[i] = model.UserInfo{
			User: model.User{
				ID:   user.ID,
				Name: user.Name,
			},
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      true,
		}
	}

	c.JSON(http.StatusOK, model.UserListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "获取关注列表成功",
		},
		UserList: user_list,
	})
	return
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	user_id := c.Query("user_id")

	followerKey := utils.GetUserFollowerKey(user_id)
	relationKey := utils.GetUserRelationKey(user_id)
	es, _ := global.RedisDB.SMembers(c, followerKey).Result()
	len := len(es)
	var user_list = make([]model.UserInfo, len, len)
	for i := 0; i < len; i++ {
		curRelationkey := utils.GetUserRelationKey(es[i])
		curFollowerkey := utils.GetUserFollowerKey(es[i])
		followCount, _ := global.RedisDB.SCard(c, curRelationkey).Result()
		followerCount, _ := global.RedisDB.SCard(c, curFollowerkey).Result()
		follow, _ := global.RedisDB.SIsMember(c, relationKey, es[i]).Result()
		var user model.User
		id, _ := strconv.ParseInt(es[i], 10, 64)
		global.DB.Where("id = ?", id).Find(&user)
		user_list[i] = model.UserInfo{
			User: model.User{
				ID:   user.ID,
				Name: user.Name,
			},
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      follow,
		}
	}
	c.JSON(http.StatusOK, model.UserListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "获取粉丝列表成功",
		},
		UserList: user_list,
	})
	return
}
