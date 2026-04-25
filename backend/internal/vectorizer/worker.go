package vectorizer

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/hpds/skill-hub/internal/client/embedding"
	"github.com/hpds/skill-hub/internal/milvus"
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/internal/repository"
	"github.com/hpds/skill-hub/pkg/logger"
)

type Worker struct {
	embedder    *embedding.Client
	milvusCli   *milvus.Client
	skillRepo   *repository.SkillRepo
	embRepo     *repository.EmbeddingRepo
	concurrency int
	queue       chan int64
	workers     sync.WaitGroup
	stopCh      chan struct{}
}

func NewWorker(
	embedder *embedding.Client,
	milvusCli *milvus.Client,
	skillRepo *repository.SkillRepo,
	embRepo *repository.EmbeddingRepo,
	concurrency int,
) *Worker {
	return &Worker{
		embedder:    embedder,
		milvusCli:   milvusCli,
		skillRepo:   skillRepo,
		embRepo:     embRepo,
		concurrency: concurrency,
		queue:       make(chan int64, 100),
		stopCh:      make(chan struct{}),
	}
}

func (w *Worker) Start(ctx context.Context) {
	for i := 0; i < w.concurrency; i++ {
		w.workers.Add(1)
		go w.runWorker(ctx, i)
	}
	logger.Info("vectorizer worker started", logger.Int("concurrency", w.concurrency))
}

func (w *Worker) Stop() {
	close(w.stopCh)
	w.workers.Wait()
	logger.Info("vectorizer worker stopped")
}

func (w *Worker) Enqueue(skillID int64) {
	select {
	case w.queue <- skillID:
	default:
		logger.Warn("vectorizer queue full, dropping", logger.Int64("skill_id", skillID))
	}
}

func (w *Worker) runWorker(ctx context.Context, id int) {
	defer w.workers.Done()
	for {
		select {
		case <-w.stopCh:
			return
		case <-ctx.Done():
			return
		case skillID := <-w.queue:
			w.processSkill(ctx, skillID)
		}
	}
}

func (w *Worker) processSkill(ctx context.Context, skillID int64) {
	start := time.Now()
	logger.Info("vectorizing skill", logger.Int64("skill_id", skillID), logger.Int("worker_id", 0))

	skill, err := w.skillRepo.GetByID(skillID)
	if err != nil {
		logger.Error("vectorizer get skill", logger.Int64("skill_id", skillID), logger.ErrorField(err))
		return
	}
	if skill == nil {
		logger.Warn("vectorizer skill not found", logger.Int64("skill_id", skillID))
		return
	}

	text := buildEmbeddingText(skill)
	contentHash := fmt.Sprintf("%x", sha256.Sum256([]byte(text)))

	existingEmbs, err := w.embRepo.GetBySkillID(skillID)
	if err == nil && len(existingEmbs) > 0 && existingEmbs[0].ContentHash == contentHash {
		logger.Debug("skill already vectorized, skipped", logger.Int64("skill_id", skillID))
		return
	}

	dims := w.embedder.Dims()
	chunks := chunkText(text, 512)

	var allVectors [][]float32
	var allChunks []string

	for i, chunk := range chunks {
		vec, err := w.embedder.Embed(chunk)
		if err != nil {
			logger.Error("vectorizer embed failed",
				logger.Int64("skill_id", skillID),
				logger.Int("chunk", i),
				logger.ErrorField(err))
			return
		}

		emb := &model.SkillEmbedding{
			SkillID:     skillID,
			Vector:      vec,
			ModelName:   w.embedder.Model(),
			ContentHash: contentHash,
			ChunkIndex:  i,
			ChunkText:   chunk,
		}

		if err := w.embRepo.Upsert(emb); err != nil {
			logger.Error("vectorizer save embedding",
				logger.Int64("skill_id", skillID),
				logger.ErrorField(err))
			continue
		}

		allVectors = append(allVectors, vec)
		allChunks = append(allChunks, chunk)
	}

	if len(allVectors) > 0 {
		ids := make([]int64, len(allVectors))
		for i := range allVectors {
			ids[i] = skillID
		}

		if err := w.milvusCli.BatchInsert(ctx, ids, allVectors, allChunks, w.embedder.Model()); err != nil {
			logger.Error("vectorizer milvus insert",
				logger.Int64("skill_id", skillID),
				logger.ErrorField(err))
		}

		logger.Info("skill vectorized successfully",
			logger.Int64("skill_id", skillID),
			logger.Int("chunks", len(allVectors)),
			logger.Int("dims", dims),
			logger.Duration("duration", time.Since(start)))
	}
}

func buildEmbeddingText(skill *model.Skill) string {
	text := skill.Name
	if skill.DisplayName != "" {
		text += " " + skill.DisplayName
	}
	if skill.Description != "" {
		text += " " + skill.Description
	}
	if skill.Readme != "" {
		maxReadme := 2000
		if len(skill.Readme) > maxReadme {
			text += " " + skill.Readme[:maxReadme]
		} else {
			text += " " + skill.Readme
		}
	}
	return text
}

func chunkText(text string, chunkSize int) []string {
	runes := []rune(text)
	if len(runes) <= chunkSize {
		return []string{text}
	}

	var chunks []string
	for i := 0; i < len(runes); i += chunkSize {
		end := i + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[i:end]))
	}
	return chunks
}
