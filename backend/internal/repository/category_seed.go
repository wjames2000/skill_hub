package repository

import (
	"github.com/hpds/skill-hub/internal/model"
	"github.com/hpds/skill-hub/pkg/logger"
)

type seedCategory struct {
	Name     string
	ZhName   string
	EnName   string
	Slug     string
	Icon     string
	ParentID int64
	Order    int
	Children []seedCategory
}

var categoryTree = []seedCategory{
	{
		Name: "内容创作类", ZhName: "内容创作类", EnName: "Content Creation",
		Slug: "content-creation", Icon: "palette", Order: 1,
		Children: []seedCategory{
			{Name: "文案撰写", ZhName: "文案撰写", EnName: "Copywriting", Slug: "copywriting", Icon: "edit_note", Order: 1},
			{Name: "视觉设计", ZhName: "视觉设计", EnName: "Visual Design", Slug: "visual-design", Icon: "design_services", Order: 2},
			{Name: "多媒体制作", ZhName: "多媒体制作", EnName: "Multimedia Production", Slug: "multimedia-production", Icon: "movie", Order: 3},
		},
	},
	{
		Name: "信息处理类", ZhName: "信息处理类", EnName: "Information Processing",
		Slug: "information-processing", Icon: "database", Order: 2,
		Children: []seedCategory{
			{Name: "数据采集", ZhName: "数据采集", EnName: "Data Collection", Slug: "data-collection", Icon: "travel_explore", Order: 1},
			{Name: "数据分析", ZhName: "数据分析", EnName: "Data Analysis", Slug: "data-analysis", Icon: "analytics", Order: 2},
			{Name: "数据清洗", ZhName: "数据清洗", EnName: "Data Cleaning", Slug: "data-cleaning", Icon: "cleaning_services", Order: 3},
			{Name: "文档处理", ZhName: "文档处理", EnName: "Document Processing", Slug: "document-processing", Icon: "description", Order: 4},
			{Name: "翻译", ZhName: "翻译", EnName: "Translation", Slug: "translation", Icon: "translate", Order: 5},
		},
	},
	{
		Name: "软件工程类", ZhName: "软件工程类", EnName: "Software Engineering",
		Slug: "software-engineering", Icon: "code", Order: 3,
		Children: []seedCategory{
			{
				Name: "代码质量", ZhName: "代码质量", EnName: "Code Quality", Slug: "code-quality", Icon: "verified", Order: 1,
				Children: []seedCategory{
					{Name: "PR 检查", ZhName: "PR 检查", EnName: "PR Review", Slug: "pr-review", Icon: "merge", Order: 1},
					{Name: "安全扫描", ZhName: "安全扫描", EnName: "Security Scan", Slug: "security-scan", Icon: "shield", Order: 2},
					{Name: "规范校验", ZhName: "规范校验", EnName: "Lint Check", Slug: "lint-check", Icon: "checklist", Order: 3},
				},
			},
			{
				Name: "模板脚手架", ZhName: "模板脚手架", EnName: "Templates & Scaffolding", Slug: "templates", Icon: "folder_copy", Order: 2,
				Children: []seedCategory{
					{Name: "项目初始化", ZhName: "项目初始化", EnName: "Project Init", Slug: "project-init", Icon: "rocket_launch", Order: 1},
					{Name: "文档生成", ZhName: "文档生成", EnName: "Docs Generation", Slug: "doc-generation", Icon: "auto_stories", Order: 2},
					{Name: "代码生成", ZhName: "代码生成", EnName: "Code Generation", Slug: "code-generation", Icon: "smartphone", Order: 3},
				},
			},
			{
				Name: "部署运维", ZhName: "部署运维", EnName: "Deployment & Operations", Slug: "deployment", Icon: "cloud_sync", Order: 3,
				Children: []seedCategory{
					{Name: "上线检查", ZhName: "上线检查", EnName: "Release Check", Slug: "release-check", Icon: "check_circle", Order: 1},
					{Name: "回滚管理", ZhName: "回滚管理", EnName: "Rollback Management", Slug: "rollback", Icon: "undo", Order: 2},
					{Name: "环境验证", ZhName: "环境验证", EnName: "Env Verification", Slug: "env-verification", Icon: "verified_user", Order: 3},
				},
			},
		},
	},
	{
		Name: "团队协作类", ZhName: "团队协作类", EnName: "Team Collaboration",
		Slug: "team-collaboration", Icon: "groups", Order: 4,
		Children: []seedCategory{
			{
				Name: "流程自动化", ZhName: "流程自动化", EnName: "Process Automation", Slug: "process-automation", Icon: "sync_alt", Order: 1,
				Children: []seedCategory{
					{Name: "需求评审", ZhName: "需求评审", EnName: "Requirement Review", Slug: "requirement-review", Icon: "rate_review", Order: 1},
					{Name: "发版检查", ZhName: "发版检查", EnName: "Release Review", Slug: "release-review", Icon: "deployed_code", Order: 2},
					{Name: "周报流程", ZhName: "周报流程", EnName: "Weekly Report", Slug: "weekly-report", Icon: "calendar_month", Order: 3},
				},
			},
			{
				Name: "知识管理", ZhName: "知识管理", EnName: "Knowledge Management", Slug: "knowledge-management", Icon: "menu_book", Order: 2,
				Children: []seedCategory{
					{Name: "内部知识库", ZhName: "内部知识库", EnName: "Knowledge Base", Slug: "knowledge-base", Icon: "library_books", Order: 1},
					{Name: "SDK 指南", ZhName: "SDK 指南", EnName: "SDK Guide", Slug: "sdk-guide", Icon: "code_blocks", Order: 2},
					{Name: "客诉 SOP", ZhName: "客诉 SOP", EnName: "Complaint SOP", Slug: "complaint-sop", Icon: "support", Order: 3},
				},
			},
			{
				Name: "项目管理", ZhName: "项目管理", EnName: "Project Management", Slug: "project-management", Icon: "track_changes", Order: 3,
				Children: []seedCategory{
					{Name: "目标追踪", ZhName: "目标追踪", EnName: "Goal Tracking", Slug: "goal-tracking", Icon: "flag", Order: 1},
					{Name: "日志生成", ZhName: "日志生成", EnName: "Log Generation", Slug: "log-generation", Icon: "summarize", Order: 2},
				},
			},
		},
	},
	{
		Name: "基础设施类", ZhName: "基础设施类", EnName: "Infrastructure",
		Slug: "infrastructure", Icon: "dns", Order: 5,
		Children: []seedCategory{
			{Name: "资源巡检", ZhName: "资源巡检", EnName: "Resource Inspection", Slug: "resource-inspection", Icon: "search_insights", Order: 1},
			{Name: "集群管理", ZhName: "集群管理", EnName: "Cluster Management", Slug: "cluster-management", Icon: "hub", Order: 2},
			{Name: "故障排查", ZhName: "故障排查", EnName: "Troubleshooting", Slug: "troubleshooting", Icon: "troubleshoot", Order: 3},
			{Name: "环境修复", ZhName: "环境修复", EnName: "Env Repair", Slug: "env-repair", Icon: "build", Order: 4},
		},
	},
	{
		Name: "AI 智能体类", ZhName: "AI 智能体类", EnName: "AI Agents",
		Slug: "ai-agent", Icon: "smart_toy", Order: 6,
		Children: []seedCategory{
			{Name: "原子型技能", ZhName: "原子型技能", EnName: "Atomic Skills", Slug: "atomic-skill", Icon: "widgets", Order: 1},
			{Name: "工作流型技能", ZhName: "工作流型技能", EnName: "Workflow Skills", Slug: "workflow-skill", Icon: "account_tree", Order: 2},
			{Name: "专属型技能", ZhName: "专属型技能", EnName: "Dedicated Skills", Slug: "dedicated-skill", Icon: "star", Order: 3},
		},
	},
	{
		Name: "参考型资料", ZhName: "参考型资料", EnName: "Reference Materials",
		Slug: "reference-materials", Icon: "book", Order: 7,
		Children: []seedCategory{
			{Name: "规则说明", ZhName: "规则说明", EnName: "Rules", Slug: "rules", Icon: "gavel", Order: 1},
			{Name: "SDK 使用方法", ZhName: "SDK 使用方法", EnName: "SDK Usage", Slug: "sdk-usage", Icon: "manual", Order: 2},
		},
	},
}

