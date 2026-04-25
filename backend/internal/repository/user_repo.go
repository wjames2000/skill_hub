package repository

import (
	"github.com/hpds/skill-hub/internal/model"
	"xorm.io/xorm"
)

type UserRepo struct {
	db *xorm.Engine
}

func NewUserRepo(db *xorm.Engine) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *model.User) error {
	_, err := r.db.Insert(user)
	return err
}

func (r *UserRepo) GetByID(id int64) (*model.User, error) {
	var user model.User
	has, err := r.db.ID(id).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &user, nil
}

func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	has, err := r.db.Where("username = ?", username).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	has, err := r.db.Where("email = ?", email).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &user, nil
}

func (r *UserRepo) GetByGitHubID(githubID string) (*model.User, error) {
	var user model.User
	has, err := r.db.Where("github_id = ?", githubID).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &user, nil
}

func (r *UserRepo) Update(user *model.User) error {
	_, err := r.db.ID(user.ID).AllCols().Update(user)
	return err
}

func (r *UserRepo) List(page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	total, err := r.db.Limit(pageSize, (page-1)*pageSize).Desc("created_at").FindAndCount(&users)
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepo) UpdateLastLogin(id int64) error {
	_, err := r.db.ID(id).Cols("last_login_at").Update(&model.User{})
	return err
}
