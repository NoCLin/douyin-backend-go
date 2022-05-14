package controller

import (
	"fmt"
	"github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	//token := c.Query("token")
	userId := c.Query("user_id")
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	commentTtext := c.Query("comment_text")
	commentId := c.Query("comment_id")

	v1 := model.Video{
		AuthorID: 122,
		PlayUrl:  "Ascas",
		CoverUrl: "asas",
	}
	v2 := model.Video{
		AuthorID: 12222,
		PlayUrl:  "Ascasaas",
		CoverUrl: "asas",
	}
	v3 := model.Video{
		AuthorID: 124,
		PlayUrl:  "Asa",
		CoverUrl: "asa",
	}

	global.DB.Create(&v1)
	global.DB.Create(&v2)
	global.DB.Create(&v3)

	//此处先省去token验证阶段

	uId, err := strconv.ParseInt(userId, 10, 64) // 这里转换为 int64 , token验证如若是会用到 userId这一步就会是多余的
	if err != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "error in userId"})
		return
	}

	action, err3 := strconv.Atoi(actionType)
	if err3 != nil {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "error action"})
		return
	}
	if action == 2 { //删除评论
		commId, err4 := strconv.ParseInt(commentId, 10, 64)
		if err4 != nil {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "delete failed"})
			return
		}

		delCommit := model.Comment{
			Id: commId, //这里只提供id便可删除
		}
		global.DB.Delete(&delCommit)
		c.JSON(http.StatusOK, model.Response{StatusCode: 0, StatusMsg: "delete success!"})
		return
	}

	if action == 1 { //发布评论
		vdeId, err2 := strconv.ParseInt(videoId, 10, 64)
		if err2 != nil { //验证视频id
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "error in videoId"})
			return
		}

		//视频首先存在
		var video model.Video
		//global.DB.Debug().First(&video, vdeId)
		global.DB.First(&video, vdeId)
		fmt.Printf("%v", video)
		if video.AuthorID == 0 {
			c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "no  video !"})
			return
		}

		commit := model.Comment{
			UserID:     uId,
			Content:    commentTtext,
			VideoId:    vdeId,
			CreateDate: time.Now().Format("2006-01-02 15:04:05"),
		}

		global.DB.Create(&commit)

		c.JSON(http.StatusOK, model.Response{StatusCode: 0, StatusMsg: "commit sucess !"})
		return

	}
	//前面两个操作都不满足
	c.JSON(http.StatusOK, model.Response{StatusCode: 0, StatusMsg: " illegal operation"})

}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	videoId := c.Query("video_id")

	vdeId, err2 := strconv.ParseInt(videoId, 10, 64)
	if err2 != nil { //验证视频id
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "error in videoId"})
		return
	}
	var video model.Video
	//global.DB.Debug().First(&video, vdeId)
	global.DB.First(&video, vdeId)
	log.Println(vdeId)
	fmt.Printf("%v", video)

	var commites []model.Comment //查询到评论列表
	global.DB.Where("video_id = ? ", vdeId).Find(&commites)

	DemoComments := make([]model.CommentResponse, len(commites))
	index := 0
	for _, co := range commites {
		var author model.User //查询到每条commit的author
		global.DB.Where("id = ?", co.UserID).Find(&author)
		co.User = author

		commentResponsrItem := model.CommentResponse{
			co,
		}
		//DemoComments = append(DemoComments, commentResponsrItem)
		DemoComments[index] = commentResponsrItem
		index++
		fmt.Printf("%v \n", co)
	}
	c.JSON(http.StatusOK, model.CommentListResponse{
		Response:    model.Response{StatusCode: 0},
		CommentList: DemoComments,
	})
}
