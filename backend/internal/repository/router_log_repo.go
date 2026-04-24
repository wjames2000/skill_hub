package repository

import (
	"github.com/hpds/skill-hub/internal/model"
	"xorm.io/xorm"
)

type RouterLogRepo struct {
	db *xorm.Engine
}

func NewRouterLogRepo(db *xorm.Engine) *RouterLogRepo {
	return &RouterLogRepo{db: db}
}

func (r *RouterLogRepo) Create(log *model.RouterLog) error {
	_, err := r.db.Insert(log)
	return err
}

func (r *RouterLogRepo) GetByID(id int64) (*model.RouterLog, error) {
	var log model.RouterLog
	has, err := r.db.ID(id).Get(&log)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &log, nil
}

func (r *RouterLogRepo) List(page, pageSize int) ([]*model.RouterLog, int64, error) {
	var logs []*model.RouterLog
	total, err := r.db.Desc("created_at").Limit(pageSize, (page-1)*pageSize).FindAndCount(&logs)
	if err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (r *RouterLogRepo) GetBySessionID(sessionID string) ([]*model.RouterLog, error) {
	var logs []*model.RouterLog
	err := r.db.Where("session_id = ?", sessionID).OrderBy("created_at").Find(&logs)
	return logs, err
}

func (r *RouterLogRepo) UpdateFeedback(id int64, score int, comment string) error {
	_, err := r.db.ID(id).Cols("feedback_score", "feedback_comment").Update(&model.RouterLog{
		FeedbackScore:   score,
		FeedbackComment: comment,
	})
	return err
}

func (r *RouterLogRepo) UpdateExecution(id int64, result string, duration int) error {
	_, err := r.db.ID(id).Cols("is_executed", "execute_result", "execute_duration").Update(&model.RouterLog{
		IsExecuted:      true,
		ExecuteResult:   result,
		ExecuteDuration: duration,
	})
	return err
}
