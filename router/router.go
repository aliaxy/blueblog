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

	v1 := r.Group("/api/v1")
	// 注册
	v1.POST("/signup", controller.SignUpHandler)

	// 登录
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middleware.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/posts", controller.PostListHandler)
		v1.GET("/posts2", controller.PostListHandler2)
		v1.GET("/post/:id", controller.PostDetailHandler)

		v1.POST("/vote", controller.PostVoteHandler)
	}

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": 404,
		})
	})

	return r
}
