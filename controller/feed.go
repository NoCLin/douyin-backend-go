package controller

import (
	"context"
	"fmt"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils/json_response"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"os"
	"time"
)

func Feed(c *gin.Context) {

	// TOOD: latest_time

	var videos []model.Video

	// TODO: 过滤已删除
	// TODO: 关联查询 user
	// TODO: 填充字段
	G.DB.Order("created_at desc").Limit(30).Find(&videos)

	responseVideos := make([]model.VideoResponse, len(videos))
	for i := 0; i < len(videos); i++ {
		responseVideos[i].Video = videos[i]
		//responseVideos[i].CreatedAt = VideoList[i].CreatedAt
		responseVideos[i].Author.Name = "xxx"
		responseVideos[i].Author.ID = uint(videos[i].AuthorID)
		responseVideos[i].Author.FollowCount = 999
		responseVideos[i].Author.FollowerCount = 999
		responseVideos[i].FavoriteCount = 1
		responseVideos[i].CommentCount = 1
		responseVideos[i].IsFavorite = true
	}
	fmt.Println(len(videos))
	if len(videos) == 0 {
		responseVideos = DemoVideos
	}
	feed := model.FeedResponse{
		VideoList: responseVideos,
		NextTime:  time.Now().Unix(),
	}
	json_response.OK(c, "ok", feed)
	return
}

func GetVideo(c *gin.Context) {
	bucketName := "bucket" + c.Param("user")
	objectName := c.Param("filename")
	obj, err := G.MinioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"message": err,
		})
		return
	}
	//fmt.Println(obj.Stat())
	//先写到本地
	localfile := "./public/" + bucketName + "_" + objectName
	f, err := os.Create(localfile)
	if _, err = io.Copy(f, obj); err != nil {
		fmt.Println(err)
		return
	}

	c.File(localfile)
}

//func GetVideo(c *gin.Context){
//	//file := G.MinioClient.GetObject(context.Background(),)
//	bucketName := "mymusic"   //一个用户名一个桶
//	objectName := "test.jpg"  //对象名称
//	objectName = c.Param("filename")
//	fmt.Println(objectName)
//	//filePath :="E:\\DeepLearning\\dataset\\BSDS300\\minio\\download.jpg"//存储路径
//	basePath := "E:\\DeepLearning\\dataset\\BSDS300\\minio\\"
//	filePath := path.Join(basePath,objectName)
//	err:=G.MinioClient.FGetObject(context.Background(),bucketName,objectName,filePath,minio.GetObjectOptions{})
//	if err!=nil{
//		log.Println("getfile failed: ",err)
//		c.JSON(http.StatusNotFound,gin.H{
//			"message": "not found this vedio",
//		})
//		return
//	}
//	c.File(filePath)
//}
