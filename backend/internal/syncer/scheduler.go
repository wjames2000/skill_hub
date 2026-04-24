package syncer

import (
	"context"
	"fmt"
	"time"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/pkg/logger"

	gocron "github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron         *gocron.Cron
	syncTaskRepo *repository.SyncTaskRepo
	queue        *TaskQueue
	config       SyncConfig
	entries      map[string]gocron.EntryID
}

func NewScheduler(syncTaskRepo *repository.SyncTaskRepo, queue *TaskQueue, config SyncConfig) *Scheduler {
	return &Scheduler{
		cron:         gocron.New(gocron.WithSeconds(), gocron.WithLogger(&cronLogger{})),
		syncTaskRepo: syncTaskRepo,
		queue:        queue,
		config:       config,
		entries:      make(map[string]gocron.EntryID),
	}
}

func (s *Scheduler) Start() error {
	if s.config.FullSyncCron != "" {
		id, err := s.cron.AddFunc(s.config.FullSyncCron, func() {
			s.triggerSync("full", "topic")
		})
		if err != nil {
			return fmt.Errorf("register full sync cron: %w", err)
		}
		s.entries["full"] = id
		logger.Info("registered full sync cron", logger.String("schedule", s.config.FullSyncCron))
	}

	if s.config.IncrementalCron != "" {
		id, err := s.cron.AddFunc(s.config.IncrementalCron, func() {
			s.triggerSync("incremental", "")
		})
		if err != nil {
			return fmt.Errorf("register incremental sync cron: %w", err)
		}
		s.entries["incremental"] = id
		logger.Info("registered incremental sync cron", logger.String("schedule", s.config.IncrementalCron))
	}

	s.cron.Start()
	logger.Info("scheduler started")
	return nil
}

func (s *Scheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
	}
	logger.Info("scheduler stopped")
}

func (s *Scheduler) triggerSync(syncType, strategy string) {
	logger.Info("cron triggering sync", logger.String("type", syncType))

	task := &model.SyncTask{
		Type:     syncType,
		Strategy: strategy,
		Status:   model.SyncStatusPending,
	}

	if err := s.syncTaskRepo.Create(task); err != nil {
		logger.Error("create sync task failed", logger.String("error", err.Error()))
		return
	}

	var err error
	if syncType == "full" {
		err = s.queue.EnqueueFullSync(task.ID, strategy)
	} else {
		err = s.queue.EnqueueIncrementalSync(task.ID)
	}

	if err != nil {
		logger.Error("enqueue sync task failed",
			logger.Int64("task_id", task.ID),
			logger.String("error", err.Error()))
	}
}

type cronLogger struct{}

func (l *cronLogger) Info(msg string, keysAndValues ...interface{}) {
	logger.Info("[cron] "+msg, logger.Any("args", keysAndValues))
}

func (l *cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	logger.Error("[cron] "+msg,
		logger.String("error", err.Error()),
		logger.Any("args", keysAndValues))
}

func CreateSyncTask(ctx context.Context, syncTaskRepo *repository.SyncTaskRepo, queue *TaskQueue, syncType, strategy string) (*model.SyncTask, error) {
	running, err := syncTaskRepo.GetRunningTask()
	if err != nil {
		return nil, fmt.Errorf("check running: %w", err)
	}
	if running != nil {
		return nil, fmt.Errorf("a sync task is already running (id=%d)", running.ID)
	}

	task := &model.SyncTask{
		Type:     syncType,
		Strategy: strategy,
		Status:   model.SyncStatusPending,
	}

	if err := syncTaskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	if syncType == "full" {
		err = queue.EnqueueFullSync(task.ID, strategy)
	} else {
		err = queue.EnqueueIncrementalSync(task.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("enqueue: %w", err)
	}

	return task, nil
}

func init() {
	// Ensure timezone is available
	_, _ = time.LoadLocation("UTC")
}
