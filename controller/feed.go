package controller

import (
	"context"
	"fmt"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	//查询mysql获取最
	//G.DB.Select()
	VideoList := make([]model.Video, 3)
	fmt.Println(len(VideoList))
	G.DB.Order("created_at desc").Limit(30).Find(&VideoList)
	fmt.Println(len(VideoList))
	fmt.Println(VideoList[0])
	fmt.Println(VideoList[1])
	fmt.Println(VideoList[2])
	//TODO: FavoriteCount CommentCount  IsFavorite
	//fmt.Println(len(VideoList),len(VideoResponse))
	nums := 30
	if len(VideoList) < nums {
		nums = len(VideoList)
	}
	VideoResponse := make([]model.VideoResponse, nums)
	for i := 0; i < nums; i++ {
		VideoResponse[i].Video = VideoList[i]
	}
	feed := model.FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: VideoResponse,
		NextTime:  time.Now().Unix(),
	}
	c.JSON(http.StatusOK, feed)
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
