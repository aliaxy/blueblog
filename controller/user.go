package controller

import (
	"errors"
	"strconv"

	"blueblog/dao/mysql"
	"blueblog/logic"
	"blueblog/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 处理注册请求
func SignUpHandler(ctx *gin.Context) {
	// 获取参数和校验
	p := new(models.ParamSignUp)
	if err := ctx.ShouldBindJSON(p); err != nil {
		// 返回错误响应
		zap.L().Error("Sign up with invalid param", zap.Error(err))

		// 判断 error 是不是 validator.ValidationErros 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParams)
		} else {
			// 进行翻译
			ResponseErrorWithMsg(ctx, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		}
		return
	}

	// 业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(ctx, CodeUserExist)
		} else {
			ResponseError(ctx, CodeServerBusy)
		}
		return
	}

	// 返回响应
	ResponseSuccess(ctx, nil)
}

// LoginHandler 处理登录请求
func LoginHandler(ctx *gin.Context) {
	// 获取参数和校验
	p := new(models.ParamLogin)
	if err := ctx.ShouldBindJSON(p); err != nil {
		// 返回错误响应
		zap.L().Error("Login with invalid param", zap.Error(err))

		// 判断 error 是不是 validator.ValidationErros 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParams)
		} else {
			ResponseErrorWithMsg(ctx, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		}
		return
	}

	// 业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(ctx, CodeUserNotExist)
		} else {
			ResponseError(ctx, CodeInvalidPassword)
		}
		return
	}

	// 返回响应
	ResponseSuccess(ctx, gin.H{
		"user_id":   strconv.FormatInt(user.UserID, 10),
		"user_name": user.Username,
		"token":     user.Token,
	})
}
