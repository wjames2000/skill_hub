package model

import "time"

type SkillEmbedding struct {
	ID          int64     `xorm:"pk autoincr 'id'" json:"id"`
	SkillID     int64     `xorm:"int not null index 'skill_id'" json:"skill_id"`
	Vector      []float32 `xorm:"json 'vector'" json:"vector"`
	ModelName   string    `xorm:"varchar(100) 'model_name'" json:"model_name"`
	ContentHash string    `xorm:"varchar(64) 'content_hash'" json:"content_hash"`
	ChunkIndex  int       `xorm:"int default 0 'chunk_index'" json:"chunk_index"`
	ChunkText   string    `xorm:"text 'chunk_text'" json:"chunk_text"`
	CreatedAt   time.Time `xorm:"created 'created_at'" json:"created_at"`
	UpdatedAt   time.Time `xorm:"updated 'updated_at'" json:"updated_at"`
}

func (s *SkillEmbedding) TableName() string {
	return "skill_embeddings"
}
