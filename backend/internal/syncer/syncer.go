package syncer

import (
	"context"
	"fmt"
	"time"

	githubclient "github.com/hpds/skill-hub/internal/client/github"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/pkg/logger"
)

type SyncConfig struct {
	FullSyncCron    string
	IncrementalCron string
	IncrementalDays int
	Concurrency     int
	SyncTimeout     int
	ScanEnabled     bool
}

type SyncOrchestrator struct {
	githubClient *githubclient.Client
	discovery    *DiscoveryManager
	parser       *SkillMetadataParser
	extractor    *SkillExtractor
	skillRepo    *repository.SkillRepo
	syncTaskRepo *repository.SyncTaskRepo
	scanner      *SecurityScanner
	queue        *TaskQueue
	config       SyncConfig
}

type SkillMetadataParser struct{}

func NewSkillMetadataParser() *SkillMetadataParser {
	return &SkillMetadataParser{}
}

func (p *SkillMetadataParser) Parse(content string) (*ParsedSkill, error) {
	return ParseSKILLMD(content)
}

func NewSyncOrchestrator(
	githubClient *githubclient.Client,
	discovery *DiscoveryManager,
	skillRepo *repository.SkillRepo,
	syncTaskRepo *repository.SyncTaskRepo,
	scanner *SecurityScanner,
	queue *TaskQueue,
	config SyncConfig,
) *SyncOrchestrator {
	return &SyncOrchestrator{
		githubClient: githubClient,
		discovery:    discovery,
		parser:       NewSkillMetadataParser(),
		extractor:    NewSkillExtractor(),
		skillRepo:    skillRepo,
		syncTaskRepo: syncTaskRepo,
		scanner:      scanner,
		queue:        queue,
		config:       config,
	}
}

func (s *SyncOrchestrator) ExecuteFullSync(ctx context.Context, taskID int64, strategy string) error {
	task, err := s.syncTaskRepo.GetByID(taskID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	now := time.Now()
	task.Status = model.SyncStatusRunning
	task.StartedAt = &now
	_ = s.syncTaskRepo.Update(task)

	logger.Info("starting full sync", logger.Int64("task_id", taskID), logger.String("strategy", strategy))

	var since time.Time
	if s.config.IncrementalDays > 0 {
		since = time.Now().AddDate(0, 0, -s.config.IncrementalDays)
	}

	taskCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.SyncTimeout)*time.Second)
	defer cancel()

	repos, err := s.discovery.DiscoverAll(taskCtx, since)
	if err != nil {
		task.Status = model.SyncStatusFailed
		task.ErrorMessage = fmt.Sprintf("discovery failed: %v", err)
		finish := time.Now()
		task.FinishedAt = &finish
		_ = s.syncTaskRepo.Update(task)
		return fmt.Errorf("discover: %w", err)
	}

	task.FoundRepos = len(repos)
	_ = s.syncTaskRepo.Update(task)

	logger.Info("discovered repos", logger.Int("count", len(repos)))

	for i, repo := range repos {
		select {
		case <-taskCtx.Done():
			task.Status = model.SyncStatusFailed
			task.ErrorMessage = "sync cancelled: " + taskCtx.Err().Error()
			finish := time.Now()
			task.FinishedAt = &finish
			_ = s.syncTaskRepo.Update(task)
			return taskCtx.Err()
		default:
		}

		logger.Info("processing repo",
			logger.String("repo", repo.FullName),
			logger.Int("progress", i+1),
			logger.Int("total", len(repos)))

		if err := s.processRepo(taskCtx, repo, task); err != nil {
			logger.Warn("failed to process repo",
				logger.String("repo", repo.FullName),
				logger.String("error", err.Error()))
			task.FailedSkills++
			_ = s.syncTaskRepo.Update(task)
		}
	}

	task.TotalRepos = len(repos)
	task.Status = model.SyncStatusCompleted
	finish := time.Now()
	task.FinishedAt = &finish
	_ = s.syncTaskRepo.Update(task)

	logger.Info("full sync completed",
		logger.Int64("task_id", taskID),
		logger.Int("total_repos", task.TotalRepos),
		logger.Int("new_skills", task.NewSkills),
		logger.Int("updated_skills", task.UpdatedSkills),
		logger.Int("failed", task.FailedSkills))

	return nil
}

