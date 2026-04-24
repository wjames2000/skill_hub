package repository

import (
	"github.com/hpds/skill-hub/internal/model"
	"xorm.io/xorm"
)

type APIKeyRepo struct {
	db *xorm.Engine
}

func NewAPIKeyRepo(db *xorm.Engine) *APIKeyRepo {
	return &APIKeyRepo{db: db}
}

func (r *APIKeyRepo) Create(key *model.APIKey) error {
	_, err := r.db.Insert(key)
	return err
}

func (r *APIKeyRepo) GetByID(id int64) (*model.APIKey, error) {
	var key model.APIKey
	has, err := r.db.ID(id).Get(&key)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &key, nil
}

func (r *APIKeyRepo) GetByKey(keyStr string) (*model.APIKey, error) {
	var key model.APIKey
	has, err := r.db.Where("key = ?", keyStr).Get(&key)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &key, nil
}

func (r *APIKeyRepo) ListByUser(userID int64) ([]*model.APIKey, error) {
	var keys []*model.APIKey
	err := r.db.Where("user_id = ?", userID).Desc("created_at").Find(&keys)
	return keys, err
}

func (r *APIKeyRepo) Revoke(id int64) error {
	_, err := r.db.ID(id).Cols("is_revoked").Update(&model.APIKey{IsRevoked: true})
	return err
}

func (r *APIKeyRepo) UpdateLastUsed(id int64) error {
	_, err := r.db.ID(id).Cols("last_used_at").Update(&model.APIKey{})
	return err
}
