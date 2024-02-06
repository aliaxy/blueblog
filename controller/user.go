package controller

import (
	"net/http"

	"main/logic"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	// 获取参数和校验
	// 业务处理
	logic.SignUp()
	// 返回响应
	c.JSON(http.StatusOK, "ok")
}
