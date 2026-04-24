package service

import (
	"context"
	"fmt"
	"time"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/internal/syncer"
	"github.com/hpds/skill-hub/pkg/logger"
)

type SyncService struct {
	skillRepo    *repository.SkillRepo
	syncTaskRepo *repository.SyncTaskRepo
	orchestrator *syncer.SyncOrchestrator
	scheduler    *syncer.Scheduler
	queue        *syncer.TaskQueue
	syncConfig   syncer.SyncConfig
}

func NewSyncService(
	skillRepo *repository.SkillRepo,
	syncTaskRepo *repository.SyncTaskRepo,
	orchestrator *syncer.SyncOrchestrator,
	scheduler *syncer.Scheduler,
	queue *syncer.TaskQueue,
	syncConfig syncer.SyncConfig,
) *SyncService {
	return &SyncService{
		skillRepo:    skillRepo,
		syncTaskRepo: syncTaskRepo,
		orchestrator: orchestrator,
		scheduler:    scheduler,
		queue:        queue,
		syncConfig:   syncConfig,
	}
}

func (s *SyncService) TriggerFullSync(ctx context.Context, strategy string) (*model.SyncTask, error) {
	return syncer.CreateSyncTask(ctx, s.syncTaskRepo, s.queue, model.SyncTypeFull, strategy)
}

func (s *SyncService) TriggerIncrementalSync(ctx context.Context) (*model.SyncTask, error) {
	return syncer.CreateSyncTask(ctx, s.syncTaskRepo, s.queue, model.SyncTypeIncremental, "")
}

func (s *SyncService) GetSyncStatus(ctx context.Context, taskID int64) (*model.SyncTask, error) {
	task, err := s.syncTaskRepo.GetByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	return task, nil
}

func (s *SyncService) ListSyncTasks(ctx context.Context, page, pageSize int) ([]*model.SyncTask, int64, error) {
	return s.syncTaskRepo.List(page, pageSize)
}

func (s *SyncService) GetLatestSyncTask(ctx context.Context, syncType string) (*model.SyncTask, error) {
	return s.syncTaskRepo.GetLatestByType(syncType)
}

func (s *SyncService) GetRunningTask(ctx context.Context) (*model.SyncTask, error) {
	return s.syncTaskRepo.GetRunningTask()
}

func (s *SyncService) GetSkillStats(ctx context.Context) (*repository.SkillStats, error) {
	return s.skillRepo.GetStats()
}

func (s *SyncService) SetScheduler(scheduler *syncer.Scheduler) {
	s.scheduler = scheduler
}

func (s *SyncService) SyncSingleRepo(ctx context.Context, owner, repo string) error {
	fullName := fmt.Sprintf("%s/%s", owner, repo)
	return s.orchestrator.SyncSingleRepo(ctx, owner, repo, fullName)
}

func (s *SyncService) HandleSyncTask(ctx context.Context, payload syncer.SyncTaskPayload) error {
	switch payload.SyncType {
	case "full":
		return s.orchestrator.ExecuteFullSync(ctx, payload.TaskID, payload.Strategy)
	case "incremental":
		return s.orchestrator.ExecuteIncrementalSync(ctx, payload.TaskID)
	default:
		return fmt.Errorf("unknown sync type: %s", payload.SyncType)
	}
}

func (s *SyncService) HandleScanTask(ctx context.Context, payload syncer.ScanTaskPayload) error {
	if !s.syncConfig.ScanEnabled {
		return nil
	}

	start := time.Now()
	logger.Info("starting security scan",
		logger.Int64("skill_id", payload.SkillID),
		logger.String("repo", payload.FullName))

	result, err := s.orchestrator.Scanner().ScanRepo(ctx, payload.RepoPath)
	if err != nil {
		return fmt.Errorf("scan repo: %w", err)
	}

	if err := s.skillRepo.UpdateScanResult(payload.SkillID, result.Passed, result.Summary); err != nil {
		return fmt.Errorf("update scan result: %w", err)
	}

	logger.Info("scan completed",
		logger.Int64("skill_id", payload.SkillID),
		logger.Bool("passed", result.Passed),
		logger.Duration("duration", time.Since(start)))

	return nil
}

func (s *SyncService) StartScheduler() error {
	return s.scheduler.Start()
}

func (s *SyncService) StopScheduler() {
	s.scheduler.Stop()
}
