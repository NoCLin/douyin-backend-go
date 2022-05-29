package controller

import (
	"github.com/NoCLin/douyin-backend-go/model"
)

var DemoVideos = []model.VideoResponse{
	{
		Video: model.Video{
			Author:   DemoUser.User,
			PlayUrl:  "https://www.w3schools.com/html/movie.mp4",
			CoverUrl: "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		},
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
}

var DemoUser = model.UserInfo{
	User: model.User{
		Name: "TestUser",
	},
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}

var usersLoginInfo = map[string]model.UserInfo{
	"zhangleidouyin": {
		User: model.User{
			Name: "zhanglei",
		},
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}
