package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"encoding/json"

	"github.com/hpds/skill-hub/internal/client/embedding"
	"github.com/hpds/skill-hub/internal/client/llm"
	"github.com/hpds/skill-hub/internal/client/reranker"
	"github.com/hpds/skill-hub/internal/milvus"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/pkg/logger"
	mls "github.com/hpds/skill-hub/pkg/meilisearch"
)

type RouterService struct {
	embedder    *embedding.Client
	llmClient   *llm.Client
	rerankerCli *reranker.Client
	milvusCli   *milvus.Client
	meiliCli    *mls.Client
	skillRepo   *repository.SkillRepo
	embRepo     *repository.EmbeddingRepo
	logRepo     *repository.RouterLogRepo
}

func NewRouterService(
	embedder *embedding.Client,
	llmClient *llm.Client,
	rerankerCli *reranker.Client,
	milvusCli *milvus.Client,
	meiliCli *mls.Client,
	skillRepo *repository.SkillRepo,
	embRepo *repository.EmbeddingRepo,
	logRepo *repository.RouterLogRepo,
) *RouterService {
	return &RouterService{
		embedder:    embedder,
		llmClient:   llmClient,
		rerankerCli: rerankerCli,
		milvusCli:   milvusCli,
		meiliCli:    meiliCli,
		skillRepo:   skillRepo,
		embRepo:     embRepo,
		logRepo:     logRepo,
	}
}

type MatchedSkill struct {
	Skill    *model.Skill `json:"skill"`
	Score    float64      `json:"score"`
	Strategy string       `json:"strategy"`
}

type MatchRequest struct {
	Query    string `json:"query"`
	TopK     int    `json:"top_k,omitempty"`
	UserID   int64  `json:"user_id,omitempty"`
	Strategy string `json:"strategy,omitempty"`
}

type MatchResponse struct {
	MatchedSkills []*MatchedSkill `json:"matched_skills"`
	Strategy      string          `json:"strategy"`
	TotalTime     int64           `json:"total_time_ms"`
}

type ExecuteRequest struct {
	SessionID string `json:"session_id"`
	Query     string `json:"query"`
	SkillID   int64  `json:"skill_id"`
	UserID    int64  `json:"user_id,omitempty"`
}

type ExecuteResponse struct {
	SessionID string `json:"session_id"`
	Result    string `json:"result"`
	Duration  int    `json:"duration_ms"`
}

type FeedbackRequest struct {
	SessionID string `json:"session_id"`
	LogID     int64  `json:"log_id"`
	Score     int    `json:"score"`
	Comment   string `json:"comment,omitempty"`
}

func (s *RouterService) Match(ctx context.Context, req *MatchRequest) (*MatchResponse, error) {
	start := time.Now()

	if req.TopK <= 0 {
		req.TopK = 10
	}
	if req.TopK > 50 {
		req.TopK = 50
	}

	strategy := req.Strategy
	if strategy == "" {
		strategy = "hybrid"
	}

	var skills []*MatchedSkill
	var err error

	switch strategy {
	case "vector":
		skills, err = s.vectorSearch(ctx, req.Query, req.TopK)
	case "keyword":
		skills, err = s.keywordSearch(ctx, req.Query, req.TopK)
	case "hybrid":
		skills, err = s.hybridSearch(ctx, req.Query, req.TopK)
	default:
		skills, err = s.hybridSearch(ctx, req.Query, req.TopK)
	}

	if err != nil {
		return nil, fmt.Errorf("match search: %w", err)
	}

	if len(skills) > 0 && s.rerankerCli != nil {
		reranked, err := s.rerankResults(ctx, req.Query, skills, req.TopK)
		if err == nil {
			skills = reranked
			strategy = "hybrid+rerank"
		}
	}

	totalTime := time.Since(start).Milliseconds()

	return &MatchResponse{
		MatchedSkills: skills,
		Strategy:      strategy,
		TotalTime:     totalTime,
	}, nil
}

func (s *RouterService) Execute(ctx context.Context, req *ExecuteRequest) (*ExecuteResponse, error) {
	start := time.Now()

	skill, err := s.skillRepo.GetByID(req.SkillID)
	if err != nil {
		return nil, fmt.Errorf("get skill: %w", err)
	}
	if skill == nil {
		return nil, fmt.Errorf("skill not found: %d", req.SkillID)
	}

	contextPrompt := s.buildContext(skill, req.Query)

	result, err := s.llmClient.Chat(contextPrompt, req.Query)
	if err != nil {
		return nil, fmt.Errorf("llm execute: %w", err)
	}

	duration := int(time.Since(start).Milliseconds())

	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("sess_%d", time.Now().UnixNano())
	}

	logEntry := &model.RouterLog{
		SessionID:        sessionID,
		Query:            req.Query,
		MatchedSkillID:   skill.ID,
		MatchedSkillName: skill.Name,
		MatchScore:       1.0,
		MatchStrategy:    "execute",
		IsExecuted:       true,
		ExecuteResult:    result,
		ExecuteDuration:  duration,
		UserID:           req.UserID,
	}

	if err := s.logRepo.Create(logEntry); err != nil {
		logger.Error("save router log", logger.ErrorField(err))
	}

	return &ExecuteResponse{
		SessionID: sessionID,
		Result:    result,
		Duration:  duration,
	}, nil
}

