package syncer

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type SkillMetadata struct {
	Name         string                 `yaml:"name"`
	DisplayName  string                 `yaml:"display_name"`
	Version      string                 `yaml:"version"`
	Description  string                 `yaml:"description"`
	Author       string                 `yaml:"author"`
	Category     string                 `yaml:"category"`
	Tags         []string               `yaml:"tags"`
	Icon         string                 `yaml:"icon"`
	Homepage     string                 `yaml:"homepage"`
	License      string                 `yaml:"license"`
	Language     string                 `yaml:"language"`
	AI           AIMetadata             `yaml:"ai"`
	Config       map[string]interface{} `yaml:"config"`
	Dependencies []string               `yaml:"dependencies"`
	Install      InstallConfig          `yaml:"install"`
	Extra        map[string]interface{} `yaml:"extra"`
}

type AIMetadata struct {
	Compatible  []string `yaml:"compatible"`
	Recommended bool     `yaml:"recommended"`
	Version     string   `yaml:"version"`
}

type InstallConfig struct {
	Method   string `yaml:"method"`
	Command  string `yaml:"command"`
	Path     string `yaml:"path"`
	Requires string `yaml:"requires"`
}

type ParsedSkill struct {
	RawContent string
	Metadata   *SkillMetadata
	Body       string
}

func ParseSKILLMD(content string) (*ParsedSkill, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, nil
	}

	content = strings.ReplaceAll(content, "\r\n", "\n")

	if !strings.HasPrefix(content, "---") {
		return &ParsedSkill{
			RawContent: content,
			Body:       content,
		}, nil
	}

	rest := content[3:]
	endIdx := strings.Index(rest, "\n---")
	if endIdx < 0 {
		endIdx = strings.Index(rest, "\n---\n")
	}
	if endIdx < 0 {
		endIdx = strings.Index(rest, "---")
		if endIdx >= 0 && endIdx < 10 {
			rest = rest[endIdx+3:]
		}
		return &ParsedSkill{
			RawContent: content,
			Body:       strings.TrimSpace(rest),
		}, nil
	}

	yamlStr := rest[:endIdx]
	body := ""
	if endIdx+5 < len(rest) {
		body = strings.TrimSpace(rest[endIdx+5:])
	} else if endIdx+4 < len(rest) {
		body = strings.TrimSpace(rest[endIdx+4:])
	}

	var metadata SkillMetadata
	if err := yaml.Unmarshal([]byte(yamlStr), &metadata); err != nil {
		return &ParsedSkill{
			RawContent: content,
			Body:       strings.TrimSpace(rest),
		}, nil
	}

	if metadata.Version == "" {
		metadata.Version = "1.0.0"
	}

	return &ParsedSkill{
		RawContent: content,
		Metadata:   &metadata,
		Body:       body,
	}, nil
}

func DetectSkillFile(files []string) string {
	for _, f := range files {
		name := strings.ToLower(f)
		if name == "skill.md" {
			return f
		}
	}

	for _, f := range files {
		if strings.EqualFold(f, "skill.md") ||
			strings.EqualFold(f, "skill.yaml") ||
			strings.EqualFold(f, "skill.yml") {
			return f
		}
	}

	return ""
}
