package controller

import (
	"strconv"

	"blueblog/logic"
	"blueblog/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

// PostListHandler
func PostListHandler(ctx *gin.Context) {
	page, size := GetPageInfo(ctx)

	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(ctx, data)
}

// PostVoteHandler
func PostVoteHandler(ctx *gin.Context) {
	// 获取数据
	p := new(models.ParamVoteData)
	if err := ctx.ShouldBindJSON(p); err != nil {
		zap.L().Error("ctx.ShouldBindJSON(p) failed", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParams)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(ctx, CodeInvalidParams, errData)
		return
	}

	// 获取用户 ID
	userID, err := getCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeNeedLogin)
		return
	}

	// 具体业务
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost(userID, p) failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(ctx, nil)
}

// PostListHandler2 升级的帖子列表接口
// 根据前端传来的参数动态获取帖子列表（创建时间 or 分数）
func PostListHandler2(ctx *gin.Context) {
	// /api/v1/posts?page=1&size=10&order=time
	// 获取分页参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	if err := ctx.ShouldBindQuery(p); err != nil {
		zap.L().Error("PostListHandler2 with invalid params", zap.Error(err))
		ResponseError(ctx, CodeInvalidParams)
		return
	}

	// 获取数据
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(ctx, data)
}