func (s *RouterService) SubmitFeedback(ctx context.Context, req *FeedbackRequest) error {
	if req.Score < 1 || req.Score > 5 {
		return fmt.Errorf("score must be between 1 and 5")
	}

	return s.logRepo.UpdateFeedback(req.LogID, req.Score, req.Comment)
}

func (s *RouterService) vectorSearch(ctx context.Context, query string, topK int) ([]*MatchedSkill, error) {
	queryVec, err := s.embedder.Embed(query)
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}

	results, err := s.embRepo.SearchSimilar(queryVec, s.embedder.Model(), topK*2)
	if err != nil {
		logger.Warn("pg vector search failed, falling back to keyword", logger.ErrorField(err))

		meiliResults, meiliErr := s.meiliCli.Search("skills", query, int64(topK))
		if meiliErr != nil {
			return nil, fmt.Errorf("keyword fallback: %w", meiliErr)
		}

		var skills []*MatchedSkill
		for _, hit := range meiliResults.Hits {
			id, score := extractHitFields(hit)
			if id <= 0 {
				continue
			}
			skill, err := s.skillRepo.GetByID(id)
			if err != nil || skill == nil {
				continue
			}
			skills = append(skills, &MatchedSkill{
				Skill:    skill,
				Score:    score,
				Strategy: "keyword_fallback",
			})
		}
		return skills, nil
	}

	skillMap := make(map[int64]*model.Skill)
	var mu sync.Mutex
	var wg sync.WaitGroup
	skillCh := make(chan *model.Skill, len(results))

	for _, r := range results {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			skill, err := s.skillRepo.GetByID(id)
			if err != nil || skill == nil {
				return
			}
			skillCh <- skill
		}(r.ID)
	}

	go func() {
		wg.Wait()
		close(skillCh)
	}()

	for skill := range skillCh {
		mu.Lock()
		skillMap[skill.ID] = skill
		mu.Unlock()
	}

	var skills []*MatchedSkill
	for _, r := range results {
		if skill, ok := skillMap[r.ID]; ok {
			skills = append(skills, &MatchedSkill{
				Skill:    skill,
				Score:    r.Score,
				Strategy: "vector",
			})
		}
	}

	if len(skills) > topK {
		skills = skills[:topK]
	}

	return skills, nil
}

func (s *RouterService) keywordSearch(ctx context.Context, query string, topK int) ([]*MatchedSkill, error) {
	if s.meiliCli == nil {
		return s.fallbackKeywordSearch(ctx, query, topK)
	}

	resp, err := s.meiliCli.Search("skills", query, int64(topK))
	if err != nil {
		return s.fallbackKeywordSearch(ctx, query, topK)
	}

	var skills []*MatchedSkill
	for _, hit := range resp.Hits {
		id, score := extractHitFields(hit)
		if id <= 0 {
			continue
		}
		skill, err := s.skillRepo.GetByID(id)
		if err != nil || skill == nil {
			continue
		}
		skills = append(skills, &MatchedSkill{
			Skill:    skill,
			Score:    score,
			Strategy: "keyword",
		})
	}

	return skills, nil
}

