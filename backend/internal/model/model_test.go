package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSkillStatus(t *testing.T) {
	assert.Equal(t, 1, SkillStatusActive)
	assert.Equal(t, 2, SkillStatusInactive)
	assert.Equal(t, 3, SkillStatusDeprecated)
}

func TestSkillDefault(t *testing.T) {
	s := &Skill{}
	assert.Equal(t, int64(0), s.ID)
	assert.Empty(t, s.Name)
	assert.Empty(t, s.DisplayName)
	assert.Empty(t, s.Description)
	assert.Equal(t, 0, s.Stars)
	assert.Equal(t, 0, s.Status)
	assert.Nil(t, s.CreatedAt)
	assert.Nil(t, s.UpdatedAt)
}

func TestSkillCategory(t *testing.T) {
	now := time.Now()
	cat := &SkillCategory{
		ID:          1,
		Name:        "data-analysis",
		DisplayName: "Data Analysis",
		Description: "Data analysis tools",
		SortOrder:   1,
		Status:      1,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	assert.Equal(t, "data-analysis", cat.Name)
	assert.Equal(t, "Data Analysis", cat.DisplayName)
	assert.Equal(t, 1, cat.SortOrder)
	assert.Equal(t, 1, cat.Status)
}

func TestUser(t *testing.T) {
	u := &User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hash",
		Role:         "user",
		Status:       1,
	}
	assert.Equal(t, "testuser", u.Username)
	assert.Equal(t, "test@example.com", u.Email)
	assert.Equal(t, "user", u.Role)
	assert.True(t, u.IsActive())

	u.Status = 0
	assert.False(t, u.IsActive())
}

func TestAPIKey(t *testing.T) {
	now := time.Now()
	key := &APIKey{
		ID:        1,
		UserID:    1,
		Key:       "sk-test-key",
		Name:      "Development",
		LastUsed:  &now,
		ExpiresAt: &now,
		Status:    1,
	}
	assert.Equal(t, "sk-test-key", key.Key)
	assert.Equal(t, "Development", key.Name)
	assert.True(t, key.IsValid())

	key.Status = 0
	assert.False(t, key.IsValid())
}

func TestRouterLog(t *testing.T) {
	now := time.Now()
	log := &RouterLog{
		ID:          1,
		SessionID:   "sess_abc",
		Query:       "analyze data",
		MatchedSkill: []byte(`[{"skill_id":1,"score":0.95}]`),
		Strategy:    "hybrid",
		Duration:    150,
		Status:      1,
		CreatedAt:   &now,
	}
	assert.Equal(t, "sess_abc", log.SessionID)
	assert.Equal(t, "analyze data", log.Query)
	assert.Equal(t, "hybrid", log.Strategy)
	assert.Equal(t, int64(150), log.Duration)
	assert.Equal(t, 1, log.Status)
}

func TestSyncState(t *testing.T) {
	state := &SyncState{
		ID:          1,
		Source:      "github",
		Status:      "running",
		SkillsCount: 42,
	}
	assert.Equal(t, "github", state.Source)
	assert.Equal(t, "running", state.Status)
	assert.Equal(t, 42, state.SkillsCount)
}
