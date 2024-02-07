package controller

import (
	"strconv"

	"blueblog/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler
func CommunityHandler(ctx *gin.Context) {
	// 查询到所有的社区 (community_id, community_name) 以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy) // 不轻易把服务器报错返回给外
		return
	}
	ResponseSuccess(ctx, data)
}

// CommunityDetailHandler 社区分类详情
func CommunityDetailHandler(ctx *gin.Context) {
	// 获取社区 ID
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(ctx, CodeInvalidParams)
		return
	}

	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy) // 不轻易把服务器报错返回给外
		return
	}
	ResponseSuccess(ctx, data)
}
