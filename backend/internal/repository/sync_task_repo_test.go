package repository

import (
	"testing"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncTaskRepo_Create(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SyncTask{})
	repo := NewSyncTaskRepo(engine)

	task := &model.SyncTask{
		Type:     model.SyncTypeFull,
		Strategy: model.StrategyTopic,
		Status:   model.SyncStatusPending,
	}

	err := repo.Create(task)
	require.NoError(t, err)
	assert.Greater(t, task.ID, int64(0))
}

func TestSyncTaskRepo_GetByID(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SyncTask{})
	repo := NewSyncTaskRepo(engine)

	_, _ = engine.Insert(&model.SyncTask{Type: model.SyncTypeFull, Status: model.SyncStatusPending})

	task, err := repo.GetByID(1)
	require.NoError(t, err)
	require.NotNil(t, task)
	assert.Equal(t, model.SyncTypeFull, task.Type)

	notFound, err := repo.GetByID(999)
	require.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestSyncTaskRepo_Update(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SyncTask{})
	repo := NewSyncTaskRepo(engine)

	_, _ = engine.Insert(&model.SyncTask{Type: model.SyncTypeFull, Status: model.SyncStatusPending})

	task, _ := repo.GetByID(1)
	task.Status = model.SyncStatusRunning
	err := repo.Update(task)
	require.NoError(t, err)

	updated, _ := repo.GetByID(1)
	assert.Equal(t, model.SyncStatusRunning, updated.Status)
}

func TestSyncTaskRepo_List(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SyncTask{})
	repo := NewSyncTaskRepo(engine)

	tasks := []*model.SyncTask{
		{Type: model.SyncTypeFull, Status: model.SyncStatusCompleted},
		{Type: model.SyncTypeIncremental, Status: model.SyncStatusCompleted},
		{Type: model.SyncTypeFull, Status: model.SyncStatusFailed},
	}
	for _, tsk := range tasks {
		_, _ = engine.Insert(tsk)
	}

	results, total, err := repo.List(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Equal(t, 3, len(results))
}

func TestSyncTaskRepo_GetLatestByType(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SyncTask{})
	repo := NewSyncTaskRepo(engine)

	_, _ = engine.Insert(&model.SyncTask{Type: model.SyncTypeFull, Status: model.SyncStatusCompleted})
	_, _ = engine.Insert(&model.SyncTask{Type: model.SyncTypeFull, Status: model.SyncStatusFailed})

	latest, err := repo.GetLatestByType(model.SyncTypeFull)
	require.NoError(t, err)
	require.NotNil(t, latest)
	assert.Equal(t, model.SyncStatusFailed, latest.Status)

	nilResult, err := repo.GetLatestByType("nonexistent")
	require.NoError(t, err)
	assert.Nil(t, nilResult)
}

func TestSyncTaskRepo_GetRunningTask(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SyncTask{})
	repo := NewSyncTaskRepo(engine)

	t.Run("no running task", func(t *testing.T) {
		task, err := repo.GetRunningTask()
		require.NoError(t, err)
		assert.Nil(t, task)
	})

	t.Run("finds running task", func(t *testing.T) {
		_, _ = engine.Insert(&model.SyncTask{Type: model.SyncTypeFull, Status: model.SyncStatusRunning})

		task, err := repo.GetRunningTask()
		require.NoError(t, err)
		require.NotNil(t, task)
		assert.Equal(t, model.SyncStatusRunning, task.Status)
	})
}

func TestSyncTaskRepo_CancelRunning(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SyncTask{})
	repo := NewSyncTaskRepo(engine)

	_, _ = engine.Insert(&model.SyncTask{Type: model.SyncTypeFull, Status: model.SyncStatusRunning})

	err := repo.CancelRunning()
	require.NoError(t, err)

	task, _ := repo.GetByID(1)
	assert.Equal(t, model.SyncStatusCancelled, task.Status)
	assert.Equal(t, "cancelled by new task", task.ErrorMessage)
}
