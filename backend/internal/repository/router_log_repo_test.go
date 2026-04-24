package repository

import (
	"context"
	"testing"
	"time"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouterLogRepo_Create(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.RouterLog{})
	repo := NewRouterLogRepo(engine)
	ctx := context.Background()

	log := &model.RouterLog{
		SessionID:    "sess_001",
		Query:        "analyze data",
		MatchedSkill: []byte(`[{"skill_id":1,"score":0.95}]`),
		Strategy:     "hybrid",
		Duration:     150,
		Status:       1,
	}

	err := repo.Create(ctx, log)
	require.NoError(t, err)
	assert.Greater(t, log.ID, int64(0))
}

func TestRouterLogRepo_List(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.RouterLog{})
	repo := NewRouterLogRepo(engine)
	ctx := context.Background()

	now := time.Now()
	logs := []*model.RouterLog{
		{SessionID: "sess_1", Query: "q1", Strategy: "hybrid", Duration: 100, Status: 1, CreatedAt: &now},
		{SessionID: "sess_2", Query: "q2", Strategy: "vector", Duration: 200, Status: 1, CreatedAt: &now},
	}
	for _, l := range logs {
		_, _ = engine.Insert(l)
	}

	t.Run("list all logs", func(t *testing.T) {
		results, total, err := repo.List(ctx, "", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, 2, len(results))
	})

	t.Run("list with session filter", func(t *testing.T) {
		results, total, err := repo.List(ctx, "sess_1", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "q1", results[0].Query)
	})
}

func TestRouterLogRepo_GetByID(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.RouterLog{})
	repo := NewRouterLogRepo(engine)
	ctx := context.Background()

	now := time.Now()
	_, _ = engine.Insert(&model.RouterLog{
		SessionID: "sess_get", Query: "find me", Strategy: "keyword",
		Duration: 50, Status: 1, CreatedAt: &now,
	})

	log, err := repo.GetByID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, "find me", log.Query)
}

func TestRouterLogRepo_UpdateFeedback(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.RouterLog{})
	repo := NewRouterLogRepo(engine)
	ctx := context.Background()

	now := time.Now()
	_, _ = engine.Insert(&model.RouterLog{
		SessionID: "sess_fb", Query: "feedback", Strategy: "hybrid",
		Duration: 100, Status: 1, CreatedAt: &now,
	})

	err := repo.UpdateFeedback(ctx, 1, 4)
	require.NoError(t, err)

	log, _ := repo.GetByID(ctx, 1)
	assert.Equal(t, 4, log.Feedback)
}

func TestRouterLogRepo_Delete(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.RouterLog{})
	repo := NewRouterLogRepo(engine)
	ctx := context.Background()

	now := time.Now()
	_, _ = engine.Insert(&model.RouterLog{
		SessionID: "sess_del", Query: "delete me", Strategy: "keyword",
		Duration: 50, Status: 1, CreatedAt: &now,
	})

	err := repo.Delete(ctx, 1)
	require.NoError(t, err)

	deleted, err := repo.GetByID(ctx, 1)
	assert.Error(t, err)
	assert.Nil(t, deleted)
}
