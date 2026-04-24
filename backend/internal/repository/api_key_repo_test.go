package repository

import (
	"context"
	"testing"
	"time"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIKeyRepo_Create(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.User{}, &model.APIKey{})
	repo := NewAPIKeyRepo(engine)
	ctx := context.Background()

	now := time.Now()
	key := &model.APIKey{
		UserID:    1,
		Key:       "sk_test_key_123",
		Name:      "Test Key",
		LastUsed:  &now,
		ExpiresAt: &now.Add(365 * 24 * time.Hour),
		Status:    1,
	}

	err := repo.Create(ctx, key)
	require.NoError(t, err)
	assert.Greater(t, key.ID, int64(0))

	// Create user first
	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	_, _ = engine.Insert(&model.User{Username: "keyuser", Email: "key@t.com", PasswordHash: string(hash), Role: "user", Status: 1})
}

func TestAPIKeyRepo_FindByKey(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.APIKey{})
	repo := NewAPIKeyRepo(engine)
	ctx := context.Background()

	now := time.Now()
	_, _ = engine.Insert(&model.APIKey{
		UserID: 1, Key: "sk_find_me", Name: "Find Key",
		ExpiresAt: &now.Add(24 * time.Hour), Status: 1,
	})

	t.Run("find existing key", func(t *testing.T) {
		key, err := repo.FindByKey(ctx, "sk_find_me")
		require.NoError(t, err)
		assert.Equal(t, "Find Key", key.Name)
	})

	t.Run("find non-existing key", func(t *testing.T) {
		key, err := repo.FindByKey(ctx, "sk_nonexistent")
		assert.Error(t, err)
		assert.Nil(t, key)
	})
}

func TestAPIKeyRepo_ListByUserID(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.APIKey{})
	repo := NewAPIKeyRepo(engine)
	ctx := context.Background()

	now := time.Now()
	_, _ = engine.Insert(&model.APIKey{UserID: 1, Key: "sk_1", Name: "Key 1", ExpiresAt: &now.Add(24 * time.Hour), Status: 1})
	_, _ = engine.Insert(&model.APIKey{UserID: 1, Key: "sk_2", Name: "Key 2", ExpiresAt: &now.Add(24 * time.Hour), Status: 1})
	_, _ = engine.Insert(&model.APIKey{UserID: 2, Key: "sk_3", Name: "Key 3", ExpiresAt: &now.Add(24 * time.Hour), Status: 1})

	keys, err := repo.ListByUserID(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, 2, len(keys))
}

func TestAPIKeyRepo_Update(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.APIKey{})
	repo := NewAPIKeyRepo(engine)
	ctx := context.Background()

	now := time.Now()
	_, _ = engine.Insert(&model.APIKey{
		UserID: 1, Key: "sk_update", Name: "Original",
		ExpiresAt: &now.Add(24 * time.Hour), Status: 1,
	})

	key, _ := repo.FindByKey(ctx, "sk_update")
	key.Name = "Updated Name"
	err := repo.Update(ctx, key)
	require.NoError(t, err)

	updated, _ := repo.FindByKey(ctx, "sk_update")
	assert.Equal(t, "Updated Name", updated.Name)
}

func TestAPIKeyRepo_Delete(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.APIKey{})
	repo := NewAPIKeyRepo(engine)
	ctx := context.Background()

	now := time.Now()
	_, _ = engine.Insert(&model.APIKey{
		UserID: 1, Key: "sk_delete", Name: "Delete Me",
		ExpiresAt: &now.Add(24 * time.Hour), Status: 1,
	})

	key, _ := repo.FindByKey(ctx, "sk_delete")
	err := repo.Delete(ctx, key.ID)
	require.NoError(t, err)

	deleted, err := repo.FindByKey(ctx, "sk_delete")
	assert.Error(t, err)
	assert.Nil(t, deleted)
}
