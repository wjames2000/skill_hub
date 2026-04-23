package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		gin.DefaultWriter.Write([]byte(c.Request.Method + " " + c.Request.URL.Path + " " + latency.String() + "\n"))
	}
}

func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
