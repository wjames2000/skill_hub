package repository

import (
	"github.com/hpds/skill-hub/internal/model"
	"xorm.io/xorm"
)

type FavoriteRepo struct {
	db *xorm.Engine
}

func NewFavoriteRepo(db *xorm.Engine) *FavoriteRepo {
	return &FavoriteRepo{db: db}
}

func (r *FavoriteRepo) Create(fav *model.Favorite) error {
	_, err := r.db.Insert(fav)
	return err
}

func (r *FavoriteRepo) Delete(userID, skillID int64) error {
	_, err := r.db.Where("user_id = ? AND skill_id = ?", userID, skillID).Delete(&model.Favorite{})
	return err
}

func (r *FavoriteRepo) Exists(userID, skillID int64) (bool, error) {
	return r.db.Where("user_id = ? AND skill_id = ?", userID, skillID).Exist(&model.Favorite{})
}

func (r *FavoriteRepo) ListByUser(userID int64, page, pageSize int) ([]*model.Favorite, int64, error) {
	var favs []*model.Favorite
	total, err := r.db.Where("user_id = ?", userID).
		Limit(pageSize, (page-1)*pageSize).
		Desc("created_at").
		FindAndCount(&favs)
	return favs, total, err
}

func (r *FavoriteRepo) CountByUser(userID int64) (int64, error) {
	return r.db.Where("user_id = ?", userID).Count(&model.Favorite{})
}

func (r *FavoriteRepo) CountBySkill(skillID int64) (int64, error) {
	return r.db.Where("skill_id = ?", skillID).Count(&model.Favorite{})
}
