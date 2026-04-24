package service

import (
	"context"
	"testing"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMatchRequest_Validate(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		req := &MatchRequest{Query: "analyze data", TopK: 5}
		assert.NoError(t, req.Validate())
	})

	t.Run("empty query", func(t *testing.T) {
		req := &MatchRequest{Query: "", TopK: 5}
		assert.Error(t, req.Validate())
	})

	t.Run("zero topk defaults to 10", func(t *testing.T) {
		req := &MatchRequest{Query: "test", TopK: 0}
		assert.NoError(t, req.Validate())
	})

	t.Run("topk capped at 50", func(t *testing.T) {
		req := &MatchRequest{Query: "test", TopK: 100}
		assert.NoError(t, req.Validate())
	})
}

func TestExecuteRequest_Validate(t *testing.T) {
	t.Run("valid execute request", func(t *testing.T) {
		req := &ExecuteRequest{Query: "analyze", SkillID: 1}
		assert.NoError(t, req.Validate())
	})

	t.Run("empty query", func(t *testing.T) {
		req := &ExecuteRequest{Query: "", SkillID: 1}
		assert.Error(t, req.Validate())
	})

	t.Run("invalid skill id", func(t *testing.T) {
		req := &ExecuteRequest{Query: "test", SkillID: 0}
		assert.Error(t, req.Validate())
	})
}

func TestFeedbackRequest_Validate(t *testing.T) {
	t.Run("valid feedback", func(t *testing.T) {
		req := &FeedbackRequest{LogID: 1, Score: 4}
		assert.NoError(t, req.Validate())
	})

	t.Run("invalid log id", func(t *testing.T) {
		req := &FeedbackRequest{LogID: 0, Score: 4}
		assert.Error(t, req.Validate())
	})

	t.Run("score too low", func(t *testing.T) {
		req := &FeedbackRequest{LogID: 1, Score: 0}
		assert.Error(t, req.Validate())
	})

	t.Run("score too high", func(t *testing.T) {
		req := &FeedbackRequest{LogID: 1, Score: 6}
		assert.Error(t, req.Validate())
	})
}

func TestMatchedSkill(t *testing.T) {
	skill := &model.Skill{ID: 1, Name: "test-skill", DisplayName: "Test Skill", Description: "A test", Stars: 100, Status: 1}
	ms := &MatchedSkill{Skill: skill, Score: 0.95, Strategy: "hybrid"}

	assert.Equal(t, "test-skill", ms.Skill.Name)
	assert.Equal(t, 0.95, ms.Score)
	assert.Equal(t, "hybrid", ms.Strategy)
}

func TestMatchResponse(t *testing.T) {
	resp := &MatchResponse{
		MatchedSkills: []*MatchedSkill{},
		Strategy:      "hybrid",
		TotalTime:     150,
	}
	assert.Equal(t, "hybrid", resp.Strategy)
	assert.Equal(t, int64(150), resp.TotalTime)
	assert.Empty(t, resp.MatchedSkills)
}

func TestExecuteResponse(t *testing.T) {
	resp := &ExecuteResponse{
		SessionID: "sess_abc123",
		Result:    "analysis complete",
		Duration:  500,
	}
	assert.Equal(t, "sess_abc123", resp.SessionID)
	assert.Equal(t, "analysis complete", resp.Result)
	assert.Equal(t, int64(500), resp.Duration)
}

type mockSkillRepo struct {
	skills []*model.Skill
	err    error
}

func (m *mockSkillRepo) ListByIDs(ctx context.Context, ids []int64) ([]*model.Skill, error) {
	if m.err != nil {
		return nil, m.err
	}
	result := make([]*model.Skill, 0)
	for _, s := range m.skills {
		for _, id := range ids {
			if s.ID == id {
				result = append(result, s)
				break
			}
		}
	}
	return result, nil
}

func TestRouterServiceWithMock(t *testing.T) {
	mockSkills := []*model.Skill{
		{ID: 1, Name: "excel-analyzer", DisplayName: "Excel Analyzer", Stars: 100, Status: 1},
		{ID: 2, Name: "ppt-generator", DisplayName: "PPT Generator", Stars: 200, Status: 1},
	}

	repo := &mockSkillRepo{skills: mockSkills}

	t.Run("list skills by ids", func(t *testing.T) {
		skills, err := repo.ListByIDs(context.Background(), []int64{1, 2})
		assert.NoError(t, err)
		assert.Equal(t, 2, len(skills))
	})

	t.Run("list empty ids returns empty", func(t *testing.T) {
		skills, err := repo.ListByIDs(context.Background(), []int64{})
		assert.NoError(t, err)
		assert.Empty(t, skills)
	})
}
