// Package controller 响应结构
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseData 响应数据
type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// ResponseError 请求失败
func ResponseError(ctx *gin.Context, code ResCode) {
	ctx.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

// ResponseErrorWithMsg 请求失败 带有指定消息
func ResponseErrorWithMsg(ctx *gin.Context, code ResCode, msg interface{}) {
	ctx.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// ResponseSuccess 请求成功
func ResponseSuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}
