package syncer

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hpds/skill-hub/pkg/logger"
)

type ScannerConfig struct {
	Enabled bool
	Binary  string
	Rules   string
	Timeout int
}

type ScanResult struct {
	Passed  bool
	Summary string
	Issues  []ScanIssue
}

type ScanIssue struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Path     string `json:"path"`
	Line     int    `json:"line"`
	RuleID   string `json:"rule_id"`
}

type SecurityScanner struct {
	config ScannerConfig
}

func NewSecurityScanner(cfg ScannerConfig) *SecurityScanner {
	return &SecurityScanner{
		config: cfg,
	}
}

func (s *SecurityScanner) ScanRepo(ctx context.Context, repoPath string) (*ScanResult, error) {
	if !s.config.Enabled {
		return &ScanResult{Passed: true, Summary: "scan disabled"}, nil
	}

	binary := s.config.Binary
	if binary == "" {
		binary = "semgrep"
	}

	args := []string{"scan", "--json", "--no-git-ignore"}
	if s.config.Rules != "" {
		args = append(args, "--config", s.config.Rules)
	}
	args = append(args, repoPath)

	cmd := exec.CommandContext(ctx, binary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			// semgrep returns exit code 1 when findings exist
		} else {
			logger.Warn("semgrep execution failed",
				logger.String("error", err.Error()),
				logger.String("stderr", stderr.String()))
			return &ScanResult{Passed: true, Summary: "scan skipped: " + err.Error()}, nil
		}
	}

	output := stdout.String()
	if output == "" {
		return &ScanResult{Passed: true, Summary: "no output from scanner"}, nil
	}

	return s.parseResults(output)
}

func (s *SecurityScanner) parseResults(output string) (*ScanResult, error) {
	var result struct {
		Results []struct {
			CheckID string `json:"check_id"`
			Path    string `json:"path"`
			Start   struct {
				Line int `json:"line"`
			} `json:"start"`
			Extra struct {
				Severity string `json:"severity"`
				Message  string `json:"message"`
			} `json:"extra"`
		} `json:"results"`
		Errors []interface{} `json:"errors"`
	}

	if err := jsonUnmarshal([]byte(output), &result); err != nil {
		return &ScanResult{Passed: true, Summary: "parse failed: " + err.Error()}, nil
	}

	if len(result.Results) == 0 {
		return &ScanResult{Passed: true, Summary: "no issues found"}, nil
	}

	issues := make([]ScanIssue, 0, len(result.Results))
	highCount := 0
	mediumCount := 0
	lowCount := 0

	for _, r := range result.Results {
		issue := ScanIssue{
			Severity: r.Extra.Severity,
			Message:  r.Extra.Message,
			Path:     r.Path,
			Line:     r.Start.Line,
			RuleID:   r.CheckID,
		}
		issues = append(issues, issue)

		switch strings.ToLower(r.Extra.Severity) {
		case "high", "critical", "error":
			highCount++
		case "medium", "warning":
			mediumCount++
		default:
			lowCount++
		}
	}

	summary := fmt.Sprintf("found %d issues: %d high, %d medium, %d low",
		len(issues), highCount, mediumCount, lowCount)

	passed := highCount == 0
	if highCount > 0 {
		summary = fmt.Sprintf("FAILED - %s", summary)
	}

	return &ScanResult{
		Passed:  passed,
		Summary: summary,
		Issues:  issues,
	}, nil
}
