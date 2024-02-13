package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimitMiddleware(fillInterval time.Duration, capacity int64) func(ctx *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, capacity)
	return func(ctx *gin.Context) {
		// 取不到令牌返回响应
		if bucket.TakeAvailable(1) < 1 {
			ctx.String(http.StatusOK, "Rate limit passed")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
