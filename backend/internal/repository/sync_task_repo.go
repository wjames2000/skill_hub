package repository

import (
	"github.com/hpds/skill-hub/internal/model"
	"xorm.io/xorm"
)

type SyncTaskRepo struct {
	db *xorm.Engine
}

func NewSyncTaskRepo(db *xorm.Engine) *SyncTaskRepo {
	return &SyncTaskRepo{db: db}
}

func (r *SyncTaskRepo) Create(task *model.SyncTask) error {
	_, err := r.db.Insert(task)
	return err
}

func (r *SyncTaskRepo) Update(task *model.SyncTask) error {
	_, err := r.db.ID(task.ID).AllCols().Update(task)
	return err
}

func (r *SyncTaskRepo) GetByID(id int64) (*model.SyncTask, error) {
	var task model.SyncTask
	has, err := r.db.ID(id).Get(&task)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &task, nil
}

func (r *SyncTaskRepo) List(page, pageSize int) ([]*model.SyncTask, int64, error) {
	var tasks []*model.SyncTask
	total, err := r.db.Desc("id").Limit(pageSize, (page-1)*pageSize).FindAndCount(&tasks)
	if err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

func (r *SyncTaskRepo) GetLatestByType(syncType string) (*model.SyncTask, error) {
	var task model.SyncTask
	has, err := r.db.Where("type = ?", syncType).Desc("id").Get(&task)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &task, nil
}

func (r *SyncTaskRepo) GetRunningTask() (*model.SyncTask, error) {
	var task model.SyncTask
	has, err := r.db.Where("status = ?", model.SyncStatusRunning).Get(&task)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &task, nil
}

func (r *SyncTaskRepo) CancelRunning() error {
	_, err := r.db.Where("status = ?", model.SyncStatusRunning).
		Cols("status", "error_message").
		Update(&model.SyncTask{Status: model.SyncStatusCancelled, ErrorMessage: "cancelled by new task"})
	return err
}
