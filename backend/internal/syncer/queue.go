package syncer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hpds/skill-hub/pkg/logger"
)

const (
	TypeFullSync        = "sync:full"
	TypeIncrementalSync = "sync:incremental"
	TypeScanSkill       = "sync:scan"
	TypeSyncSingle      = "sync:single"

	QueueSync    = "sync"
	QueueScan    = "scan"
	QueueDefault = "default"
)

type SyncTaskPayload struct {
	TaskID   int64  `json:"task_id"`
	SyncType string `json:"sync_type"`
	Strategy string `json:"strategy"`
}

type ScanTaskPayload struct {
	SkillID  int64  `json:"skill_id"`
	RepoPath string `json:"repo_path"`
	FullName string `json:"full_name"`
}

type SingleSyncPayload struct {
	Owner    string `json:"owner"`
	Repo     string `json:"repo"`
	FullName string `json:"full_name"`
}

type TaskQueue struct {
	client  *asynq.Client
	enabled bool
}

func NewTaskQueue(redisAddr, redisPassword string, db int, enabled bool) (*TaskQueue, error) {
	if !enabled {
		return &TaskQueue{enabled: false}, nil
	}

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       db,
	})

	return &TaskQueue{
		client:  client,
		enabled: true,
	}, nil
}

func (q *TaskQueue) Close() error {
	if q.client != nil {
		return q.client.Close()
	}
	return nil
}

func (q *TaskQueue) EnqueueFullSync(taskID int64, strategy string) error {
	if !q.enabled {
		logger.Warn("queue disabled, skipping enqueue")
		return nil
	}

	payload := SyncTaskPayload{
		TaskID:   taskID,
		SyncType: "full",
		Strategy: strategy,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeFullSync, data,
		asynq.MaxRetry(3),
		asynq.Timeout(30*time.Minute),
		asynq.Queue(QueueSync),
	)

	info, err := q.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("enqueue: %w", err)
	}

	logger.Info("enqueued full sync task",
		logger.String("id", info.ID),
		logger.Int64("task_id", taskID))
	return nil
}

func (q *TaskQueue) EnqueueIncrementalSync(taskID int64) error {
	if !q.enabled {
		return nil
	}

	payload := SyncTaskPayload{
		TaskID:   taskID,
		SyncType: "incremental",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeIncrementalSync, data,
		asynq.MaxRetry(3),
		asynq.Timeout(30*time.Minute),
		asynq.Queue(QueueSync),
	)

	info, err := q.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("enqueue: %w", err)
	}

	logger.Info("enqueued incremental sync task",
		logger.String("id", info.ID),
		logger.Int64("task_id", taskID))
	return nil
}

func (q *TaskQueue) EnqueueScan(skillID int64, repoPath, fullName string) error {
	if !q.enabled {
		return nil
	}

	payload := ScanTaskPayload{
		SkillID:  skillID,
		RepoPath: repoPath,
		FullName: fullName,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	task := asynq.NewTask(TypeScanSkill, data,
		asynq.MaxRetry(2),
		asynq.Timeout(10*time.Minute),
		asynq.Queue(QueueScan),
	)

	info, err := q.client.Enqueue(task)
	if err != nil {
		return fmt.Errorf("enqueue scan: %w", err)
	}

	logger.Info("enqueued scan task", logger.String("id", info.ID), logger.Int64("skill_id", skillID))
	return nil
}

func NewSyncServer(redisAddr, redisPassword string, db int, concurrency int) (*asynq.Server, error) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       db,
		},
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				QueueSync:    6,
				QueueScan:    3,
				QueueDefault: 1,
			},
			StrictPriority: false,
		},
	)
	return srv, nil
}

func HandleSyncTask(ctx context.Context, t *asynq.Task, handler func(context.Context, SyncTaskPayload) error) error {
	var payload SyncTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	logger.Info("processing sync task",
		logger.String("type", t.Type()),
		logger.Int64("task_id", payload.TaskID))

	return handler(ctx, payload)
}

func HandleScanTask(ctx context.Context, t *asynq.Task, handler func(context.Context, ScanTaskPayload) error) error {
	var payload ScanTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal scan payload: %w", err)
	}

	logger.Info("processing scan task",
		logger.Int64("skill_id", payload.SkillID),
		logger.String("repo", payload.FullName))

	return handler(ctx, payload)
}
