package repository

import (
	"testing"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReviewRepo_Create(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Review{}, &model.Skill{}, &model.User{})
	repo := NewReviewRepo(engine)

	review := &model.Review{
		UserID:  1,
		SkillID: 1,
		Score:   5,
		Content: "Great skill!",
	}

	err := repo.Create(review)
	require.NoError(t, err)
	assert.Greater(t, review.ID, int64(0))
}

func TestReviewRepo_GetByID(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Review{})
	repo := NewReviewRepo(engine)

	_, _ = engine.Insert(&model.Review{UserID: 1, SkillID: 1, Score: 4, Content: "Good"})

	review, err := repo.GetByID(1)
	require.NoError(t, err)
	require.NotNil(t, review)
	assert.Equal(t, 4, review.Score)
	assert.Equal(t, "Good", review.Content)

	notFound, err := repo.GetByID(999)
	require.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestReviewRepo_ListBySkill(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Review{})
	repo := NewReviewRepo(engine)

	_, _ = engine.Insert(&model.Review{UserID: 1, SkillID: 1, Score: 5, Content: "Excellent"})
	_, _ = engine.Insert(&model.Review{UserID: 2, SkillID: 1, Score: 3, Content: "Average"})
	_, _ = engine.Insert(&model.Review{UserID: 1, SkillID: 2, Score: 4, Content: "Other skill"})

	reviews, total, err := repo.ListBySkill(1, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, len(reviews))
}

func TestReviewRepo_ListByUser(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Review{})
	repo := NewReviewRepo(engine)

	_, _ = engine.Insert(&model.Review{UserID: 1, SkillID: 1, Score: 5, Content: "A"})
	_, _ = engine.Insert(&model.Review{UserID: 1, SkillID: 2, Score: 4, Content: "B"})
	_, _ = engine.Insert(&model.Review{UserID: 1, SkillID: 3, Score: 3, Content: "C"})

	reviews, total, err := repo.ListByUser(1, 1, 2)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Equal(t, 2, len(reviews))
}

func TestReviewRepo_ListAll(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Review{})
	repo := NewReviewRepo(engine)

	_, _ = engine.Insert(&model.Review{UserID: 1, SkillID: 1, Score: 5, Content: "A"})
	_, _ = engine.Insert(&model.Review{UserID: 2, SkillID: 2, Score: 4, Content: "B"})

	reviews, total, err := repo.ListAll(1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, len(reviews))
}

func TestReviewRepo_Update(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Review{})
	repo := NewReviewRepo(engine)

	_, _ = engine.Insert(&model.Review{UserID: 1, SkillID: 1, Score: 3, Content: "Original"})

	review, _ := repo.GetByID(1)
	review.Score = 5
	review.Content = "Updated"
	err := repo.Update(review)
	require.NoError(t, err)

	updated, _ := repo.GetByID(1)
	assert.Equal(t, 5, updated.Score)
	assert.Equal(t, "Updated", updated.Content)
}

func TestReviewRepo_Delete(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Review{})
	repo := NewReviewRepo(engine)

	_, _ = engine.Insert(&model.Review{UserID: 1, SkillID: 1, Score: 4, Content: "Delete me"})

	err := repo.Delete(1)
	require.NoError(t, err)

	deleted, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.Nil(t, deleted)
}

func TestReviewRepo_GetAvgScoreBySkill(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Review{})
	repo := NewReviewRepo(engine)

	_, _ = engine.Insert(&model.Review{UserID: 1, SkillID: 1, Score: 5, Content: "A"})
	_, _ = engine.Insert(&model.Review{UserID: 2, SkillID: 1, Score: 3, Content: "B"})
	_, _ = engine.Insert(&model.Review{UserID: 3, SkillID: 1, Score: 4, Content: "C"})

	avg, err := repo.GetAvgScoreBySkill(1)
	require.NoError(t, err)
	assert.InDelta(t, 4.0, avg, 0.01)

	noReviews, err := repo.GetAvgScoreBySkill(999)
	require.NoError(t, err)
	assert.Equal(t, 0.0, noReviews)
}

func TestReviewRepo_CountBySkill(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.Review{})
	repo := NewReviewRepo(engine)

	_, _ = engine.Insert(&model.Review{UserID: 1, SkillID: 1, Score: 5, Content: "A"})
	_, _ = engine.Insert(&model.Review{UserID: 2, SkillID: 1, Score: 4, Content: "B"})

	count, err := repo.CountBySkill(1)
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}
