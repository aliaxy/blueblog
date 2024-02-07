package logic

import (
	"blueblog/dao/mysql"
	"blueblog/models"
	"blueblog/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	// 生成 post ID
	post.ID = snowflake.GenID()

	// 保存到数据库
	return mysql.CreatePost(post)
}

func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}

	// 根据作者 ID 查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}

	// 根据社区 ID 查询社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}

	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}

	return
}
