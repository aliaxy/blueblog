// Package mysql 社区相关操作
package mysql

import (
	"database/sql"

	"blueblog/models"

	"go.uber.org/zap"
)

// GetCommunityList 获取社区列表
func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

// GetCommunityDetailByID 通过 id 查询社区详情
func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	sqlStr := `select
		community_id, community_name, introduction, create_time
		from community
		where community_id = ?`
	community = new(models.CommunityDetail)
	if err = db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return
}
