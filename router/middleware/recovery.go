package middleware

import (
	"seed-sync/common"
	"seed-sync/log"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// gin全局异常中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := string(debug.Stack())

				// 记录错误日志
				log.Error("[GIN Recovery from panic]",
					zap.Any("error", err),
					zap.String("stack", stack),
				)

				// 直接使用 panic 的错误信息
				errMsg := fmt.Sprintf("%v", err)

				c.JSON(http.StatusOK, common.FailResult(errMsg))
				c.Abort()
			}
		}()
		c.Next()
	}
}
