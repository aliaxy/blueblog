package logic

import (
	"blueblog/dao/mysql"
	"blueblog/models"
	"blueblog/pkg/snowflake"
)

func CreatePost(post *models.Post) (err error) {
	// 生成 post ID
	post.ID = snowflake.GenID()

	// 保存到数据库
	return mysql.CreatePost(post)
}

func GetPostByID(pid int64) (data *models.Post, err error) {
	return mysql.GetPostByID(pid)
}
