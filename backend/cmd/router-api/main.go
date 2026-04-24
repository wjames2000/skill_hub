package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/internal/client/embedding"
	"github.com/hpds/skill-hub/internal/client/llm"
	"github.com/hpds/skill-hub/internal/client/reranker"
	"github.com/hpds/skill-hub/internal/handler"
	"github.com/hpds/skill-hub/internal/middleware"
	milvusCli "github.com/hpds/skill-hub/internal/milvus"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/internal/router"
	"github.com/hpds/skill-hub/internal/service"
	"github.com/hpds/skill-hub/internal/vectorizer"
	"github.com/hpds/skill-hub/pkg/config"
	"github.com/hpds/skill-hub/pkg/consul"
	"github.com/hpds/skill-hub/pkg/db"
	"github.com/hpds/skill-hub/pkg/logger"
	mls "github.com/hpds/skill-hub/pkg/meilisearch"
	rds "github.com/hpds/skill-hub/pkg/redis"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
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

	embedder := embedding.New(cfg.Embedding.BaseURL, cfg.Embedding.APIKey, cfg.Embedding.Model, cfg.Embedding.Dims)
	logger.Info("embedding client initialized", logger.String("model", embedder.Model()))

	var llmClient *llm.Client
	if cfg.LLM.BaseURL != "" {
		llmClient = llm.New(cfg.LLM.BaseURL, cfg.LLM.APIKey, cfg.LLM.Model, cfg.LLM.MaxTokens, cfg.LLM.Temperature)
		logger.Info("llm client initialized", logger.String("model", llmClient.Model()))
	}

	var rerankerClient *reranker.Client
	if cfg.Reranker.BaseURL != "" {
		rerankerClient = reranker.New(cfg.Reranker.BaseURL, cfg.Reranker.APIKey, cfg.Reranker.Model)
		logger.Info("reranker client initialized", logger.String("model", rerankerClient.Model()))
	}

	milvusCfg := milvusCli.Config{
		Address:  cfg.Milvus.Address,
		User:     cfg.Milvus.User,
		Password: cfg.Milvus.Password,
		DBName:   cfg.Milvus.DBName,
	}
	milvusClient, err := milvusCli.New(milvusCfg)
	if err != nil {
		logger.Warn("milvus unavailable, vector search disabled", logger.String("error", err.Error()))
		milvusClient = nil
	} else {
		logger.Info("milvus client initialized")
	}

	skillRepo := repository.NewSkillRepo(dbEngine)
	embRepo := repository.NewEmbeddingRepo(dbEngine)
	logRepo := repository.NewRouterLogRepo(dbEngine)
	userRepo := repository.NewUserRepo(dbEngine)
	apiKeyRepo := repository.NewAPIKeyRepo(dbEngine)
	favoriteRepo := repository.NewFavoriteRepo(dbEngine)
	reviewRepo := repository.NewReviewRepo(dbEngine)
	categoryRepo := repository.NewCategoryRepo(dbEngine)

	vectorWorker := vectorizer.NewWorker(embedder, milvusClient, skillRepo, embRepo, 3)

	vectorSvc := service.NewVectorService(
		embedder, llmClient, milvusClient, skillRepo, embRepo, vectorWorker,
	)

	routerSvc := service.NewRouterService(
		embedder, llmClient, rerankerClient, milvusClient, meiliClient, skillRepo, logRepo,
	)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(middleware.CORS())
	r.Use(middleware.StructuredLog())
	r.Use(middleware.Recovery())

	if redisClient != nil {
		r.Use(middleware.IPRateLimiter(redisClient, 100, time.Minute))
	}

	r.Use(middleware.OptionalAuth(cfg.JWT.Secret))

	api := r.Group("/api/v1")
	{
		skillHandler := handler.NewSkillHandler(skillRepo, categoryRepo, meiliClient)
		skillHandler.RegisterRoutes(api)

		routerHandler := handler.NewRouterHandler(routerSvc)
		routerHandler.RegisterRoutes(api)

		authHandler := handler.NewAuthHandler(userRepo, apiKeyRepo, cfg.JWT.Secret, cfg.JWT.ExpireHour)
		authHandler.RegisterRoutes(api)

		auth := api.Group("")
		auth.Use(middleware.AuthRequired(cfg.JWT.Secret))
		{
			userHandler := handler.NewUserHandler(userRepo, favoriteRepo, reviewRepo, apiKeyRepo, skillRepo)
			userHandler.RegisterRoutes(auth)

			categoryHandler := handler.NewCategoryHandler(categoryRepo)
			categoryHandler.RegisterRoutes(auth)

			statsHandler := handler.NewStatsHandler(skillRepo, categoryRepo, favoriteRepo)
			statsHandler.RegisterRoutes(auth)
		}
	}

	pluginAPI := r.Group("/api/v1")
	pluginAPI.Use(middleware.APIKeyAuth(apiKeyRepo))
	{
		pluginHandler := handler.NewPluginHandler(skillRepo, categoryRepo, meiliClient, favoriteRepo, reviewRepo)
		pluginHandler.RegisterRoutes(pluginAPI)
	}

	router.SetupRoutes(r)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vectorSvc.StartWorker(ctx)

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
	vectorSvc.StopWorker()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	logger.Sync()
}
