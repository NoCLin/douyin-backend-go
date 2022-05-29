package model

type CommentResponse struct {
	Comment
}

type VideoResponse struct {
	Video
	FavoriteCount int64 `json:"favorite_count"`
	CommentCount  int64 `json:"comment_count"    `
	IsFavorite    bool  `json:"is_favorite" `
}

type UserInfo struct {
	User          `json:"user"`
	FollowCount   int64 `json:"follow_count,omitempty"`
	FollowerCount int64 `json:"follower_count,omitempty"`
	IsFollow      bool  `json:"is_follow,omitempty"`
}

type UserLoginResponse struct {
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	User UserInfo `json:"user"`
}

type UserListResponse struct {
	UserList []UserInfo `json:"user_list"`
}

type VideoListResponse struct {
	VideoList []VideoResponse `json:"video_list"`
}

type FeedResponse struct {
	VideoList []VideoResponse `json:"video_list,omitempty"`
	NextTime  int64           `json:"next_time,omitempty"`
}

type CommentListResponse struct {
	CommentList []CommentResponse `json:"comment_list,omitempty"`
}
