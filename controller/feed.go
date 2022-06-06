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
	"strconv"
	"time"
)

func Feed(c *gin.Context) {

	// TODO: latest_time

	var videos []model.Video

	// TODO: 过滤已删除
	// TODO: 关联查询 user
	// TODO: 填充字段
	G.DB.Preload("Author").Order("created_at desc").Limit(30).Find(&videos)

	userId := c.GetString("userID")

	var responseVideos []model.VideoResponse
	if len(videos) == 0 {
		responseVideos = DemoVideos
	} else {
		responseVideos = make([]model.VideoResponse, len(videos))
		for i := 0; i < len(videos); i++ {
			v := videos[i]

			responseVideos[i].Video = v
			//TODO:判断是否登录，未登录isfollow都为false
			v.Author.IsFollow = isFollow(c, userId, strconv.Itoa(int(v.AuthorID)))
			responseVideos[i].Author = v.Author
			// TODO: real data
			responseVideos[i].FavoriteCount = favouriteCount(c, strconv.Itoa(int(v.ID)))
			var count int64 //评论数量
			G.DB.Model(&model.Comment{}).Where("video_id = ? ", v.ID).Count(&count)
			responseVideos[i].CommentCount = count

			//TODO:判断是否登录，未登录isFavorite都为false
			responseVideos[i].IsFavorite = isFavourite(c, strconv.Itoa(int(v.ID)), userId)
		}
	}
	var returnTime int64
	if len(videos) == 0 {
		returnTime = time.Now().Unix()
	} else {
		returnTime = videos[0].CreatedAt.Unix()
	}
	feed := model.FeedResponse{
		VideoList: responseVideos,
		NextTime:  returnTime, //本次返回视频的最新时间
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
