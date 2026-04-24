package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hpds/skill-hub/pkg/config"
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
	logger.Info("sync-worker starting", logger.String("app", cfg.App.Name))

	dbEngine, err := db.WaitForDB(cfg.DB, 30*time.Second)
	if err != nil {
		logger.Fatal("db init", logger.String("error", err.Error()))
	}
	defer dbEngine.Close()
	logger.Info("database connected")

	redisClient, err := rds.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Warn("redis unavailable", logger.String("error", err.Error()))
	} else {
		defer redisClient.Close()
		logger.Info("redis connected")
	}

	_ = dbEngine
	_ = redisClient

	logger.Info("sync-worker ready, waiting for tasks...")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("sync-worker shutting down...")
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.Sync()
}
