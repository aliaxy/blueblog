// Package controller 请求参数相关
package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CtxUserIDKey 用户 id
const CtxUserIDKey = "userID"

// ErrorUserNotLogin 用户未登录错误
var ErrorUserNotLogin = errors.New("用户未登录")

// 获取当前登录用户的 ID
func getCurrentUser(ctx *gin.Context) (userID int64, err error) {
	uid, ok := ctx.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}

	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

// GetPageInfo  获取分页参数
func GetPageInfo(ctx *gin.Context) (int64, int64) {
	pagetStr := ctx.DefaultQuery("page", "1")
	sizeStr := ctx.DefaultQuery("size", "10")

	page, _ := strconv.ParseInt(pagetStr, 10, 64)
	size, _ := strconv.ParseInt(sizeStr, 10, 64)

	return page, size
}
