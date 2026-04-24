package repository

import (
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

func (r *EmbeddingRepo) GetSkillIDsNeedingVectorization(limit int) ([]int64, error) {
	var ids []int64
	err := r.db.SQL(`SELECT s.id FROM skills s 
		LEFT JOIN skill_embeddings e ON s.id = e.skill_id 
		WHERE e.id IS NULL AND s.status = 1 
		LIMIT ?`, limit).Find(&ids)
	return ids, err
}
