package controller

import (
	"strconv"

	"blueblog/logic"
	"blueblog/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子
func CreatePostHandler(ctx *gin.Context) {
	// 获取参数并校验
	post := new(models.Post)
	if err := ctx.ShouldBindJSON(post); err != nil {
		zap.L().Debug("ctx.ShouldBindJSON(post) error", zap.Error(err))
		zap.L().Error("create post with invalid param")
		ResponseError(ctx, CodeInvalidParams)
		return
	}

	// 从 ctx 中取到当前发请求的用户 ID
	userID, err := getCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}
	post.AuthorID = userID

	// 创建帖子
	if err := logic.CreatePost(post); err != nil {
		zap.L().Error("logic.CreatePost(post) failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(ctx, nil)
}

// PostDetailHandler 获取帖子详情
func PostDetailHandler(ctx *gin.Context) {
	// 获取路径参数 帖子 ID
	idStr := ctx.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParams)
		return
	}

	// 根据 ID 查找帖子数据
	data, err := logic.GetPostByID(id)
	if err != nil {
		zap.L().Error("logic.GetPostByID(id) failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(ctx, data)
}