func (s *RouterService) fallbackKeywordSearch(ctx context.Context, query string, topK int) ([]*MatchedSkill, error) {
	skills, err := s.skillRepo.SearchByName(query, topK)
	if err != nil {
		return nil, err
	}

	var result []*MatchedSkill
	queryLower := strings.ToLower(query)
	for _, skill := range skills {
		score := 0.0
		nameLower := strings.ToLower(skill.Name)
		descLower := strings.ToLower(skill.Description)

		if strings.Contains(nameLower, queryLower) {
			score += 10.0
		}
		if strings.Contains(descLower, queryLower) {
			score += 5.0
		}

		result = append(result, &MatchedSkill{
			Skill:    skill,
			Score:    score,
			Strategy: "keyword",
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	if len(result) > topK {
		result = result[:topK]
	}

	return result, nil
}

func (s *RouterService) hybridSearch(ctx context.Context, query string, topK int) ([]*MatchedSkill, error) {
	var wg sync.WaitGroup
	var vectorSkills, keywordSkills []*MatchedSkill
	var vectorErr, keywordErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		vectorSkills, vectorErr = s.vectorSearch(ctx, query, topK*2)
	}()
	go func() {
		defer wg.Done()
		keywordSkills, keywordErr = s.keywordSearch(ctx, query, topK*2)
	}()
	wg.Wait()

	if vectorErr != nil && keywordErr != nil {
		return nil, fmt.Errorf("both searches failed: vector=%v keyword=%v", vectorErr, keywordErr)
	}

	if vectorErr != nil || len(vectorSkills) == 0 {
		if len(keywordSkills) > topK {
			keywordSkills = keywordSkills[:topK]
		}
		return keywordSkills, nil
	}
	if keywordErr != nil || len(keywordSkills) == 0 {
		if len(vectorSkills) > topK {
			vectorSkills = vectorSkills[:topK]
		}
		return vectorSkills, nil
	}

	fused := s.rrfFusion(vectorSkills, keywordSkills, topK)
	for _, s := range fused {
		s.Strategy = "hybrid"
	}

	return fused, nil
}

func (s *RouterService) rrfFusion(vector, keyword []*MatchedSkill, topK int) []*MatchedSkill {
	k := 60.0

	scoreMap := make(map[int64]*MatchedSkill)

	for i, item := range vector {
		rank := float64(i + 1)
		rrfScore := 1.0 / (k + rank)

		if existing, ok := scoreMap[item.Skill.ID]; ok {
			existing.Score += rrfScore
		} else {
			scoreMap[item.Skill.ID] = &MatchedSkill{
				Skill:    item.Skill,
				Score:    rrfScore,
				Strategy: "hybrid",
			}
		}
	}

	for i, item := range keyword {
		rank := float64(i + 1)
		rrfScore := 1.0 / (k + rank)

		if existing, ok := scoreMap[item.Skill.ID]; ok {
			existing.Score += rrfScore
		} else {
			scoreMap[item.Skill.ID] = &MatchedSkill{
				Skill:    item.Skill,
				Score:    rrfScore,
				Strategy: "hybrid",
			}
		}
	}

	var result []*MatchedSkill
	for _, item := range scoreMap {
		result = append(result, item)
	}

	sort.Slice(result, func(i, j int) bool {
		if math.Abs(result[i].Score-result[j].Score) < 0.0001 {
			return result[i].Skill.Stars > result[j].Skill.Stars
		}
		return result[i].Score > result[j].Score
	})

	if len(result) > topK {
		result = result[:topK]
	}

	return result
}

func (s *RouterService) rerankResults(ctx context.Context, query string, skills []*MatchedSkill, topK int) ([]*MatchedSkill, error) {
	documents := make([]string, len(skills))
	for i, ms := range skills {
		documents[i] = fmt.Sprintf("%s: %s", ms.Skill.Name, truncateString(ms.Skill.Description, 200))
	}

	results, err := s.rerankerCli.Rerank(query, documents, topK)
	if err != nil {
		return nil, fmt.Errorf("rerank: %w", err)
	}

	reranked := make([]*MatchedSkill, len(results))
	for i, r := range results {
		if r.Index >= 0 && r.Index < len(skills) {
			reranked[i] = &MatchedSkill{
				Skill:    skills[r.Index].Skill,
				Score:    r.Score,
				Strategy: "rerank",
			}
		}
	}

	return reranked, nil
}

func (s *RouterService) buildContext(skill *model.Skill, query string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("You are now acting as the skill/assistant: **%s**.\n\n", skill.DisplayName))
	if skill.DisplayName == "" {
		b.Reset()
		b.WriteString(fmt.Sprintf("You are now acting as the skill/assistant: **%s**.\n\n", skill.Name))
	}

	if skill.Description != "" {
		b.WriteString(fmt.Sprintf("## Description\n%s\n\n", skill.Description))
	}

	if skill.Readme != "" {
		maxReadme := 3000
		readmeContent := skill.Readme
		if len(readmeContent) > maxReadme {
			readmeContent = readmeContent[:maxReadme] + "\n...(truncated)"
		}
		b.WriteString(fmt.Sprintf("## Instructions / README\n%s\n\n", readmeContent))
	}

	b.WriteString(fmt.Sprintf("## Metadata\n- Repository: %s\n- Author: %s\n- Category: %s\n- Tags: %v\n- Language: %s\n- Version: %s\n\n",
		skill.Repository, skill.Author, skill.Category, skill.Tags, skill.Language, skill.Version))

	b.WriteString("## Guidelines\n")
	b.WriteString("1. Answer based on the skill's documented capabilities and instructions.\n")
	b.WriteString("2. If the user's request is outside the skill's scope, politely explain what this skill can do.\n")
	b.WriteString("3. Provide practical, actionable responses.\n")
	b.WriteString("4. When appropriate, include code examples or configuration snippets.\n")
	b.WriteString(fmt.Sprintf("5. The current user query is: %s\n", query))

	return b.String()
}

func extractHitFields(hit map[string]json.RawMessage) (skillID int64, score float64) {
	if raw, ok := hit["id"]; ok {
		var idFloat float64
		if err := json.Unmarshal(raw, &idFloat); err == nil {
			skillID = int64(idFloat)
		}
	}
	if raw, ok := hit["_rankingScore"]; ok {
		_ = json.Unmarshal(raw, &score)
	}
	return
}
