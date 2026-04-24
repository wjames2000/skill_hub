package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVectorSearchRequest(t *testing.T) {
	req := &VectorSearchRequest{
		Query:   "analyze data trends",
		TopK:    10,
		Filters: map[string]interface{}{"category": "data-analysis"},
	}
	assert.Equal(t, "analyze data trends", req.Query)
	assert.Equal(t, 10, req.TopK)
	assert.Equal(t, "data-analysis", req.Filters["category"])
}

func TestVectorSearchResponse(t *testing.T) {
	resp := &VectorSearchResponse{
		Results: []VectorSearchResult{
			{SkillID: 1, Score: 0.95, Text: "excel analyzer"},
			{SkillID: 2, Score: 0.87, Text: "ppt generator"},
		},
		TotalTime: 50,
	}
	assert.Equal(t, 2, len(resp.Results))
	assert.Equal(t, int64(1), resp.Results[0].SkillID)
	assert.Equal(t, 0.95, resp.Results[0].Score)
	assert.Equal(t, "ppt generator", resp.Results[1].Text)
	assert.Equal(t, int64(50), resp.TotalTime)
}

type mockVectorService struct {
	searchFunc func(ctx context.Context, req *VectorSearchRequest) (*VectorSearchResponse, error)
}

func (m *mockVectorService) Search(ctx context.Context, req *VectorSearchRequest) (*VectorSearchResponse, error) {
	if m.searchFunc != nil {
		return m.searchFunc(ctx, req)
	}
	return &VectorSearchResponse{Results: []VectorSearchResult{}, TotalTime: 0}, nil
}

func TestVectorServiceMock(t *testing.T) {
	svc := &mockVectorService{
		searchFunc: func(ctx context.Context, req *VectorSearchRequest) (*VectorSearchResponse, error) {
			return &VectorSearchResponse{
				Results: []VectorSearchResult{
					{SkillID: 1, Score: 0.98, Text: "best match"},
				},
				TotalTime: 10,
			}, nil
		},
	}

	resp, err := svc.Search(context.Background(), &VectorSearchRequest{Query: "test", TopK: 5})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Results))
	assert.Equal(t, float64(0.98), resp.Results[0].Score)
	assert.Equal(t, int64(10), resp.TotalTime)
}
