// Package controller 为 api 定义
package controller

import "blueblog/models"

type _ResponsePostList struct {
	Code    ResCode                 `json:"code"`    // 业务响应状态码
	Message string                  `json:"message"` // 响应消息
	Data    []*models.APIPostDetail `json:"data"`    // 数据
}
