package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/handler"
	"github.com/hpds/skill-hub/internal/middleware"
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())
	r.Use(middleware.RequestLog())
	r.Use(middleware.Recovery())

	admin := r.Group("/api/v1/admin")
	{
		handler.RegisterAdminRoutes(admin)
	}

	r.Run(":8081")
}
