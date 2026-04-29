package syncer

import (
	"time"

	githubclient "github.com/hpds/skill-hub/internal/client/github"
	"github.com/hpds/skill-hub/internal/model"
)

type SkillExtractor struct{}

func NewSkillExtractor() *SkillExtractor {
	return &SkillExtractor{}
}

func (e *SkillExtractor) Extract(repo *githubclient.RepoInfo, parsed *ParsedSkill, readme string) *model.Skill {
	skill := &model.Skill{
		Repository:    repo.FullName,
		RepoOwner:     repo.Owner,
		RepoName:      repo.Name,
		DefaultBranch: repo.DefaultBranch,
		Stars:         repo.Stars,
		Forks:         repo.Forks,
		OpenIssues:    repo.OpenIssues,
		Language:      repo.Language,
		Topics:        repo.Topics,
		License:       repo.License,
		AvatarURL:     repo.AvatarURL,
		Homepage:      repo.Homepage,
		IsArchived:    repo.Archived,
		Status:        model.SkillStatusActive,
		Score:         e.calculateScore(repo.Stars, repo.Forks),
	}

	if repo.Description != "" {
		skill.Description = repo.Description
	}

	if readme != "" {
		skill.Readme = truncateString(readme, 50000)
	}

	if parsed != nil && parsed.Metadata != nil {
		m := parsed.Metadata

		if m.Name != "" {
			skill.Name = m.Name
		} else {
			skill.Name = repo.Name
		}

		if m.DisplayName != "" {
			skill.DisplayName = m.DisplayName
		}

		if m.Version != "" {
			skill.Version = m.Version
		} else {
			skill.Version = "1.0.0"
		}

		if m.Description != "" {
			skill.Description = m.Description
		}

		if m.Author != "" {
			skill.Author = m.Author
		} else {
			skill.Author = repo.Owner
		}

		if m.Category != "" {
			skill.Category = m.Category
		} else {
			skill.Category = e.inferCategory(repo.Topics, repo.Description)
		}

		if len(m.Tags) > 0 {
			skill.Tags = m.Tags
		} else {
			skill.Tags = repo.Topics
		}

		if m.Homepage != "" {
			skill.Homepage = m.Homepage
		}

		if m.License != "" {
			skill.License = m.License
		}

		if m.Language != "" {
			skill.Language = m.Language
		}

		if parsed.Body != "" && skill.Readme == "" {
			skill.Readme = parsed.Body
		}
	} else {
		skill.Name = repo.Name
		skill.Author = repo.Owner
		skill.Category = e.inferCategory(repo.Topics, repo.Description)
		skill.Tags = repo.Topics
		skill.Version = "1.0.0"
	}

	skill.LastSyncAt = time.Now()

	return skill
}

func (e *SkillExtractor) calculateScore(stars, forks int) float64 {
	score := float64(stars)*1.0 + float64(forks)*0.5
	if score > 10000 {
		score = 10000
	}
	return score
}

func (e *SkillExtractor) inferCategory(topics []string, description string) string {
	categoryMap := map[string]string{
		// AI Agent sub-categories
		"ai": "atomic-skill", "artificial-intelligence": "atomic-skill", "machine-learning": "atomic-skill",
		"nlp": "atomic-skill", "llm": "atomic-skill", "gpt": "atomic-skill", "claude": "atomic-skill",
		"openai": "atomic-skill", "rag": "atomic-skill", "prompt": "atomic-skill",
		"agent": "dedicated-skill", "assistant": "dedicated-skill",
		"workflow": "workflow-skill", "pipeline": "workflow-skill",

		// Software Engineering sub-categories
		"code-review": "pr-review", "pr": "pr-review", "merge": "pr-review",
		"testing": "lint-check", "lint": "lint-check", "quality": "code-quality",
		"security-scan": "security-scan", "semgrep": "security-scan",
		"cli": "project-init", "terminal": "project-init",
		"scaffold": "project-init", "boilerplate": "project-init", "template": "project-init",
		"project-init":   "project-init",
		"code-generator": "code-generation", "codegen": "code-generation",
		"doc-generation": "doc-generation", "docs": "doc-generation",
		"deploy": "release-check", "ci-cd": "release-check", "release": "release-check",
		"rollback": "rollback",
		"devops":   "env-verification", "developer-tools": "env-verification", "devtools": "env-verification",

		// Content Creation sub-categories
		"content": "copywriting", "writing": "copywriting",
		"copywriting": "copywriting", "blog": "copywriting",
		"documentation": "doc-generation", "markdown": "doc-generation",
		"visual-design": "visual-design", "design": "visual-design",
		"multimedia": "multimedia-production", "video": "multimedia-production",

		// Information Processing sub-categories
		"data": "data-analysis", "database": "data-analysis",
		"analytics": "data-analysis", "etl": "data-analysis",
		"crawler": "data-collection", "scraper": "data-collection",
		"visualization": "data-analysis", "report": "data-analysis",
		"translation": "translation", "translate": "translation",
		"document-processing": "document-processing", "ocr": "document-processing",

		// Infrastructure sub-categories
		"security": "resource-inspection", "monitoring": "resource-inspection",
		"observability": "resource-inspection",
		"kubernetes":    "cluster-management", "docker": "cluster-management",
		"troubleshooting": "troubleshooting", "diagnostic": "troubleshooting",
		"env-repair": "env-repair", "repair": "env-repair",

		// Team Collaboration sub-categories
		"collaboration": "process-automation",
		"requirement":   "requirement-review", "sprint": "release-review",
		"project-management": "goal-tracking",
		"knowledge":          "knowledge-base", "wiki": "knowledge-base",
		"sdk":       "sdk-guide",
		"complaint": "complaint-sop", "support": "complaint-sop",

		// Reference Materials sub-categories
		"reference": "rules", "tutorial": "rules", "guide": "rules",
		"rule": "rules", "spec": "rules",
		"sdk-usage": "sdk-usage",
	}

	for _, topic := range topics {
		t := toLower(topic)
		if cat, ok := categoryMap[t]; ok {
			return cat
		}
	}

	if description != "" {
		desc := toLower(description)
		for keyword, cat := range categoryMap {
			if contains(desc, keyword) {
				return cat
			}
		}
	}

	return "reference-materials"
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsString(s, substr)
}

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
