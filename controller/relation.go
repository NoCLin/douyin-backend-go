package controller

import (
	"github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")
	toUserId := c.Query("to_user_id")
	actionType := c.Query("action_type")
	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	relationKey := utils.GetUserRelationKey(userId)
	followerKey := utils.GetUserFollowerKey(toUserId)
	if actionType == "1" {
		pipe := global.RedisDB.TxPipeline()
		//关注者的关注列表
		global.RedisDB.SAdd(relationKey, toUserId)
		//被关注者粉丝列表
		global.RedisDB.SAdd(followerKey, userId)
		_, err := pipe.Exec()
		if err != nil {
		}
	} else if actionType == "2" {
		pipe := global.RedisDB.TxPipeline()
		//关注者的关注列表
		global.RedisDB.SRem(relationKey, toUserId)
		//被关注者粉丝列表
		global.RedisDB.SRem(followerKey, userId)
		_, err := pipe.Exec()
		if err != nil {
		}
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "无该操作类型"})
	}

}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	token := c.Query("token")
	user_id := c.Query("user_id")
	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	relationKey := utils.GetUserRelationKey(user_id)
	es, _ := global.RedisDB.SMembers(relationKey).Result()
	len := len(es)
	var user_list = make([]model.UserInfo, len, len)
	for i := 0; i < len; i++ {
		curRelationkey := utils.GetUserRelationKey(es[i])
		curFollowerkey := utils.GetUserFollowerKey(es[i])
		followCount, _ := global.RedisDB.SCard(curRelationkey).Result()
		followerCount, _ := global.RedisDB.SCard(curFollowerkey).Result()
		var user model.User

		//id, _ := strconv.ParseInt(es[i], 10, 64)

		global.DB.Where("id = ?", es[i]).Find(&user)
		user_list[i] = model.UserInfo{
			User: model.User{
				//ID:   id,
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
	token := c.Query("token")
	user_id := c.Query("user_id")
	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	followerKey := utils.GetUserFollowerKey(user_id)
	relationKey := utils.GetUserRelationKey(user_id)
	es, _ := global.RedisDB.SMembers(followerKey).Result()
	len := len(es)
	var user_list = make([]model.UserInfo, len, len)
	for i := 0; i < len; i++ {
		curRelationkey := utils.GetUserRelationKey(es[i])
		curFollowerkey := utils.GetUserFollowerKey(es[i])
		followCount, _ := global.RedisDB.SCard(curRelationkey).Result()
		followerCount, _ := global.RedisDB.SCard(curFollowerkey).Result()
		follow, _ := global.RedisDB.SIsMember(relationKey, es[i]).Result()
		var user model.User
		id, _ := strconv.ParseInt(es[i], 10, 64)
		global.DB.Where("id = ?", id).Find(&user)
		user_list[i] = model.UserInfo{
			User: model.User{
				//Id:   id,
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
