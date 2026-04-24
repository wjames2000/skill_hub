package repository

import (
	"testing"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbeddingRepo_Create(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SkillEmbedding{}, &model.Skill{})
	repo := NewEmbeddingRepo(engine)

	emb := &model.SkillEmbedding{
		SkillID:     1,
		ModelName:   "test-model",
		ContentHash: "hash123",
		ChunkIndex:  0,
		ChunkText:   "test chunk",
	}

	err := repo.Create(emb)
	require.NoError(t, err)
	assert.Greater(t, emb.ID, int64(0))
}

func TestEmbeddingRepo_Upsert(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SkillEmbedding{})
	repo := NewEmbeddingRepo(engine)

	emb := &model.SkillEmbedding{
		SkillID:     1,
		ChunkIndex:  0,
		ModelName:   "test-model",
		ContentHash: "hash1",
		ChunkText:   "original",
	}

	t.Run("insert on new embedding", func(t *testing.T) {
		err := repo.Upsert(emb)
		require.NoError(t, err)
		assert.Greater(t, emb.ID, int64(0))
	})

	t.Run("update on existing embedding", func(t *testing.T) {
		emb.ChunkText = "updated"
		err := repo.Upsert(emb)
		require.NoError(t, err)
		assert.Greater(t, emb.ID, int64(0))

		embs, _ := repo.GetBySkillID(1)
		require.Equal(t, 1, len(embs))
		assert.Equal(t, "updated", embs[0].ChunkText)
	})
}

func TestEmbeddingRepo_GetBySkillID(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SkillEmbedding{})
	repo := NewEmbeddingRepo(engine)

	_, _ = engine.Insert(&model.SkillEmbedding{SkillID: 1, ChunkIndex: 0, ChunkText: "chunk0", ModelName: "m"})
	_, _ = engine.Insert(&model.SkillEmbedding{SkillID: 1, ChunkIndex: 1, ChunkText: "chunk1", ModelName: "m"})
	_, _ = engine.Insert(&model.SkillEmbedding{SkillID: 2, ChunkIndex: 0, ChunkText: "other", ModelName: "m"})

	embs, err := repo.GetBySkillID(1)
	require.NoError(t, err)
	assert.Equal(t, 2, len(embs))
	assert.Equal(t, "chunk0", embs[0].ChunkText)
	assert.Equal(t, "chunk1", embs[1].ChunkText)
}

func TestEmbeddingRepo_DeleteBySkillID(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SkillEmbedding{})
	repo := NewEmbeddingRepo(engine)

	_, _ = engine.Insert(&model.SkillEmbedding{SkillID: 1, ChunkIndex: 0, ChunkText: "delete me", ModelName: "m"})

	err := repo.DeleteBySkillID(1)
	require.NoError(t, err)

	embs, _ := repo.GetBySkillID(1)
	assert.Empty(t, embs)
}

func TestEmbeddingRepo_CountByModel(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SkillEmbedding{})
	repo := NewEmbeddingRepo(engine)

	_, _ = engine.Insert(&model.SkillEmbedding{SkillID: 1, ChunkIndex: 0, ModelName: "m1", ChunkText: "a"})
	_, _ = engine.Insert(&model.SkillEmbedding{SkillID: 2, ChunkIndex: 0, ModelName: "m1", ChunkText: "b"})
	_, _ = engine.Insert(&model.SkillEmbedding{SkillID: 3, ChunkIndex: 0, ModelName: "m2", ChunkText: "c"})

	count, err := repo.CountByModel("m1")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestEmbeddingRepo_GetSkillIDsNeedingVectorization(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.SkillEmbedding{}, &model.Skill{})
	repo := NewEmbeddingRepo(engine)

	// Skills without embeddings should appear
	_, _ = engine.Insert(&model.Skill{ID: 1, Name: "s1", Status: model.SkillStatusActive})
	_, _ = engine.Insert(&model.Skill{ID: 2, Name: "s2", Status: model.SkillStatusActive})
	_, _ = engine.Insert(&model.Skill{ID: 3, Name: "s3", Status: model.SkillStatusActive})
	// Skill 4 is inactive, should not appear
	_, _ = engine.Insert(&model.Skill{ID: 4, Name: "s4", Status: model.SkillStatusInactive})
	// Add embedding for skill 2
	_, _ = engine.Insert(&model.SkillEmbedding{SkillID: 2, ChunkIndex: 0, ModelName: "m", ChunkText: "t"})

	ids, err := repo.GetSkillIDsNeedingVectorization(10)
	require.NoError(t, err)
	assert.Contains(t, ids, int64(1))
	assert.Contains(t, ids, int64(3))
	assert.NotContains(t, ids, int64(2)) // has embedding
	assert.NotContains(t, ids, int64(4)) // inactive
}
