package model

import (
	"gorm.io/gorm"
)

// TODO: https://gorm.io/zh_CN/docs/indexes.html

type Video struct {
	gorm.Model

	// TODO: HashID
	AuthorID int64
	Author   User   `json:"author"     `
	PlayUrl  string `json:"play_url"            `
	CoverUrl string `json:"cover_url,omitempty"   `
}

type Comment struct {
	gorm.Model

	UserID     int64
	User       User   `json:"user" `
	Content    string `json:"content,omitempty" `
	CreateDate string `json:"create_date,omitempty" `
}

type User struct {
	gorm.Model
	ID   uint   `gorm:"primarykey;AUTO_INCREMENT"`
	Name string `json:"name"`
	// TODO: ignore Password in JSON
	Password       string `json:"-" `
	PasswordHashed string `json:"-" `
}

type Follow struct {
	gorm.Model

	Name       string `json:"name,omitempty"`
	FollowerId int64  `json:"follower_id,omitempty" ` // 关注人
	FolloweeId int64  `json:"followee_id,omitempty" ` // 被关注人
	IsFollow   bool   `json:"is_follow,omitempty"`
}
