package repository

import (
	"fmt"
	"time"

	"github.com/hpds/skill-hub/internal/model"
	"xorm.io/xorm"
)

type SkillRepo struct {
	db *xorm.Engine
}

func NewSkillRepo(db *xorm.Engine) *SkillRepo {
	return &SkillRepo{db: db}
}

func (r *SkillRepo) Create(skill *model.Skill) error {
	_, err := r.db.Insert(skill)
	return err
}

func (r *SkillRepo) Update(skill *model.Skill) error {
	_, err := r.db.ID(skill.ID).AllCols().Update(skill)
	return err
}

func (r *SkillRepo) GetByID(id int64) (*model.Skill, error) {
	var skill model.Skill
	has, err := r.db.ID(id).Get(&skill)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &skill, nil
}

func (r *SkillRepo) GetByRepo(repo string) (*model.Skill, error) {
	var skill model.Skill
	has, err := r.db.Where("repository = ?", repo).Get(&skill)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &skill, nil
}

func (r *SkillRepo) Upsert(skill *model.Skill) (bool, error) {
	existing, err := r.GetByRepo(skill.Repository)
	if err != nil {
		return false, fmt.Errorf("get existing: %w", err)
	}

	if existing == nil {
		skill.LastSyncAt = time.Now()
		skill.Status = model.SkillStatusActive
		_, err = r.db.Insert(skill)
		if err != nil {
			return false, fmt.Errorf("insert: %w", err)
		}
		return true, nil
	}

	skill.ID = existing.ID
	skill.LastSyncAt = time.Now()
	skill.CreatedAt = existing.CreatedAt
	skill.Installs = existing.Installs
	skill.Score = existing.Score
	_, err = r.db.ID(existing.ID).AllCols().Update(skill)
	if err != nil {
		return false, fmt.Errorf("update: %w", err)
	}
	return false, nil
}

func (r *SkillRepo) List(page, pageSize int, status int) ([]*model.Skill, int64, error) {
	var skills []*model.Skill
	sess := r.db.Where("1=1")
	if status >= 0 {
		sess = sess.Where("status = ?", status)
	}
	total, err := sess.Limit(pageSize, (page-1)*pageSize).Desc("stars").FindAndCount(&skills)
	if err != nil {
		return nil, 0, err
	}
	return skills, total, nil
}

func (r *SkillRepo) Delete(id int64) error {
	_, err := r.db.ID(id).Delete(&model.Skill{})
	return err
}

func (r *SkillRepo) SetStatus(id int64, status int) error {
	_, err := r.db.ID(id).Cols("status").Update(&model.Skill{Status: status})
	return err
}

func (r *SkillRepo) UpdateScanResult(id int64, passed bool, report string) error {
	_, err := r.db.ID(id).Cols("scan_passed", "scan_report").Update(&model.Skill{
		ScanPassed: passed,
		ScanReport: report,
	})
	return err
}

func (r *SkillRepo) ListQuery() *xorm.Session {
	return r.db.Where("1=1")
}

func (r *SkillRepo) UpdateScore(id int64, score float64) error {
	_, err := r.db.ID(id).Cols("score").Update(&model.Skill{Score: score})
	return err
}

func (r *SkillRepo) IncrementInstalls(id int64) error {
	_, err := r.db.Exec("UPDATE skills SET installs = installs + 1 WHERE id = ?", id)
	return err
}

func (r *SkillRepo) ListByIDs(ids []int64) ([]*model.Skill, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var skills []*model.Skill
	err := r.db.In("id", ids).Find(&skills)
	return skills, err
}

func (r *SkillRepo) SearchByName(name string, limit int) ([]*model.Skill, error) {
	var skills []*model.Skill
	err := r.db.Where("name LIKE ?", "%"+name+"%").
		Or("display_name LIKE ?", "%"+name+"%").
		Limit(limit).Find(&skills)
	if err != nil {
		return nil, err
	}
	return skills, nil
}

type SkillStats struct {
	TotalSkills   int64 `json:"total_skills"`
	ActiveSkills  int64 `json:"active_skills"`
	TotalStars    int64 `json:"total_stars"`
	TotalInstalls int64 `json:"total_installs"`
}

func (r *SkillRepo) GetStats() (*SkillStats, error) {
	stats := &SkillStats{}
	_, err := r.db.SQL("SELECT COUNT(*) as total_skills FROM skills").Get(stats)
	if err != nil {
		return nil, err
	}
	_, err = r.db.SQL("SELECT COUNT(*) as active_skills FROM skills WHERE status = ?", model.SkillStatusActive).Get(stats)
	if err != nil {
		return nil, err
	}
	row := r.db.SQL("SELECT COALESCE(SUM(stars), 0) as total_stars, COALESCE(SUM(installs), 0) as total_installs FROM skills")
	_, err = row.Get(stats)
	return stats, err
}

func (r *SkillRepo) ListIDsNeedingUpdate(since time.Time, limit int) ([]int64, error) {
	var ids []int64
	err := r.db.Table(&model.Skill{}).
		Where("last_sync_at < ? OR last_sync_at IS NULL", since).
		Limit(limit).
		Cols("id").
		Find(&ids)
	return ids, err
}

func (r *SkillRepo) CountByStatus(status int) (int64, error) {
	return r.db.Where("status = ?", status).Count(&model.Skill{})
}
