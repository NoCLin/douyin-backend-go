package controller

import (
	"github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/gin-gonic/gin"
	"net/http"

)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	//relationKey := utils.GetUserRelationKey(userId)
	//followerKey := utils.GetUserFollowerKey(toUserId)
	userFavoriteKey := utils.GetVideoFavoriteKey(userId)
	videoBeFavouriteKey :=utils.GetVideoFavoriteNumKey(videoId)
	if actionType == "1" {
		pipe := global.RedisDB.TxPipeline()
		//用户喜欢的视频列表
		global.RedisDB.SAdd(userFavoriteKey, videoId)
		//视频被哪些粉丝点赞
		global.RedisDB.SAdd(videoBeFavouriteKey, userId)
		_, err := pipe.Exec()
		if err != nil {
		}
	} else if actionType == "2" {
		pipe := global.RedisDB.TxPipeline()
		//用户喜欢的视频列表
		global.RedisDB.SRem(userFavoriteKey, videoId)
		//视频被哪些粉丝点赞
		global.RedisDB.SRem(videoBeFavouriteKey, userId)
		_, err := pipe.Exec()
		if err != nil {
		}
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "无该操作类型"})
	}

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")
	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	userFavoriteKey := utils.GetVideoFavoriteKey(userId)
	es, _ := global.RedisDB.SMembers(userFavoriteKey).Result()
	len := len(es)
	var favorite_list = make([]model.Video, len, len)
	for i := 0; i < len; i++ {
		curVideoFavoriteNum := utils.GetVideoFavoriteNumKey(es[i])
		//curFollowerkey := utils.GetUserFollowerKey(es[i])
		videoFavoriteNum, _ := global.RedisDB.SCard(curVideoFavoriteNum).Result()
		//followerCount, _ := global.RedisDB.SCard(curFollowerkey).Result()
		var video model.Video



		global.DB.Where("id = ?", es[i]).Find(&video)
		favorite_list[i] = model.VideoResponse{
			Video:model.Video{
				AuthorID: video.AuthorID,
				Author:   video.Author,
				PlayUrl:  video.PlayUrl,
				CoverUrl: video.CoverUrl,
			},
			FavoriteCount :videoFavoriteNum,
			CommentCount  :0,
			IsFavorite :true,
			}

	}
		c.JSON(http.StatusOK, model.VideoListResponse{
			Response: model.Response{
				StatusCode: 0,
				StatusMsg:  "获取点赞列表成功",
			},
			VideoList:favorite_list,
		})



		return
	}



}
