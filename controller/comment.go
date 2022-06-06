package controller

import (
	"fmt"
	"github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils/json_response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {

	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	commentText := c.Query("comment_text")
	commentId := c.Query("comment_id")

	//userId已经有中间件验证
	uId, _ := strconv.ParseInt(c.GetString("userID"), 10, 64)

	action, err3 := strconv.Atoi(actionType)
	if err3 != nil { //错误操作
		json_response.Error(c, 1, "faulty operation")
		return
	}
	if action == 2 { //删除评论
		commId, err4 := strconv.ParseInt(commentId, 10, 64)
		if err4 != nil {
			json_response.Error(c, 1, "comment does not exist")
			return
		}

		delCommit := model.Comment{
			Model: model.Model{
				ID: uint(commId), //这里只提供id便可删除
			},
		}

		global.DB.Delete(&delCommit)
		json_response.OK(c, "ok", nil)
		return
	}

	if action == 1 { //发布评论
		vdeId, err2 := strconv.ParseInt(videoId, 10, 64)
		if err2 != nil { //验证视频id
			json_response.Error(c, 1, "video does not exist")
			return
		}

		//视频首先存在
		var video model.Video
		//global.DB.Debug().First(&video, vdeId)
		global.DB.First(&video, vdeId)
		fmt.Printf("%v", video)
		if video.AuthorID == 0 {
			json_response.Error(c, 1, "video does not exist")
			return
		}

		//过滤评论
		commentText = global.WordFilter.Replace(commentText, '*')

		commit := model.Comment{
			UserID:  uId,
			Content: commentText,
			VideoId: vdeId,
		}

		global.DB.Create(&commit)
		comm := make([]model.CommentResponse, 1)
		comm[0] = model.CommentResponse{
			Id:          commit.ID,
			User:        commit.User,
			Content:     commit.Content,
			CreatedDate: commit.CreatedAt.Format("01-02"),
		}
		json_response.OK(c, "ok", model.CommentActionResponse{
			CommentResponse: model.CommentResponse{
				Id:          commit.ID,
				User:        commit.User,
				Content:     commit.Content,
				CreatedDate: commit.CreatedAt.Format("01-02"),
			},
		})

		return

	}
	//前面两个操作都不满足
	json_response.Error(c, 1, " illegal operation")
	return

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	videoID := c.Query("video_id")

	vdeId, err2 := strconv.ParseInt(videoID, 10, 64)
	if err2 != nil { //验证视频id
		json_response.Error(c, -1, "invalid videoID")
		return
	}

	// TODO: make sure video exists

	var comments []model.Comment //查询到评论列表
	global.DB.Preload("User").Where("video_id = ? ", vdeId).Find(&comments)

	ret := make([]model.CommentResponse, len(comments))
	index := 0
	for _, co := range comments {
		//var author model.User //查询到每条commit的author
		//global.DB.Where("id = ?", co.UserID).Find(&author)
		//co.User = author
		//ret = append(ret, commentResponsrItem)

		ret[index] = model.CommentResponse{
			Id:          co.ID,
			User:        co.User,
			Content:     co.Content,
			CreatedDate: co.CreatedAt.Format("01-02"),
		}
		index++
	}
	json_response.OK(c, "ok", model.CommentListResponse{
		CommentList: ret,
	})
}
