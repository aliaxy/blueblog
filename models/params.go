// Package models 请求参数相关
package models

const (
	OrderTime  = "time"  // OrderTime 时间排序
	OrderScore = "score" // OrderScore 分数排序
)

// 定义请求参数的结构体

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	PostID    string `json:"post_id" binding:"required,min=1"`
	Direction int8   `json:"direction" binding:"oneof=1 0 -1"`
}

// ParamPostList 获取帖子列表参数
type ParamPostList struct {
	CommunityID int64  `json:"community_id,string" form:"community_id"` // 可以为空
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
}

// ParamCommunityPostList 根据社区
// type ParamCommunityPostList struct {
// 	*ParamPostList
// }
