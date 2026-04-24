package service

import (
	"context"
	"fmt"

	"github.com/hpds/skill-hub/internal/client/embedding"
	"github.com/hpds/skill-hub/internal/client/llm"
	"github.com/hpds/skill-hub/internal/milvus"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/internal/vectorizer"
	"github.com/hpds/skill-hub/pkg/logger"
)

type VectorService struct {
	embedder  *embedding.Client
	llmClient *llm.Client
	milvusCli *milvus.Client
	skillRepo *repository.SkillRepo
	embRepo   *repository.EmbeddingRepo
	worker    *vectorizer.Worker
}

func NewVectorService(
	embedder *embedding.Client,
	llmClient *llm.Client,
	milvusCli *milvus.Client,
	skillRepo *repository.SkillRepo,
	embRepo *repository.EmbeddingRepo,
	worker *vectorizer.Worker,
) *VectorService {
	return &VectorService{
		embedder:  embedder,
		llmClient: llmClient,
		milvusCli: milvusCli,
		skillRepo: skillRepo,
		embRepo:   embRepo,
		worker:    worker,
	}
}

func (s *VectorService) VectorizeSkill(ctx context.Context, skillID int64) error {
	s.worker.Enqueue(skillID)
	return nil
}

func (s *VectorService) RevectorizeAll(ctx context.Context) error {
	ids, err := s.embRepo.GetSkillIDsNeedingVectorization(100)
	if err != nil {
		return fmt.Errorf("get unvectorized skills: %w", err)
	}
	for _, id := range ids {
		s.worker.Enqueue(id)
	}
	logger.Info("revectorize all queued", logger.Int("count", len(ids)))
	return nil
}

func (s *VectorService) GenerateEnhancedDescription(ctx context.Context, skillID int64) (string, error) {
	skill, err := s.skillRepo.GetByID(skillID)
	if err != nil {
		return "", fmt.Errorf("get skill: %w", err)
	}
	if skill == nil {
		return "", fmt.Errorf("skill not found: %d", skillID)
	}

	systemPrompt := `You are a skill documentation expert. Your task is to generate a concise, informative enhanced description for a skill/tool based on its metadata and README. Focus on:
1. What the skill does (core functionality)
2. Key features and capabilities
3. Use cases and scenarios
4. Requirements and dependencies

Keep the description under 500 words. Use clear, technical language.`

	userPrompt := fmt.Sprintf(`Generate an enhanced description for the following skill:

Name: %s
Display Name: %s
Description: %s
Category: %s
Tags: %v
Language: %s
README: %s

Please provide a well-structured enhanced description.`,
		skill.Name, skill.DisplayName, skill.Description,
		skill.Category, skill.Tags, skill.Language, truncateString(skill.Readme, 3000))

	enhanced, err := s.llmClient.Chat(systemPrompt, userPrompt)
	if err != nil {
		return "", fmt.Errorf("llm enhance: %w", err)
	}

	return enhanced, nil
}

func (s *VectorService) SearchVectors(ctx context.Context, query string, topK int) ([]milvus.SearchResult, error) {
	queryVec, err := s.embedder.Embed(query)
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}

	results, err := s.milvusCli.Search(ctx, queryVec, topK)
	if err != nil {
		return nil, fmt.Errorf("milvus search: %w", err)
	}

	return results, nil
}

func (s *VectorService) GetEmbeddingCount(ctx context.Context) (int64, error) {
	return s.embRepo.CountByModel(s.embedder.Model())
}

func (s *VectorService) StartWorker(ctx context.Context) {
	s.worker.Start(ctx)
}

func (s *VectorService) StopWorker() {
	s.worker.Stop()
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
