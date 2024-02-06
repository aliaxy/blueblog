package controller

import (
	"net/http"

	"blueblog/logic"
	"blueblog/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 处理注册请求
func SignUpHandler(ctx *gin.Context) {
	// 获取参数和校验
	p := new(models.ParamSignUp)
	if err := ctx.ShouldBindJSON(p); err != nil {
		// 返回错误响应
		zap.L().Error("Sign up with invalid param", zap.Error(err))
		// 判断 error 是不是 validator.ValidationErros 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"msg": removeTopStruct(errs.Translate(trans)), // 进行翻译
			})
		}
		return
	}
	// 业务处理
	if err := logic.SignUp(p); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "注册失败",
		})
		return
	}
	// 返回响应
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
