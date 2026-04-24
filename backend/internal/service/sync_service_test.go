package service

import (
	"testing"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestSyncTaskType(t *testing.T) {
	task := &SyncTask{
		ID:      "task_001",
		Type:    SyncTypeGitHub,
		Payload: map[string]interface{}{"repo": "owner/repo"},
	}
	assert.Equal(t, "github", string(task.Type))
	assert.Equal(t, "task_001", task.ID)
}

func TestSyncResult(t *testing.T) {
	result := &SyncResult{
		TaskID:        "task_001",
		Success:       true,
		SkillsAdded:   5,
		SkillsUpdated: 3,
		SkillsRemoved: 1,
		Errors:        nil,
	}
	assert.True(t, result.Success)
	assert.Equal(t, 5, result.SkillsAdded)
	assert.Equal(t, 3, result.SkillsUpdated)
	assert.Empty(t, result.Errors)
}

func TestSyncResultWithErrors(t *testing.T) {
	result := &SyncResult{
		TaskID:        "task_002",
		Success:       false,
		SkillsAdded:   0,
		SkillsUpdated: 0,
		SkillsRemoved: 0,
		Errors:        []string{"failed to fetch repo: timeout", "invalid manifest"},
	}
	assert.False(t, result.Success)
	assert.Equal(t, 2, len(result.Errors))
	assert.Contains(t, result.Errors[0], "timeout")
}

func TestSyncTaskValidation(t *testing.T) {
	t.Run("github type has correct string", func(t *testing.T) {
		assert.Equal(t, "github", string(SyncTypeGitHub))
	})

	t.Run("gitlab type has correct string", func(t *testing.T) {
		assert.Equal(t, "gitlab", string(SyncTypeGitLab))
	})
}

func TestModelSyncState(t *testing.T) {
	state := &model.SyncState{
		ID:          1,
		Source:      "github",
		LastSyncAt:  nil,
		Status:      "pending",
		SkillsCount: 10,
	}
	assert.Equal(t, int64(1), state.ID)
	assert.Equal(t, "github", state.Source)
	assert.Equal(t, "pending", state.Status)
	assert.Equal(t, 10, state.SkillsCount)
	assert.Nil(t, state.LastSyncAt)
}
