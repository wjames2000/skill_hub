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
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/internal/service"
	"github.com/hpds/skill-hub/internal/syncer"
	"github.com/hpds/skill-hub/pkg/config"
	"github.com/hpds/skill-hub/pkg/consul"
	"github.com/hpds/skill-hub/pkg/db"
	"github.com/hpds/skill-hub/pkg/logger"
	rds "github.com/hpds/skill-hub/pkg/redis"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		panic("load config: " + err.Error())
	}

	logger.Init(cfg.App.Env)
	logger.Info("admin-api starting", logger.String("app", cfg.App.Name), logger.String("env", cfg.App.Env))

	dbEngine, err := db.WaitForDB(cfg.DB, 30*time.Second)
	if err != nil {
		logger.Fatal("db init", logger.String("error", err.Error()))
	}
	defer dbEngine.Close()
	logger.Info("database connected")

	redisClient, err := rds.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Warn("redis unavailable", logger.String("error", err.Error()))
		redisClient = nil
	} else {
		defer redisClient.Close()
		logger.Info("redis connected")
	}

	consulClient, err := consul.New(cfg.Consul.Addr, cfg.Consul.Token)
	if err != nil {
		logger.Warn("consul unavailable", logger.String("error", err.Error()))
		consulClient = nil
	} else {
		serviceID := fmt.Sprintf("admin-api-%d", cfg.App.Port)
		if err := consulClient.RegisterService(
			cfg.App.Name+"-admin", serviceID, "localhost", cfg.App.Port, []string{"admin", "api"},
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

	skillRepo := repository.NewSkillRepo(dbEngine)
	syncTaskRepo := repository.NewSyncTaskRepo(dbEngine)
	userRepo := repository.NewUserRepo(dbEngine)
	categoryRepo := repository.NewCategoryRepo(dbEngine)
	favoriteRepo := repository.NewFavoriteRepo(dbEngine)
	reviewRepo := repository.NewReviewRepo(dbEngine)

	syncSvc := service.NewSyncService(
		skillRepo,
		syncTaskRepo,
		nil,
		nil,
		nil,
		syncer.SyncConfig{},
	)

	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.AdminRequired(cfg.JWT.Secret))
	{
		syncAdmin := handler.NewSyncAdminHandler(syncSvc)
		syncAdmin.RegisterRoutes(admin)

		adminHandler := handler.NewAdminHandler(skillRepo, categoryRepo, userRepo, favoriteRepo, reviewRepo)
		adminHandler.RegisterRoutes(admin)
	}

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
