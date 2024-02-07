package middleware

import (
	"strings"

	"blueblog/controller"
	"blueblog/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于 JWT 的中间件认证
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			controller.ResponseError(ctx, controller.CodeNeedLogin)
			ctx.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(ctx, controller.CodeInvalidToken)
			ctx.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(ctx, controller.CodeInvalidToken)
			ctx.Abort()
			return
		}

		ctx.Set(controller.CtxUserIDKey, mc.UserID)
		ctx.Next()
	}
}
