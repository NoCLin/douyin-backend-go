package controller

import (
	"context"
	"fmt"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	//fmt.Println("publish")
	//fmt.Println(c.PostForm("user_id"))
	//fmt.Println(c.PostForm("token"))
	//c.FormFile("data")
	//fmt.Println("publish")
	//token:= c.Query("token")
	//fmt.Println(token)
	token := c.PostForm("token")
	userInfo, exist := usersLoginInfo[token]

	if !exist {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	data, err := c.FormFile("data")

	//fmt.Println("data",data)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//filename := filepath.Base(data.Filename)
	//user := usersLoginInfo[token]
	temp := strings.Split(data.Filename, ".")
	filetype := temp[1]

	//saveFile := filepath.Join("./public/", finalName)

	//if err := c.SaveUploadedFile(data, saveFile); err != nil {
	//	c.JSON(http.StatusOK, model.Response{
	//		StatusCode: 1,
	//		StatusMsg:  err.Error(),
	//	})
	//	return
	//}
	//fmt.Println(data.Size)
	src, err := data.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	userId := c.PostForm("user_id")
	userIdNum, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		fmt.Printf("invalid userId")
	}
	//uuid,err := uuid.NewRandom()
	if err != nil {
		log.Printf("gen uuid error: %v", err)
	}
	rand.Seed(time.Now().UnixNano())

	//uuidNum,err := strconv.ParseInt(uuid,10,2)
	uniqueId := rand.Int63()
	uniqueIdStr := strconv.FormatInt(uniqueId, 10)
	uniqueIdStr += "." + filetype
	video := &model.Video{
		AuthorID: userIdNum,
		Author:   userInfo.User,
		PlayUrl:  "127.0.0.1:8080/video/" + userId + "/" + uniqueIdStr,
		CoverUrl: "",
	}
	G.DB.Create(video)
	bucketName := "bucket" + userId //bucket不能短于3个字符
	//objectName:=data.Filename
	objectName := uniqueIdStr
	ok, err := G.MinioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: 1,
			StatusMsg:  "BucketExistsError:" + err.Error(),
		})
	}
	if ok == false {

		err = G.MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "cn-east-1"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.Response{
				StatusCode: 1,
				StatusMsg:  "MakeBucketError:" + err.Error(),
			})
		}
	}
	_, err = G.MinioClient.PutObject(context.Background(), bucketName, objectName, src, data.Size, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
		StatusMsg:  "uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	c.JSON(http.StatusOK, model.VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
