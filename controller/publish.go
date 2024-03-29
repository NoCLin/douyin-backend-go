package controller

import (
	"context"
	"encoding/json"
	"fmt"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/NoCLin/douyin-backend-go/utils/json_response"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	"log"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {

	userId := c.GetString("userID")

	data, err := c.FormFile("data")
	if err != nil {
		log.Println("upload error", err)
		json_response.Error(c, 1, "upload error")
		return
	}

	var user model.User
	G.DB.Where("id = ?", userId).First(&user)

	filename := filepath.Base(data.Filename)

	finalName := fmt.Sprintf("%d_%s", user.ID, filename)
	temp := strings.Split(data.Filename, ".")
	filetype := temp[1]

	//fmt.Println(data.Size)
	src, err := data.Open()
	if err != nil {
		log.Println("upload error", err)
		json_response.Error(c, 1, "upload error")
		return
	}

	rand.Seed(time.Now().UnixNano())

	//uuidNum,err := strconv.ParseInt(uuid,10,2)
	uniqueId := rand.Int63()
	uniqueIdStr := strconv.FormatInt(uniqueId, 10)
	imgStr := uniqueIdStr + ".jpeg"
	uniqueIdStr += "." + filetype
	userIdStr := strconv.Itoa(int(user.ID))
	video := &model.Video{
		AuthorID: int64(user.ID),
		Author:   user,
		//PlayUrl: "http://192.168.31.222:9000/bucket" + userIdStr + "/" + uniqueIdStr,
		//CoverUrl: "http://192.168.31.222:9000/bucket" + userIdStr + "/" + imgStr,
		PlayUrl:  G.Config.MinIO.UserAccessUrl + "/bucket" + userIdStr + "/" + uniqueIdStr,
		CoverUrl: G.Config.MinIO.UserAccessUrl + "/bucket" + userIdStr + "/" + imgStr,
	}
	G.DB.Create(video)
	bucketName := "bucket" + userIdStr //bucket不能短于3个字符
	//objectName:=data.Filename
	objectName := uniqueIdStr

	/****************************/
	//上传的视频写到磁盘 因为要调用ffmpeg 没办法
	saveFile := filepath.Join("./public/", finalName)
	fmt.Printf(saveFile)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		log.Println("oss upload error", err)
		json_response.Error(c, 1, "upload error")
		return
	}
	img, imgSize := utils.ExampleReadFrameAsJpeg(saveFile, 1) //根据上传的视频 然后获取第一帧制作封面

	/*****************************/
	ok, err := G.MinioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Println("oss upload error", err)
		json_response.Error(c, 1, "upload error")
		return
	}

	if ok == false {
		err = G.MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "cn-east-1"})
		if err != nil {
			log.Println("oss upload error", err)
			json_response.Error(c, 1, "upload error")
			return
		}
		policy_ := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetBucketLocation","s3:ListBucket","s3:ListBucketMultipartUploads"],"Resource":["arn:aws:s3:::` + bucketName + `"]},{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:AbortMultipartUpload","s3:DeleteObject","s3:GetObject","s3:ListMultipartUploadParts","s3:PutObject"],"Resource":["arn:aws:s3:::` + bucketName + `/*"]}]}`
		//policy,err:=G.MinioClient.GetBucketPolicy( context.Background(), bucketName)
		err := G.MinioClient.SetBucketPolicy(context.Background(), bucketName, policy_)
		if err != nil {
			log.Println("oss upload error", err)
			json_response.Error(c, 1, "upload error")
			return
		}
	}
	//fmt.Printf("%#v", data.Header)
	//textproto.MIMEHeader{"Content-Disposition":[]string{"form-data; name=\"data\"; filename=\"毕业季.mp4\"; filename*=UTF-8''%E6%AF%95%E4%B8%9A%E5%AD%A3.mp4"}, "Content-Type":[]string{"video/mp4"}}
	opt := minio.PutObjectOptions{
		ContentType: data.Header.Get("Content-Type"),
	}
	//fmt.Printf("contentType: ",opt.ContentType)
	//fmt.Printf("objectName: ",objectName)
	_, err = G.MinioClient.PutObject(context.Background(), bucketName, objectName, src, data.Size, opt)
	if err != nil {
		log.Println("oss upload error", err)
		json_response.Error(c, 1, "upload error")
		return
	}
	opt.ContentType = "image/jpeg"
	_, err = G.MinioClient.PutObject(context.Background(), bucketName, imgStr, img, imgSize, opt)
	if err != nil {
		log.Println("oss upload error", err)
		json_response.Error(c, 1, "upload error")
		return
	}
	err = G.RedisDB.Del(context.Background(), "publishlist").Err()
	if err != nil {
		log.Println("del redis cache failed")
	}

	json_response.OK(c, "uploaded successfully", nil)
}

//PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	userId := c.GetString("userID")
	fmt.Printf("publishList   userId: ", userId)
	//userId = "1"
	userIdNum, err := strconv.ParseInt(userId, 10, 64)
	_ = userIdNum
	if err != nil {
		log.Println("string to int failed! ", err)
		json_response.Error(c, -1, "invalid")
		return
	}
	var videos []model.Video
	//G.DB.Where("author_id = ?", userIdNum).Find(&videos)

	redisCache, err := G.RedisDB.Get(context.Background(), "publishlist").Result()
	fmt.Println("redisCache: ", redisCache)
	fmt.Println(len(redisCache))
	if err == redis.Nil || redisCache == "null" { //去mysql拿
		G.DB.Table("videos").Preload("Author").Select("*").Where("author_id = ?", userId).Scan(&videos)
		tempStr, err := json.Marshal(videos)
		if err != nil {
			log.Println("struct to json failed! ", err)
		}
		err = G.RedisDB.Set(context.Background(), "publishlist", tempStr, time.Second*50).Err()
		if err != nil {
			log.Println("set info  to redis failed! ", err)
		}
	} else {
		fmt.Println("get from redis")
		json.Unmarshal([]byte(redisCache), &videos)
	}
	fmt.Println("len:  ", len(videos))
	response := make([]model.VideoResponse, len(videos))
	//交换封面顺序 让前端正常显示
	for i, j := 0, len(videos)-1; i < j; i, j = i+1, j-1 {
		temp := videos[i].CoverUrl
		videos[i].CoverUrl = videos[j].CoverUrl
		videos[j].CoverUrl = temp
	}
	for i := 0; i < len(videos); i++ {
		response[i].Video = videos[i]
		response[i].FavoriteCount = favouriteCount(c, strconv.Itoa(int(videos[i].ID)))

		var count int64 //评论数量
		G.DB.Model(&model.Comment{}).Where("video_id = ? ", videos[i].ID).Count(&count)

		response[i].CommentCount = count
		response[i].IsFavorite = isFavourite(c, strconv.Itoa(int(videos[i].ID)), userId)
	}
	json_response.OK(c, "ok", model.VideoListResponse{
		VideoList: response,
	})

}