func (r *CategoryRepo) SeedDefaultCategories() error {
	count, err := r.db.Count(&model.SkillCategory{})
	if err != nil {
		return err
	}

	var upsert func(list []seedCategory, parentID int64) ([]int64, error)
	upsert = func(list []seedCategory, parentID int64) ([]int64, error) {
		var ids []int64
		for _, sc := range list {
			existing, _ := r.GetBySlug(sc.Slug)
			if existing != nil {
				if existing.ZhName == "" || existing.EnName == "" {
					if _, err := r.db.ID(existing.ID).Cols("zh_name", "en_name").Update(&model.SkillCategory{
						ZhName: sc.ZhName,
						EnName: sc.EnName,
					}); err != nil {
						logger.Warn("update category bilingual name failed",
							logger.String("slug", sc.Slug), logger.String("error", err.Error()))
					}
				}
				childIDs, err := upsert(sc.Children, existing.ID)
				if err != nil {
					return ids, err
				}
				ids = append(ids, childIDs...)
			} else {
				cat := &model.SkillCategory{
					Name:      sc.Name,
					ZhName:    sc.ZhName,
					EnName:    sc.EnName,
					Slug:      sc.Slug,
					Icon:      sc.Icon,
					ParentID:  parentID,
					SortOrder: sc.Order,
				}
				if err := r.Create(cat); err != nil {
					logger.Warn("seed category failed",
						logger.String("name", sc.Name), logger.String("error", err.Error()))
					continue
				}
				childIDs, err := upsert(sc.Children, cat.ID)
				if err != nil {
					return ids, err
				}
				ids = append(ids, childIDs...)
			}
		}
		return ids, nil
	}

	_, err = upsert(categoryTree, 0)
	if err != nil {
		return err
	}

	if count == 0 {
		logger.Info("default categories seeded successfully")
	} else {
		logger.Info("default categories bilingual names updated")
	}
	return nil
}
