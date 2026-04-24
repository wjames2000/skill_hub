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
		"ai": "ai", "artificial-intelligence": "ai", "machine-learning": "ai",
		"nlp": "ai", "llm": "ai", "gpt": "ai", "claude": "ai", "copilot": "ai",
		"chatgpt": "ai", "openai": "ai", "agent": "ai", "rag": "ai",

		"developer-tools": "developer-tools", "devtools": "developer-tools",
		"cli": "developer-tools", "terminal": "developer-tools",
		"productivity": "productivity", "automation": "productivity",

		"code-review": "code-review", "testing": "testing",
		"documentation": "documentation", "docs": "documentation",
		"security": "security",

		"database": "database", "data": "database",
		"monitoring": "monitoring", "observability": "monitoring",

		"frontend": "frontend", "react": "frontend", "vue": "frontend",
		"backend": "backend", "api": "backend", "server": "backend",
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

	return "uncategorized"
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
