package repository

import (
	"math"
	"sort"

	"github.com/hpds/skill-hub/internal/milvus"
	"github.com/hpds/skill-hub/internal/model"
	"xorm.io/xorm"
)

type EmbeddingRepo struct {
	db *xorm.Engine
}

func NewEmbeddingRepo(db *xorm.Engine) *EmbeddingRepo {
	return &EmbeddingRepo{db: db}
}

func (r *EmbeddingRepo) Create(emb *model.SkillEmbedding) error {
	_, err := r.db.Insert(emb)
	return err
}

func (r *EmbeddingRepo) Upsert(emb *model.SkillEmbedding) error {
	existing := &model.SkillEmbedding{}
	has, err := r.db.Where("skill_id = ? AND chunk_index = ?", emb.SkillID, emb.ChunkIndex).Get(existing)
	if err != nil {
		return err
	}
	if has {
		emb.ID = existing.ID
		emb.CreatedAt = existing.CreatedAt
		_, err = r.db.ID(existing.ID).AllCols().Update(emb)
		return err
	}
	_, err = r.db.Insert(emb)
	return err
}

func (r *EmbeddingRepo) GetBySkillID(skillID int64) ([]*model.SkillEmbedding, error) {
	var embs []*model.SkillEmbedding
	err := r.db.Where("skill_id = ?", skillID).OrderBy("chunk_index").Find(&embs)
	return embs, err
}

func (r *EmbeddingRepo) DeleteBySkillID(skillID int64) error {
	_, err := r.db.Where("skill_id = ?", skillID).Delete(&model.SkillEmbedding{})
	return err
}

func (r *EmbeddingRepo) CountByModel(modelName string) (int64, error) {
	return r.db.Where("model_name = ?", modelName).Count(&model.SkillEmbedding{})
}

func (r *EmbeddingRepo) SearchSimilar(queryVec []float32, modelName string, topK int) ([]milvus.SearchResult, error) {
	var embs []*model.SkillEmbedding
	if err := r.db.Where("model_name = ?", modelName).Find(&embs); err != nil {
		return nil, err
	}

	best := make(map[int64]float64)

	for _, emb := range embs {
		score := cosineSimilarity(queryVec, emb.Vector)
		if existing, ok := best[emb.SkillID]; !ok || score > existing {
			best[emb.SkillID] = score
		}
	}

	type scoredSkill struct {
		skillID int64
		score   float64
	}
	var sorted []scoredSkill
	for id, score := range best {
		sorted = append(sorted, scoredSkill{skillID: id, score: score})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].score > sorted[j].score
	})
	if len(sorted) > topK {
		sorted = sorted[:topK]
	}

	results := make([]milvus.SearchResult, len(sorted))
	for i, s := range sorted {
		results[i] = milvus.SearchResult{ID: s.skillID, Score: s.score}
	}
	return results, nil
}

func (r *EmbeddingRepo) GetSkillIDsNeedingVectorization(limit int) ([]int64, error) {
	var ids []int64
	err := r.db.SQL(`SELECT s.id FROM skills s 
		LEFT JOIN skill_embeddings e ON s.id = e.skill_id 
		WHERE e.id IS NULL AND s.status = 1 
		LIMIT ?`, limit).Find(&ids)
	return ids, err
}

func cosineSimilarity(a, b []float32) float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	var dot, normA, normB float64
	for i := 0; i < n; i++ {
		va := float64(a[i])
		vb := float64(b[i])
		dot += va * vb
		normA += va * va
		normB += vb * vb
	}
	denom := math.Sqrt(normA) * math.Sqrt(normB)
	if denom == 0 {
		return 0
	}
	return dot / denom
}
