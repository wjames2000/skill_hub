package repository

import (
	"github.com/hpds/skill-hub/internal/model"
	"xorm.io/xorm"
)

type CategoryTreeNode struct {
	ID       int64               `json:"id"`
	Name     string              `json:"name"`
	Slug     string              `json:"slug"`
	Icon     string              `json:"icon"`
	Count    int64               `json:"count"`
	Children []*CategoryTreeNode `json:"children"`
}

type CategoryRepo struct {
	db *xorm.Engine
}

func NewCategoryRepo(db *xorm.Engine) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(cat *model.SkillCategory) error {
	_, err := r.db.Insert(cat)
	return err
}

func (r *CategoryRepo) GetByID(id int64) (*model.SkillCategory, error) {
	var cat model.SkillCategory
	has, err := r.db.ID(id).Get(&cat)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &cat, nil
}

func (r *CategoryRepo) GetBySlug(slug string) (*model.SkillCategory, error) {
	var cat model.SkillCategory
	has, err := r.db.Where("slug = ?", slug).Get(&cat)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &cat, nil
}

func (r *CategoryRepo) List() ([]*model.SkillCategory, error) {
	var cats []*model.SkillCategory
	err := r.db.Asc("sort_order").Find(&cats)
	return cats, err
}

func (r *CategoryRepo) GetTree() ([]*CategoryTreeNode, error) {
	all, err := r.List()
	if err != nil {
		return nil, err
	}

	byParent := make(map[int64][]*model.SkillCategory)
	for _, cat := range all {
		byParent[cat.ParentID] = append(byParent[cat.ParentID], cat)
	}

	var build func(parentID int64) []*CategoryTreeNode
	build = func(parentID int64) []*CategoryTreeNode {
		var nodes []*CategoryTreeNode
		for _, cat := range byParent[parentID] {
			count, _ := r.GetSkillCountByCategory(cat.Name)
			node := &CategoryTreeNode{
				ID:    cat.ID,
				Name:  cat.Name,
				Slug:  cat.Slug,
				Icon:  cat.Icon,
				Count: count,
			}
			node.Children = build(cat.ID)
			nodes = append(nodes, node)
		}
		return nodes
	}

	return build(0), nil
}

func (r *CategoryRepo) Update(cat *model.SkillCategory) error {
	_, err := r.db.ID(cat.ID).AllCols().Update(cat)
	return err
}

func (r *CategoryRepo) Delete(id int64) error {
	_, err := r.db.ID(id).Delete(&model.SkillCategory{})
	return err
}

func (r *CategoryRepo) UpdateSkillCount(id int64, count int) error {
	_, err := r.db.ID(id).Cols("skill_count").Update(&model.SkillCategory{SkillCount: count})
	return err
}

func (r *CategoryRepo) GetSkillCountByCategory(category string) (int64, error) {
	return r.db.Where("category = ? AND status = 1", category).Count(&model.Skill{})
}
