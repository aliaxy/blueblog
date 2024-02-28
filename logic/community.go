// Package logic 社区服务层
package logic

import (
	"blueblog/dao/mysql"
	"blueblog/models"
)

// GetCommunityList 获取社区列表
func GetCommunityList() ([]*models.Community, error) {
	// 查数据库 查找到所有的 community 并返回
	return mysql.GetCommunityList()
}

// GetCommunityDetail 获取社区详情
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
