package repository

import (
	"context"
	"testing"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"xorm.io/xorm"
)

func newTestEngine() *xorm.Engine {
	engine, err := xorm.NewEngine("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	engine.ShowSQL(false)
	_ = engine.Sync2(&model.Skill{}, &model.SkillCategory{})
	return engine
}

func seedSkills(engine *xorm.Engine) {
	skills := []*model.Skill{
		{Name: "skill-a", DisplayName: "Skill A", Description: "First skill", Stars: 100, Status: model.SkillStatusActive, Category: "cat1"},
		{Name: "skill-b", DisplayName: "Skill B", Description: "Second skill", Stars: 200, Status: model.SkillStatusActive, Category: "cat2"},
		{Name: "inactive-skill", DisplayName: "Inactive", Description: "Inactive skill", Stars: 50, Status: model.SkillStatusInactive, Category: "cat1"},
	}
	for _, s := range skills {
		_, _ = engine.Insert(s)
	}
}

func TestSkillRepo_List(t *testing.T) {
	engine := newTestEngine()
	seedSkills(engine)
	repo := NewSkillRepo(engine)

	t.Run("list all active skills", func(t *testing.T) {
		skills, total, err := repo.List(context.Background(), "", "", "stars", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, 2, len(skills))
	})

	t.Run("list with category filter", func(t *testing.T) {
		skills, total, err := repo.List(context.Background(), "cat1", "", "stars", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, 1, len(skills))
		assert.Equal(t, "skill-a", skills[0].Name)
	})

	t.Run("list with pagination", func(t *testing.T) {
		skills, total, err := repo.List(context.Background(), "", "", "stars", 1, 1)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, 1, len(skills))
	})

	t.Run("list with sort by name", func(t *testing.T) {
		skills, total, err := repo.List(context.Background(), "", "", "name", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, "skill-a", skills[0].Name)
		assert.Equal(t, "skill-b", skills[1].Name)
	})

	t.Run("list with sort by stars desc", func(t *testing.T) {
		skills, total, err := repo.List(context.Background(), "", "", "stars", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, "skill-b", skills[0].Name) // 200 stars first
		assert.Equal(t, "skill-a", skills[1].Name) // 100 stars second
	})
}

func TestSkillRepo_GetByID(t *testing.T) {
	engine := newTestEngine()
	seedSkills(engine)
	repo := NewSkillRepo(engine)

	t.Run("get existing skill", func(t *testing.T) {
		skill, err := repo.GetByID(context.Background(), 1)
		require.NoError(t, err)
		assert.Equal(t, "skill-a", skill.Name)
	})

	t.Run("get non-existing skill", func(t *testing.T) {
		skill, err := repo.GetByID(context.Background(), 999)
		assert.Error(t, err)
		assert.Nil(t, skill)
	})
}

func TestSkillRepo_Search(t *testing.T) {
	engine := newTestEngine()
	seedSkills(engine)
	repo := NewSkillRepo(engine)

	t.Run("search by name", func(t *testing.T) {
		skills, total, err := repo.Search(context.Background(), "skill-a", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, 1, len(skills))
		assert.Equal(t, "skill-a", skills[0].Name)
	})

	t.Run("search by description", func(t *testing.T) {
		skills, total, err := repo.Search(context.Background(), "First", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, 1, len(skills))
	})

	t.Run("search with no results", func(t *testing.T) {
		skills, total, err := repo.Search(context.Background(), "zzzzz", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, skills)
	})

	t.Run("search only returns active skills", func(t *testing.T) {
		skills, total, err := repo.Search(context.Background(), "Inactive", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, skills)
	})

	t.Run("search with empty query returns all active", func(t *testing.T) {
		skills, total, err := repo.Search(context.Background(), "", 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, 2, len(skills))
	})
}

func TestSkillRepo_ListByIDs(t *testing.T) {
	engine := newTestEngine()
	seedSkills(engine)
	repo := NewSkillRepo(engine)

	t.Run("list by valid ids", func(t *testing.T) {
		skills, err := repo.ListByIDs(context.Background(), []int64{1, 2})
		require.NoError(t, err)
		assert.Equal(t, 2, len(skills))
	})

	t.Run("list by invalid ids returns empty", func(t *testing.T) {
		skills, err := repo.ListByIDs(context.Background(), []int64{999, 1000})
		require.NoError(t, err)
		assert.Empty(t, skills)
	})

	t.Run("list by empty ids returns empty", func(t *testing.T) {
		skills, err := repo.ListByIDs(context.Background(), []int64{})
		require.NoError(t, err)
		assert.Empty(t, skills)
	})
}

func TestSkillRepo_CRUD(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Skill{})
	repo := NewSkillRepo(engine)
	ctx := context.Background()

	t.Run("create skill", func(t *testing.T) {
		skill := &model.Skill{
			Name:        "new-skill",
			DisplayName: "New Skill",
			Stars:       50,
			Status:      model.SkillStatusActive,
		}
		err := repo.Create(ctx, skill)
		require.NoError(t, err)
		assert.Greater(t, skill.ID, int64(0))
	})

	t.Run("update skill", func(t *testing.T) {
		skill, _ := repo.GetByID(ctx, 1)
		skill.Stars = 999
		err := repo.Update(ctx, skill)
		require.NoError(t, err)

		updated, _ := repo.GetByID(ctx, 1)
		assert.Equal(t, 999, updated.Stars)
	})

	t.Run("delete skill", func(t *testing.T) {
		err := repo.Delete(ctx, 1)
		require.NoError(t, err)

		deleted, err := repo.GetByID(ctx, 1)
		assert.Error(t, err)
		assert.Nil(t, deleted)
	})
}

func TestNewSkillRepo(t *testing.T) {
	engine := newTestEngine()
	repo := NewSkillRepo(engine)
	assert.NotNil(t, repo)
}
