package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/handler"
	"github.com/hpds/skill-hub/internal/middleware"
	"github.com/hpds/skill-hub/internal/router"
	"github.com/hpds/skill-hub/pkg/config"
	"github.com/hpds/skill-hub/pkg/consul"
	"github.com/hpds/skill-hub/pkg/db"
	"github.com/hpds/skill-hub/pkg/logger"
	mls "github.com/hpds/skill-hub/pkg/meilisearch"
	"github.com/hpds/skill-hub/pkg/minio"
	rds "github.com/hpds/skill-hub/pkg/redis"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		panic("load config: " + err.Error())
	}

	logger.Init(cfg.App.Env)
	logger.Info("router-api starting", logger.String("app", cfg.App.Name), logger.String("env", cfg.App.Env))

	dbEngine, err := db.WaitForDB(cfg.DB, 30*time.Second)
	if err != nil {
		logger.Fatal("db init", logger.String("error", err.Error()))
	}
	defer dbEngine.Close()
	logger.Info("database connected")

	redisClient, err := rds.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Warn("redis unavailable, rate limiter disabled", logger.String("error", err.Error()))
		redisClient = nil
	} else {
		defer redisClient.Close()
		logger.Info("redis connected")
	}

	meiliClient, err := mls.New(cfg.Meili.Host, cfg.Meili.APIKey)
	if err != nil {
		logger.Warn("meilisearch unavailable, search disabled", logger.String("error", err.Error()))
		meiliClient = nil
	} else {
		logger.Info("meilisearch connected")
	}

	minioClient, err := minio.New(cfg.Minio.Endpoint, cfg.Minio.AccessKey, cfg.Minio.SecretKey, cfg.Minio.UseSSL, cfg.Minio.Bucket)
	if err != nil {
		logger.Warn("minio unavailable, file upload disabled", logger.String("error", err.Error()))
		minioClient = nil
	} else {
		logger.Info("minio connected")
	}

	consulClient, err := consul.New(cfg.Consul.Addr, cfg.Consul.Token)
	if err != nil {
		logger.Warn("consul unavailable, service discovery disabled", logger.String("error", err.Error()))
		consulClient = nil
	} else {
		serviceID := fmt.Sprintf("router-api-%d", cfg.App.Port)
		if err := consulClient.RegisterService(
			cfg.App.Name+"-router", serviceID, "localhost", cfg.App.Port, []string{"router", "api"},
		); err != nil {
			logger.Warn("consul register failed", logger.String("error", err.Error()))
		} else {
			defer consulClient.Deregister(serviceID)
			logger.Info("consul registered")
		}
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(middleware.CORS())
	r.Use(middleware.StructuredLog())
	r.Use(middleware.Recovery())

	if redisClient != nil {
		r.Use(middleware.IPRateLimiter(redisClient, 100, time.Minute))
	}

	_ = dbEngine
	_ = meiliClient
	_ = minioClient

	api := r.Group("/api/v1")
	{
		handler.RegisterSkillRoutes(api)
		handler.RegisterSearchRoutes(api)
		handler.RegisterRouterRoutes(api)
		handler.RegisterAuthRoutes(api)
		handler.RegisterUserRoutes(api)
	}

	router.SetupRoutes(r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: r,
	}

	go func() {
		logger.Info("HTTP server listening", logger.Int("port", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", logger.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	logger.Sync()
}
