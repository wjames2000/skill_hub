package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/handler"
	"github.com/hpds/skill-hub/internal/middleware"
	"github.com/hpds/skill-hub/internal/router"
)

func main() {
	r := gin.Default()

	r.Use(cors.Default())
	r.Use(middleware.RequestLog())
	r.Use(middleware.Recovery())

	api := r.Group("/api/v1")
	{
		handler.RegisterSkillRoutes(api)
		handler.RegisterSearchRoutes(api)
		handler.RegisterRouterRoutes(api)
		handler.RegisterAuthRoutes(api)
		handler.RegisterUserRoutes(api)
	}

	router.SetupRoutes(r)

	r.Run(":8080")
}
