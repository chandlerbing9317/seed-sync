package middleware

import (
	"seed-sync/log"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TraceLogger 自定义 Gin 的日志中间件
func TraceLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		// 打印请求开始日志
		log.Info(fmt.Sprintf("[GIN  %s begin]", path),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
		)

		c.Next()

		// 打印请求结束日志
		log.Info(fmt.Sprintf("[GIN  %s end]", path),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("cost", time.Since(start)),
			zap.String("ip", c.ClientIP()),
			zap.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}
