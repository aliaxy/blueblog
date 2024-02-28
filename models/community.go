// Package models 社区模型
package models

import "time"

// Community 社区结构
type Community struct {
	ID   int64  `json:"id,string" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

// CommunityDetail 社区详情
type CommunityDetail struct {
	ID           int64     `json:"id,string" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
