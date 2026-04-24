package repository

import (
	"github.com/hpds/skill-hub/internal/model"
	"xorm.io/xorm"
)

type ReviewRepo struct {
	db *xorm.Engine
}

func NewReviewRepo(db *xorm.Engine) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) Create(review *model.Review) error {
	_, err := r.db.Insert(review)
	return err
}

func (r *ReviewRepo) GetByID(id int64) (*model.Review, error) {
	var review model.Review
	has, err := r.db.ID(id).Get(&review)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &review, nil
}

func (r *ReviewRepo) ListBySkill(skillID int64, page, pageSize int) ([]*model.Review, int64, error) {
	var reviews []*model.Review
	total, err := r.db.Where("skill_id = ?", skillID).
		Limit(pageSize, (page-1)*pageSize).
		Desc("created_at").
		FindAndCount(&reviews)
	return reviews, total, err
}

func (r *ReviewRepo) ListByUser(userID int64, page, pageSize int) ([]*model.Review, int64, error) {
	var reviews []*model.Review
	total, err := r.db.Where("user_id = ?", userID).
		Limit(pageSize, (page-1)*pageSize).
		Desc("created_at").
		FindAndCount(&reviews)
	return reviews, total, err
}

func (r *ReviewRepo) ListAll(page, pageSize int) ([]*model.Review, int64, error) {
	var reviews []*model.Review
	total, err := r.db.
		Limit(pageSize, (page-1)*pageSize).
		Desc("created_at").
		FindAndCount(&reviews)
	return reviews, total, err
}

func (r *ReviewRepo) Update(review *model.Review) error {
	_, err := r.db.ID(review.ID).AllCols().Update(review)
	return err
}

func (r *ReviewRepo) Delete(id int64) error {
	_, err := r.db.ID(id).Delete(&model.Review{})
	return err
}

func (r *ReviewRepo) GetAvgScoreBySkill(skillID int64) (float64, error) {
	stats := struct {
		Avg float64 `xorm:"avg(score)"`
	}{}
	_, err := r.db.SQL("SELECT COALESCE(AVG(score), 0) as avg FROM reviews WHERE skill_id = ?", skillID).Get(&stats)
	return stats.Avg, err
}

func (r *ReviewRepo) CountBySkill(skillID int64) (int64, error) {
	return r.db.Where("skill_id = ?", skillID).Count(&model.Review{})
}
