package controller

import (
	"context"
	"fmt"
	G "github.com/NoCLin/douyin-backend-go/config/global"
	"github.com/NoCLin/douyin-backend-go/model"
	"github.com/NoCLin/douyin-backend-go/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	tokenString := c.PostForm("token")
	claim := &utils.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return G.TokenSecret, nil
	})
	_=token
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	userIdPar  := claim.UserID

	var user model.User
	G.DB.Where("id = ?", userIdPar).First(&user)

	//userInfo:=model.UserInfo{
	//	User: user,
	//}
	filename := filepath.Base(data.Filename)
	//user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.ID, filename)
	temp := strings.Split(data.Filename, ".")
	filetype := temp[1]




	//fmt.Println(data.Size)
	src, err := data.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	//userId := c.PostForm("user_id")
	//userIdNum, err := strconv.ParseInt(userId, 10, 64)
	//if err != nil {
	//	fmt.Printf("invalid userId")
	//}
	//uuid,err := uuid.NewRandom()
	//if err != nil {
	//	log.Printf("gen uuid error: %v", err)
	//}
	rand.Seed(time.Now().UnixNano())

	//uuidNum,err := strconv.ParseInt(uuid,10,2)
	uniqueId := rand.Int63()
	uniqueIdStr := strconv.FormatInt(uniqueId, 10)
	imgStr := uniqueIdStr+".jpeg"
	uniqueIdStr += "." + filetype
	userIdStr := strconv.Itoa(int(user.ID))
	video := &model.Video{
		AuthorID: int64(user.ID),
		Author:   user,
		//PlayUrl:  "127.0.0.1:8080/video/" + userId + "/" + uniqueIdStr,
		//PlayUrl:  "192.168.252.100:9000/bucket" + userId + "/" + uniqueIdStr,
		PlayUrl:  "http://192.168.31.222:9000/bucket" + userIdStr + "/" + uniqueIdStr,
		//PlayUrl:  "http://192.168.31.222:9000/bucket27/bear.mp4" ,
		//PlayUrl: "192.168.31.222:8080/video/1/movie1-3.mp4",
		CoverUrl: "http://192.168.31.222:9000/bucket" + userIdStr + "/"+imgStr,
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
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	img,imgSize := utils.ExampleReadFrameAsJpeg(saveFile,1) //根据上传的视频 然后获取第一帧制作封面

	/*****************************/
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
		policy_ := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetBucketLocation","s3:ListBucket","s3:ListBucketMultipartUploads"],"Resource":["arn:aws:s3:::`+bucketName+`"]},{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:AbortMultipartUpload","s3:DeleteObject","s3:GetObject","s3:ListMultipartUploadParts","s3:PutObject"],"Resource":["arn:aws:s3:::`+bucketName+`/*"]}]}`
		//policy,err:=G.MinioClient.GetBucketPolicy( context.Background(), bucketName)
		err := G.MinioClient.SetBucketPolicy(context.Background(), bucketName, policy_)
		if err != nil {
			fmt.Println(err)
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
	opt.ContentType = "image/jpeg"
	_, err = G.MinioClient.PutObject(context.Background(), bucketName, imgStr, img, imgSize,  opt)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {

	tokenString := c.Query("token")
	fmt.Printf("tokenString: ",tokenString)
	claim := &utils.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return G.TokenSecret, nil
	})
	if err!=nil{
		c.JSON(http.StatusInternalServerError, model.VideoListResponse{
			Response: model.Response{
				StatusCode: 1,
			},
			VideoList: nil,
		})
	}
	_=token
	userId := claim.UserID
	fmt.Printf("userId: ",userId)
	//userId = "1"
	userIdNum,err := strconv.ParseInt(userId,10,64)
	_ = userIdNum
	if err!=nil{
		log.Println("string to int failed! ",err)
	}
	var videos []model.Video
	//G.DB.Where("author_id = ?", userIdNum).Find(&videos)
	G.DB.Table("videos").Select("*").Where("author_id = ?", userId).Scan(&videos)
	fmt.Println( "len:  ",len(videos))
	//fmt.Printf("%#v",videos[0])

	response := make([]model.VideoResponse,len(videos))

	for i:=0;i<len(videos);i++{
		response[i].Video = videos[i]
		response[i].Author.Name = claim.Username
		response[i].Author.ID = 999
		response[i].Author.FollowCount = 999
		response[i].Author.FollowerCount = 999
		response[i].FavoriteCount =100
		response[i].CommentCount = 100
		response[i].IsFavorite = false
	}
	fmt.Printf("%#v",response)
	c.JSON(http.StatusOK, model.VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg: "get list success",
		},
		VideoList: response,
	})
}
