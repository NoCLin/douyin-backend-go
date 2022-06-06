package controller

import (
	"github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/NoCLin/douyin-backend-go/utils/json_response"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	userId := c.GetString("userID")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	//relationKey := utils.GetUserRelationKey(userId)
	//followerKey := utils.GetUserFollowerKey(toUserId)
	userFavoriteKey := utils.GetVideoFavoriteKey(userId)
	videoBeFavouriteKey := utils.GetVideoFavoriteNumKey(videoId)
	if !(actionType == "1" || actionType == "2") {
		json_response.Error(c, 1, "无该操作类型")
		return
	}
	pipe := global.RedisDB.TxPipeline()
	if actionType == "1" {
		//用户喜欢的视频列表
		global.RedisDB.SAdd(c, userFavoriteKey, videoId)
		//视频被哪些粉丝点赞
		global.RedisDB.SAdd(c, videoBeFavouriteKey, userId)

	} else {
		//用户喜欢的视频列表
		global.RedisDB.SRem(c, userFavoriteKey, videoId)
		//视频被哪些粉丝点赞
		global.RedisDB.SRem(c, videoBeFavouriteKey, userId)
	}
	_, err := pipe.Exec(c)
	if err != nil {
		json_response.Error(c, -1, "error")
		return
	}
	json_response.OK(c, "ok", nil)
	return

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	userId := c.Query("user_id")

	userFavoriteKey := utils.GetVideoFavoriteKey(userId)
	es, _ := global.RedisDB.SMembers(c, userFavoriteKey).Result()
	length := len(es)
	var favoriteList = make([]model.VideoResponse, length, length)

	for i := 0; i < length; i++ {
		curVideoFavoriteNum := utils.GetVideoFavoriteNumKey(es[i])
		videoFavoriteNum, _ := global.RedisDB.SCard(c, curVideoFavoriteNum).Result()

		var video model.Video
		global.DB.Preload("Author").Where("id = ?", es[i]).Find(&video)
		video.Author.IsFollow = isFollow(c, userId, strconv.Itoa(int(video.AuthorID)))

		var commentcount int64 //评论数量
		global.DB.Model(&model.Comment{}).Where("video_id = ? ", video.ID).Count(&commentcount)

		favoriteList[i] = model.VideoResponse{
			Video: model.Video{
				AuthorID: video.AuthorID,
				Author:   video.Author,
				PlayUrl:  video.PlayUrl,
				CoverUrl: video.CoverUrl,
			},
			FavoriteCount: videoFavoriteNum,
			CommentCount:  commentcount,
			IsFavorite:    true,
		}
	}

	json_response.OK(c, "ok", model.VideoListResponse{
		VideoList: favoriteList,
	})

	return
}

//返回有isFavourite字段的调这个函数，userId为当前用户id，videoId为视频Id
func isFavourite(c *gin.Context, videoId string, userId string) bool {
	videoBeFavouriteKey := utils.GetVideoFavoriteNumKey(videoId)
	isFavourite, _ := global.RedisDB.SIsMember(c, videoBeFavouriteKey, userId).Result()

	log.Println(isFavourite)
	return isFavourite
}

//从redis返回实时更新的点赞数
func favouriteCount(c *gin.Context, videoId string) int64 {
	videoBeFavouriteKey := utils.GetVideoFavoriteNumKey(videoId)
	favouriteCount, _ := global.RedisDB.SCard(c, videoBeFavouriteKey).Result()
	log.Println(favouriteCount)
	return favouriteCount
}
