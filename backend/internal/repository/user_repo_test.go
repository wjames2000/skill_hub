package repository

import (
	"context"
	"testing"

	"github.com/hpds/skill-hub/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUserRepo_CreateAndFind(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.User{}, &model.APIKey{})
	repo := NewUserRepo(engine)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: string(hash),
		Role:         "user",
		Status:       1,
	}

	err := repo.Create(ctx, user)
	require.NoError(t, err)
	assert.Greater(t, user.ID, int64(0))
}

func TestUserRepo_FindByUsername(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.User{})
	repo := NewUserRepo(engine)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	_, _ = engine.Insert(&model.User{Username: "findme", Email: "find@test.com", PasswordHash: string(hash), Role: "user", Status: 1})

	t.Run("find existing user", func(t *testing.T) {
		user, err := repo.FindByUsername(ctx, "findme")
		require.NoError(t, err)
		assert.Equal(t, "findme", user.Username)
	})

	t.Run("find non-existing user", func(t *testing.T) {
		user, err := repo.FindByUsername(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepo_FindByEmail(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.User{})
	repo := NewUserRepo(engine)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	_, _ = engine.Insert(&model.User{Username: "emailuser", Email: "email@test.com", PasswordHash: string(hash), Role: "user", Status: 1})

	t.Run("find existing email", func(t *testing.T) {
		user, err := repo.FindByEmail(ctx, "email@test.com")
		require.NoError(t, err)
		assert.Equal(t, "emailuser", user.Username)
	})

	t.Run("find non-existing email", func(t *testing.T) {
		user, err := repo.FindByEmail(ctx, "no@test.com")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepo_FindByUsernameOrEmail(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.User{})
	repo := NewUserRepo(engine)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	_, _ = engine.Insert(&model.User{Username: "mixed", Email: "mixed@test.com", PasswordHash: string(hash), Role: "user", Status: 1})

	t.Run("find by username", func(t *testing.T) {
		user, err := repo.FindByUsernameOrEmail(ctx, "mixed")
		require.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("find by email", func(t *testing.T) {
		user, err := repo.FindByUsernameOrEmail(ctx, "mixed@test.com")
		require.NoError(t, err)
		assert.NotNil(t, user)
	})

	t.Run("not found", func(t *testing.T) {
		user, err := repo.FindByUsernameOrEmail(ctx, "noone")
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepo_List(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.User{})
	repo := NewUserRepo(engine)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	_, _ = engine.Insert(&model.User{Username: "u1", Email: "u1@t.com", PasswordHash: string(hash), Role: "user", Status: 1})
	_, _ = engine.Insert(&model.User{Username: "u2", Email: "u2@t.com", PasswordHash: string(hash), Role: "admin", Status: 1})

	users, total, err := repo.List(ctx, "", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, len(users))
}

func TestUserRepo_Update(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.User{})
	repo := NewUserRepo(engine)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	_, _ = engine.Insert(&model.User{Username: "updateuser", Email: "update@t.com", PasswordHash: string(hash), Role: "user", Status: 1})

	user, _ := repo.FindByUsername(ctx, "updateuser")
	user.Role = "admin"
	err := repo.Update(ctx, user)
	require.NoError(t, err)

	updated, _ := repo.FindByUsername(ctx, "updateuser")
	assert.Equal(t, "admin", updated.Role)
}

func TestUserRepo_Delete(t *testing.T) {
	engine := newTestEngine()
	_ = engine.Sync2(&model.User{})
	repo := NewUserRepo(engine)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
	_, _ = engine.Insert(&model.User{Username: "deleteuser", Email: "delete@t.com", PasswordHash: string(hash), Role: "user", Status: 1})

	user, _ := repo.FindByUsername(ctx, "deleteuser")
	err := repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	deleted, err := repo.FindByUsername(ctx, "deleteuser")
	assert.Error(t, err)
	assert.Nil(t, deleted)
}
