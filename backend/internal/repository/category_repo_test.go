package repository

import (
	"context"
	"testing"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryRepo_List(t *testing.T) {
	engine := newTestEngine()
	repo := NewCategoryRepo(engine)
	ctx := context.Background()

	// Seed categories
	cats := []*model.SkillCategory{
		{Name: "cat1", DisplayName: "Category 1", SortOrder: 2, Status: 1},
		{Name: "cat2", DisplayName: "Category 2", SortOrder: 1, Status: 1},
		{Name: "cat3", DisplayName: "Category 3", SortOrder: 3, Status: 0},
	}
	for _, c := range cats {
		_, _ = engine.Insert(c)
	}

	t.Run("list only active categories sorted by sort_order", func(t *testing.T) {
		categories, err := repo.List(ctx)
		require.NoError(t, err)
		assert.Equal(t, 2, len(categories))
		assert.Equal(t, "cat2", categories[0].Name)  // sort_order=1 first
		assert.Equal(t, "cat1", categories[1].Name)  // sort_order=2 second
	})

	t.Run("inactive categories excluded", func(t *testing.T) {
		categories, err := repo.List(ctx)
		require.NoError(t, err)
		for _, c := range categories {
			assert.Equal(t, 1, c.Status)
		}
	})
}

func TestCategoryRepo_GetByID(t *testing.T) {
	engine := newTestEngine()
	repo := NewCategoryRepo(engine)
	ctx := context.Background()

	cat := &model.SkillCategory{Name: "test-cat", DisplayName: "Test", SortOrder: 1, Status: 1}
	_, _ = engine.Insert(cat)

	t.Run("get existing category", func(t *testing.T) {
		c, err := repo.GetByID(ctx, cat.ID)
		require.NoError(t, err)
		assert.Equal(t, "test-cat", c.Name)
	})

	t.Run("get non-existing category", func(t *testing.T) {
		c, err := repo.GetByID(ctx, 999)
		assert.Error(t, err)
		assert.Nil(t, c)
	})
}

func TestCategoryRepo_Create(t *testing.T) {
	engine := newTestEngine()
	repo := NewCategoryRepo(engine)
	ctx := context.Background()

	cat := &model.SkillCategory{
		Name:        "new-cat",
		DisplayName: "New Category",
		Description: "A new test category",
		SortOrder:   5,
		Status:      1,
	}

	err := repo.Create(ctx, cat)
	require.NoError(t, err)
	assert.Greater(t, cat.ID, int64(0))

	saved, _ := repo.GetByID(ctx, cat.ID)
	assert.Equal(t, "new-cat", saved.Name)
}

func TestCategoryRepo_Update(t *testing.T) {
	engine := newTestEngine()
	repo := NewCategoryRepo(engine)
	ctx := context.Background()

	cat := &model.SkillCategory{Name: "update-cat", DisplayName: "Original", SortOrder: 1, Status: 1}
	_, _ = engine.Insert(cat)

	cat.DisplayName = "Updated"
	err := repo.Update(ctx, cat)
	require.NoError(t, err)

	updated, _ := repo.GetByID(ctx, cat.ID)
	assert.Equal(t, "Updated", updated.DisplayName)
}

func TestCategoryRepo_Delete(t *testing.T) {
	engine := newTestEngine()
	repo := NewCategoryRepo(engine)
	ctx := context.Background()

	cat := &model.SkillCategory{Name: "delete-cat", DisplayName: "To Delete", SortOrder: 1, Status: 1}
	_, _ = engine.Insert(cat)

	err := repo.Delete(ctx, cat.ID)
	require.NoError(t, err)

	deleted, err := repo.GetByID(ctx, cat.ID)
	assert.Error(t, err)
	assert.Nil(t, deleted)
}
