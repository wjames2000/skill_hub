package repository

import (
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/pkg/logger"
)

type seedCategory struct {
	Name     string
	Slug     string
	Icon     string
	ParentID int64
	Order    int
	Children []seedCategory
}

var categoryTree = []seedCategory{
	{
		Name: "内容创作类", Slug: "content-creation", Icon: "palette", Order: 1,
		Children: []seedCategory{
			{Name: "文案撰写", Slug: "copywriting", Icon: "edit_note", Order: 1},
			{Name: "视觉设计", Slug: "visual-design", Icon: "design_services", Order: 2},
			{Name: "多媒体制作", Slug: "multimedia-production", Icon: "movie", Order: 3},
		},
	},
	{
		Name: "信息处理类", Slug: "information-processing", Icon: "database", Order: 2,
		Children: []seedCategory{
			{Name: "数据采集", Slug: "data-collection", Icon: "travel_explore", Order: 1},
			{Name: "数据分析", Slug: "data-analysis", Icon: "analytics", Order: 2},
			{Name: "数据清洗", Slug: "data-cleaning", Icon: "cleaning_services", Order: 3},
			{Name: "文档处理", Slug: "document-processing", Icon: "description", Order: 4},
			{Name: "翻译", Slug: "translation", Icon: "translate", Order: 5},
		},
	},
	{
		Name: "软件工程类", Slug: "software-engineering", Icon: "code", Order: 3,
		Children: []seedCategory{
			{
				Name: "代码质量", Slug: "code-quality", Icon: "verified", Order: 1,
				Children: []seedCategory{
					{Name: "PR 检查", Slug: "pr-review", Icon: "merge", Order: 1},
					{Name: "安全扫描", Slug: "security-scan", Icon: "shield", Order: 2},
					{Name: "规范校验", Slug: "lint-check", Icon: "checklist", Order: 3},
				},
			},
			{
				Name: "模板脚手架", Slug: "templates", Icon: "folder_copy", Order: 2,
				Children: []seedCategory{
					{Name: "项目初始化", Slug: "project-init", Icon: "rocket_launch", Order: 1},
					{Name: "文档生成", Slug: "doc-generation", Icon: "auto_stories", Order: 2},
					{Name: "代码生成", Slug: "code-generation", Icon: "smartphone", Order: 3},
				},
			},
			{
				Name: "部署运维", Slug: "deployment", Icon: "cloud_sync", Order: 3,
				Children: []seedCategory{
					{Name: "上线检查", Slug: "release-check", Icon: "check_circle", Order: 1},
					{Name: "回滚管理", Slug: "rollback", Icon: "undo", Order: 2},
					{Name: "环境验证", Slug: "env-verification", Icon: "verified_user", Order: 3},
				},
			},
		},
	},
	{
		Name: "团队协作类", Slug: "team-collaboration", Icon: "groups", Order: 4,
		Children: []seedCategory{
			{
				Name: "流程自动化", Slug: "process-automation", Icon: "sync_alt", Order: 1,
				Children: []seedCategory{
					{Name: "需求评审", Slug: "requirement-review", Icon: "rate_review", Order: 1},
					{Name: "发版检查", Slug: "release-review", Icon: "deployed_code", Order: 2},
					{Name: "周报流程", Slug: "weekly-report", Icon: "calendar_month", Order: 3},
				},
			},
			{
				Name: "知识管理", Slug: "knowledge-management", Icon: "menu_book", Order: 2,
				Children: []seedCategory{
					{Name: "内部知识库", Slug: "knowledge-base", Icon: "library_books", Order: 1},
					{Name: "SDK 指南", Slug: "sdk-guide", Icon: "code_blocks", Order: 2},
					{Name: "客诉 SOP", Slug: "complaint-sop", Icon: "support", Order: 3},
				},
			},
			{
				Name: "项目管理", Slug: "project-management", Icon: "track_changes", Order: 3,
				Children: []seedCategory{
					{Name: "目标追踪", Slug: "goal-tracking", Icon: "flag", Order: 1},
					{Name: "日志生成", Slug: "log-generation", Icon: "summarize", Order: 2},
				},
			},
		},
	},
	{
		Name: "基础设施类", Slug: "infrastructure", Icon: "dns", Order: 5,
		Children: []seedCategory{
			{Name: "资源巡检", Slug: "resource-inspection", Icon: "search_insights", Order: 1},
			{Name: "集群管理", Slug: "cluster-management", Icon: "hub", Order: 2},
			{Name: "故障排查", Slug: "troubleshooting", Icon: "troubleshoot", Order: 3},
			{Name: "环境修复", Slug: "env-repair", Icon: "build", Order: 4},
		},
	},
	{
		Name: "AI 智能体类", Slug: "ai-agent", Icon: "smart_toy", Order: 6,
		Children: []seedCategory{
			{Name: "原子型技能", Slug: "atomic-skill", Icon: "widgets", Order: 1},
			{Name: "工作流型技能", Slug: "workflow-skill", Icon: "account_tree", Order: 2},
			{Name: "专属型技能", Slug: "dedicated-skill", Icon: "star", Order: 3},
		},
	},
	{
		Name: "参考型资料", Slug: "reference-materials", Icon: "book", Order: 7,
		Children: []seedCategory{
			{Name: "规则说明", Slug: "rules", Icon: "gavel", Order: 1},
			{Name: "SDK 使用方法", Slug: "sdk-usage", Icon: "manual", Order: 2},
		},
	},
}

func (r *CategoryRepo) SeedDefaultCategories() error {
	count, err := r.db.Count(&model.SkillCategory{})
	if err != nil {
		return err
	}
	if count > 0 {
		logger.Info("categories already seeded, skipping", logger.Int64("count", count))
		return nil
	}

	var insert func(list []seedCategory, parentID int64) error
	insert = func(list []seedCategory, parentID int64) error {
		for _, sc := range list {
			cat := &model.SkillCategory{
				Name:      sc.Name,
				Slug:      sc.Slug,
				Icon:      sc.Icon,
				ParentID:  parentID,
				SortOrder: sc.Order,
			}
			if err := r.Create(cat); err != nil {
				logger.Warn("seed category failed",
					logger.String("name", sc.Name),
					logger.String("error", err.Error()))
				continue
			}
			if len(sc.Children) > 0 {
				if err := insert(sc.Children, cat.ID); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if err := insert(categoryTree, 0); err != nil {
		return err
	}

	logger.Info("default categories seeded successfully")
	return nil
}
