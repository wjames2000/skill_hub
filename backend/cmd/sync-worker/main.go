package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	githubclient "github.com/hpds/skill-hub/internal/client/github"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/internal/service"
	"github.com/hpds/skill-hub/internal/syncer"
	"github.com/hpds/skill-hub/pkg/config"
	"github.com/hpds/skill-hub/pkg/db"
	"github.com/hpds/skill-hub/pkg/logger"
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

	syncCfg := syncer.SyncConfig{
		FullSyncCron:    cfg.Sync.FullSyncCron,
		IncrementalCron: cfg.Sync.IncrementalCron,
		IncrementalDays: cfg.Sync.IncrementalDays,
		Concurrency:     cfg.Sync.Concurrency,
		SyncTimeout:     cfg.Sync.SyncTimeout,
		ScanEnabled:     cfg.Sync.ScanEnabled,
	}

	ghClient := githubclient.New(
		cfg.GitHub.Tokens,
		cfg.GitHub.MaxPerPage,
		cfg.GitHub.RequestDelay,
	)
	logger.Info("github client initialized", logger.Int("tokens", len(cfg.GitHub.Tokens)))

	skillRepo := repository.NewSkillRepo(dbEngine)
	syncTaskRepo := repository.NewSyncTaskRepo(dbEngine)

	scannerCfg := syncer.ScannerConfig{
		Enabled: cfg.Semgrep.Enabled,
		Binary:  cfg.Semgrep.Binary,
		Rules:   cfg.Semgrep.Rules,
		Timeout: cfg.Semgrep.Timeout,
	}
	scanner := syncer.NewSecurityScanner(scannerCfg)

	topicDiscovery := syncer.NewTopicDiscovery(ghClient, cfg.GitHub.SearchTopics)
	awesomeDiscovery := syncer.NewAwesomeDiscovery(ghClient)
	discoveryMgr := syncer.NewDiscoveryManager(topicDiscovery, awesomeDiscovery)

	queue, err := syncer.NewTaskQueue(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, cfg.Asynq.Enabled)
	if err != nil {
		logger.Warn("task queue init failed", logger.String("error", err.Error()))
		queue, _ = syncer.NewTaskQueue(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, false)
	}
	if queue != nil {
		defer queue.Close()
	}

	orchestrator := syncer.NewSyncOrchestrator(
		ghClient,
		discoveryMgr,
		skillRepo,
		syncTaskRepo,
		scanner,
		queue,
		syncCfg,
	)

	syncSvc := service.NewSyncService(
		skillRepo,
		syncTaskRepo,
		orchestrator,
		nil,
		queue,
		syncCfg,
	)

	scheduler := syncer.NewScheduler(syncTaskRepo, queue, syncCfg)
	syncSvc.SetScheduler(scheduler)

	if err := scheduler.Start(); err != nil {
		logger.Warn("scheduler start failed", logger.String("error", err.Error()))
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg.Asynq.Enabled {
		go func() {
			asynqSrv, err := syncer.NewSyncServer(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, cfg.Asynq.Concurrency)
			if err != nil {
				logger.Fatal("asynq server init", logger.String("error", err.Error()))
			}

			mux := asynq.NewServeMux()
			mux.HandleFunc(syncer.TypeFullSync, func(ctx context.Context, t *asynq.Task) error {
				return syncer.HandleSyncTask(ctx, t, func(ctx context.Context, p syncer.SyncTaskPayload) error {
					return syncSvc.HandleSyncTask(ctx, p)
				})
			})
			mux.HandleFunc(syncer.TypeIncrementalSync, func(ctx context.Context, t *asynq.Task) error {
				return syncer.HandleSyncTask(ctx, t, func(ctx context.Context, p syncer.SyncTaskPayload) error {
					return syncSvc.HandleSyncTask(ctx, p)
				})
			})
			mux.HandleFunc(syncer.TypeScanSkill, func(ctx context.Context, t *asynq.Task) error {
				return syncer.HandleScanTask(ctx, t, func(ctx context.Context, p syncer.ScanTaskPayload) error {
					return syncSvc.HandleScanTask(ctx, p)
				})
			})

			logger.Info("asynq server starting",
				logger.Int("concurrency", cfg.Asynq.Concurrency))

			if err := asynqSrv.Run(mux); err != nil {
				logger.Error("asynq server error", logger.String("error", err.Error()))
			}
		}()
	}

	logger.Info(fmt.Sprintf("sync-worker ready | scheduler: %s/%s | queue: %v | scan: %v",
		syncCfg.FullSyncCron, syncCfg.IncrementalCron, cfg.Asynq.Enabled, cfg.Semgrep.Enabled))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("sync-worker shutting down...")
	scheduler.Stop()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = shutdownCtx
	logger.Sync()
}
