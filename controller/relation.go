package controller

import (
	"fmt"
	"github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/NoCLin/douyin-backend-go/utils/json_response"
	"github.com/gin-gonic/gin"
	"strconv"
)

type UserRes struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type UserListRes struct {
	UserResList []UserRes `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {

	userId := c.GetString("userID")
	toUserId := c.Query("to_user_id")
	actionType := c.Query("action_type")

	if userId == toUserId {
		json_response.Error(c, 1, "invalid")
		return
	}

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
		followCount, _ := global.RedisDB.SCard(c, relationKey).Result()
		followerCount, _ := global.RedisDB.SCard(c, followerKey).Result()
		var user model.User
		global.DB.Model(&user).Where("id = ?", userId).Update("follow_count", followCount)
		global.DB.Model(&user).Where("id = ?", toUserId).Update("follower_count", followerCount)
		json_response.OK(c, "ok", nil)
		return
	} else if actionType == "2" {
		pipe := global.RedisDB.TxPipeline()
		//关注者的关注列表
		global.RedisDB.SRem(c, relationKey, toUserId)
		//被关注者粉丝列表
		global.RedisDB.SRem(c, followerKey, userId)

		_, err := pipe.Exec(c)
		if err != nil {
		}
		followCount, _ := global.RedisDB.SCard(c, relationKey).Result()
		followerCount, _ := global.RedisDB.SCard(c, followerKey).Result()
		var user model.User
		global.DB.Model(&user).Where("id = ?", userId).Update("follow_count", followCount)
		global.DB.Model(&user).Where("id = ?", toUserId).Update("follower_count", followerCount)
		json_response.OK(c, "ok", nil)
		return
	} else {
		json_response.Error(c, 1, "invalid")
		return
	}

}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {

	userId := c.Query("user_id") //这里的user_id 指的是被查询主页信息的 用户 既有登录用户，也有视频作者
	relationKey := utils.GetUserRelationKey(userId)
	es, _ := global.RedisDB.SMembers(c, relationKey).Result()
	length := len(es)
	var user_list = make([]UserRes, length, length)
	for i := 0; i < length; i++ {
		curRelationkey := utils.GetUserRelationKey(es[i])
		curFollowerkey := utils.GetUserFollowerKey(es[i])
		followCount, _ := global.RedisDB.SCard(c, curRelationkey).Result()
		followerCount, _ := global.RedisDB.SCard(c, curFollowerkey).Result()
		var user model.User
		//id, _ := strconv.ParseInt(es[i], 10, 64)
		global.DB.Where("id = ?", es[i]).Find(&user)
		fmt.Println(user.Name)
		user_list[i] = UserRes{
			ID:            user.ID,
			Name:          user.Name,
			FollowCount:   uint(followCount),
			FollowerCount: uint(followerCount),
			IsFollow:      true,
		}
	}
	json_response.OK(c, "ok", UserListRes{
		UserResList: user_list,
	})
	return
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {

	//userId := c.GetString("userID")
	userId := c.Query("user_id") //理由同上
	followerKey := utils.GetUserFollowerKey(userId)
	relationKey := utils.GetUserRelationKey(userId)

	es, _ := global.RedisDB.SMembers(c, followerKey).Result()
	length := len(es)
	var user_list = make([]UserRes, length, length)
	for i := 0; i < length; i++ {
		curRelationkey := utils.GetUserRelationKey(es[i])
		curFollowerkey := utils.GetUserFollowerKey(es[i])
		followCount, _ := global.RedisDB.SCard(c, curRelationkey).Result()
		followerCount, _ := global.RedisDB.SCard(c, curFollowerkey).Result()
		follow, _ := global.RedisDB.SIsMember(c, relationKey, es[i]).Result()
		var user model.User
		id, _ := strconv.ParseInt(es[i], 10, 64)
		global.DB.Where("id = ?", id).Find(&user)
		user_list[i] = UserRes{
			ID:            user.ID,
			Name:          user.Name,
			FollowCount:   uint(followCount),
			FollowerCount: uint(followerCount),
			IsFollow:      follow,
		}
	}
	json_response.OK(c, "ok", UserListRes{
		UserResList: user_list,
	})
	return
}

//返回有isFollow字段的调这个函数，userId为当前用户id，toUserId为要检查是否关注的ID
func isFollow(c *gin.Context, userId string, toUserid string) bool {
	token := c.Query("token")
	if token == "" {
		token = c.PostForm("token")
	}
	if token ==""{
		return false
	}
	relationKey := utils.GetUserRelationKey(userId)
	follow, _ := global.RedisDB.SIsMember(c, relationKey, toUserid).Result()
	return follow
}