func (s *SyncOrchestrator) ExecuteIncrementalSync(ctx context.Context, taskID int64) error {
	since := time.Now().AddDate(0, 0, -s.config.IncrementalDays)

	task, err := s.syncTaskRepo.GetByID(taskID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	now := time.Now()
	task.Status = model.SyncStatusRunning
	task.StartedAt = &now
	_ = s.syncTaskRepo.Update(task)

	taskCtx, cancel := context.WithTimeout(ctx, time.Duration(s.config.SyncTimeout)*time.Second)
	defer cancel()

	ids, err := s.skillRepo.ListIDsNeedingUpdate(since, 500)
	if err != nil {
		task.Status = model.SyncStatusFailed
		task.ErrorMessage = fmt.Sprintf("list stale skills: %v", err)
		finish := time.Now()
		task.FinishedAt = &finish
		_ = s.syncTaskRepo.Update(task)
		return fmt.Errorf("list stale: %w", err)
	}

	task.FoundRepos = len(ids)
	_ = s.syncTaskRepo.Update(task)

	logger.Info("incremental sync found stale skills", logger.Int("count", len(ids)))

	for _, skillID := range ids {
		select {
		case <-taskCtx.Done():
			task.Status = model.SyncStatusFailed
			task.ErrorMessage = "sync cancelled"
			finish := time.Now()
			task.FinishedAt = &finish
			_ = s.syncTaskRepo.Update(task)
			return taskCtx.Err()
		default:
		}

		skill, err := s.skillRepo.GetByID(skillID)
		if err != nil || skill == nil {
			continue
		}

		repoInfo, err := s.githubClient.GetRepo(taskCtx, skill.RepoOwner, skill.RepoName)
		if err != nil {
			logger.Warn("incremental: failed to get repo info",
				logger.String("repo", skill.Repository),
				logger.String("error", err.Error()))
			continue
		}

		var parsed *ParsedSkill
		content, err := s.githubClient.GetRepoContent(taskCtx, skill.RepoOwner, skill.RepoName, "SKILL.md", skill.DefaultBranch)
		if err == nil && content != nil {
			parsed, _ = s.parser.Parse(content.Content)
		}

		readme, _ := s.githubClient.GetReadme(taskCtx, skill.RepoOwner, skill.RepoName, skill.DefaultBranch)

		updatedSkill := s.extractor.Extract(repoInfo, parsed, readme)
		isNew, err := s.skillRepo.Upsert(updatedSkill)
		if err != nil {
			logger.Warn("incremental: upsert failed",
				logger.String("repo", skill.Repository),
				logger.String("error", err.Error()))
			task.FailedSkills++
			continue
		}

		if isNew {
			task.NewSkills++
		} else {
			task.UpdatedSkills++
		}
		_ = s.syncTaskRepo.Update(task)
	}

	task.Status = model.SyncStatusCompleted
	finish := time.Now()
	task.FinishedAt = &finish
	_ = s.syncTaskRepo.Update(task)

	logger.Info("incremental sync completed",
		logger.Int("processed", len(ids)),
		logger.Int("new", task.NewSkills),
		logger.Int("updated", task.UpdatedSkills))

	return nil
}

func (s *SyncOrchestrator) processRepo(ctx context.Context, repo DiscoveredRepo, task *model.SyncTask) error {
	repoInfo, err := s.githubClient.GetRepo(ctx, repo.Owner, repo.Name)
	if err != nil {
		return fmt.Errorf("get repo info: %w", err)
	}

	if repoInfo.Archived {
		logger.Debug("skipping archived repo", logger.String("repo", repo.FullName))
		return nil
	}

	content, err := s.githubClient.GetRepoContent(ctx, repo.Owner, repo.Name, "SKILL.md", repoInfo.DefaultBranch)
	if err != nil {
		return fmt.Errorf("get skill.md: %w", err)
	}

	var parsed *ParsedSkill
	if content != nil {
		skillFile := content
		task.ParsedSkills++
		_ = s.syncTaskRepo.Update(task)

		parsed, err = s.parser.Parse(skillFile.Content)
		if err != nil {
			logger.Warn("parse skill.md failed",
				logger.String("repo", repo.FullName),
				logger.String("error", err.Error()))
		}

		if parsed != nil && parsed.Metadata != nil {
			skill := s.skillRepo
			_ = skill
		}
	}

	readme, _ := s.githubClient.GetReadme(ctx, repo.Owner, repo.Name, repoInfo.DefaultBranch)

	skillModel := s.extractor.Extract(repoInfo, parsed, readme)

	if parsed != nil && parsed.Metadata != nil {
		skillModel.SkillFileSHA = content.SHA
		skillModel.SkillPath = "SKILL.md"
	}

	isNew, err := s.skillRepo.Upsert(skillModel)
	if err != nil {
		return fmt.Errorf("save skill: %w", err)
	}

	if isNew {
		task.NewSkills++
	} else {
		task.UpdatedSkills++
	}
	_ = s.syncTaskRepo.Update(task)

	if s.config.ScanEnabled && s.scanner != nil {
		cloneDir := fmt.Sprintf("/tmp/skill-scan/%s", repo.FullName)
		_ = cloneDir

		if s.queue.enabled {
			_ = s.queue.EnqueueScan(skillModel.ID, cloneDir, repo.FullName)
		} else {
			s.runLocalScan(ctx, skillModel)
		}
	}

	return nil
}

func (s *SyncOrchestrator) Scanner() *SecurityScanner {
	return s.scanner
}

func (s *SyncOrchestrator) runLocalScan(ctx context.Context, skill *model.Skill) {
	result, err := s.scanner.ScanRepo(ctx, "")
	if err != nil {
		logger.Warn("scan failed", logger.String("repo", skill.Repository), logger.String("error", err.Error()))
		return
	}

	_ = s.skillRepo.UpdateScanResult(skill.ID, result.Passed, result.Summary)
}

func (s *SyncOrchestrator) SyncSingleRepo(ctx context.Context, owner, repo, fullName string) error {
	repoInfo, err := s.githubClient.GetRepo(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("get repo: %w", err)
	}

	content, err := s.githubClient.GetRepoContent(ctx, owner, repo, "SKILL.md", repoInfo.DefaultBranch)
	if err != nil {
		return fmt.Errorf("get skill.md: %w", err)
	}

	var parsed *ParsedSkill
	if content != nil {
		parsed, _ = s.parser.Parse(content.Content)
	}

	readme, _ := s.githubClient.GetReadme(ctx, owner, repo, repoInfo.DefaultBranch)

	skillModel := s.extractor.Extract(repoInfo, parsed, readme)

	_, err = s.skillRepo.Upsert(skillModel)
	if err != nil {
		return fmt.Errorf("save skill: %w", err)
	}

	logger.Info("synced single repo", logger.String("repo", fullName))
	return nil
}
