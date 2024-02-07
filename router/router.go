package router

import (
	"net/http"

	"blueblog/controller"
	"blueblog/logger"
	"blueblog/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.POST("/signup", controller.SignUpHandler)

	r.POST("/login", controller.LoginHandler)

	r.GET("/ping", middleware.JWTAuthMiddleware(), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": 404,
		})
	})

	return r
}
