package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/pkg/errno"
	"github.com/hpds/skill-hub/pkg/logger"
	"github.com/hpds/skill-hub/pkg/response"
)

func StructuredLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info("[GIN]",
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.Int("status", status),
			logger.Duration("latency", latency),
			logger.String("client_ip", c.ClientIP()),
			logger.String("query", query),
			logger.String("user_agent", c.Request.UserAgent()),
		)
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("[GIN] panic recovered",
					logger.Any("error", r),
					logger.String("path", c.Request.URL.Path),
				)
				response.Error(c, errno.InternalError)
				c.Abort()
			}
		}()
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
