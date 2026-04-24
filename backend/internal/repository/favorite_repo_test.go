package repository

import (
	"testing"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFavoriteRepo_Create(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Favorite{}, &model.Skill{}, &model.User{})
	repo := NewFavoriteRepo(engine)

	fav := &model.Favorite{
		UserID:  1,
		SkillID: 1,
	}

	err := repo.Create(fav)
	require.NoError(t, err)
	assert.Greater(t, fav.ID, int64(0))
}

func TestFavoriteRepo_Exists(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Favorite{})
	repo := NewFavoriteRepo(engine)

	_, _ = engine.Insert(&model.Favorite{UserID: 1, SkillID: 1})

	exists, err := repo.Exists(1, 1)
	require.NoError(t, err)
	assert.True(t, exists)

	notExists, err := repo.Exists(1, 999)
	require.NoError(t, err)
	assert.False(t, notExists)
}

func TestFavoriteRepo_Delete(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Favorite{})
	repo := NewFavoriteRepo(engine)

	_, _ = engine.Insert(&model.Favorite{UserID: 1, SkillID: 1})

	err := repo.Delete(1, 1)
	require.NoError(t, err)

	exists, _ := repo.Exists(1, 1)
	assert.False(t, exists)
}

func TestFavoriteRepo_ListByUser(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Favorite{})
	repo := NewFavoriteRepo(engine)

	favs := []*model.Favorite{
		{UserID: 1, SkillID: 1},
		{UserID: 1, SkillID: 2},
		{UserID: 2, SkillID: 3},
	}
	for _, f := range favs {
		_, _ = engine.Insert(f)
	}

	t.Run("list user's favorites", func(t *testing.T) {
		results, total, err := repo.ListByUser(1, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, 2, len(results))
	})

	t.Run("list with pagination", func(t *testing.T) {
		results, total, err := repo.ListByUser(1, 1, 1)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, 1, len(results))
	})

	t.Run("other user not affected", func(t *testing.T) {
		results, total, err := repo.ListByUser(2, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, 1, len(results))
		assert.Equal(t, int64(3), results[0].SkillID)
	})
}

func TestFavoriteRepo_Count(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Favorite{})
	repo := NewFavoriteRepo(engine)

	_, _ = engine.Insert(&model.Favorite{UserID: 1, SkillID: 1})
	_, _ = engine.Insert(&model.Favorite{UserID: 1, SkillID: 2})
	_, _ = engine.Insert(&model.Favorite{UserID: 2, SkillID: 1})

	userCount, err := repo.CountByUser(1)
	require.NoError(t, err)
	assert.Equal(t, int64(2), userCount)

	skillCount, err := repo.CountBySkill(1)
	require.NoError(t, err)
	assert.Equal(t, int64(2), skillCount)
}
