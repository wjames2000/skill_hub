package syncer

import (
	"bytes"
	"context"
	"encoding/json"
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

type skillGuardOutput struct {
	ScanTime   string `json:"scan_time"`
	Duration   string `json:"duration"`
	TotalFiles int    `json:"total_files"`
	TotalIssue int    `json:"total_issues"`
	Results    []struct {
		RuleID      string `json:"rule_id"`
		Severity    string `json:"severity"`
		FilePath    string `json:"file_path"`
		LineNumber  int    `json:"line_number"`
		LineContent string `json:"line_content"`
		MatchType   string `json:"match_type"`
	} `json:"results"`
	Summary struct {
		Critical int `json:"critical"`
		High     int `json:"high"`
		Medium   int `json:"medium"`
		Low      int `json:"low"`
	} `json:"summary"`
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
		binary = "skill-guard"
	}

	if _, err := exec.LookPath(binary); err != nil {
		return &ScanResult{Passed: true, Summary: "scan skipped: " + binary + " not found"}, nil
	}

	args := []string{repoPath, "--json"}
	if s.config.Rules != "" {
		args = append(args, "--rules", s.config.Rules)
	}

	cmd := exec.CommandContext(ctx, binary, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 && stdout.Len() > 0 {
				// skill-guard exits with code 1 when findings exist
			} else {
				logger.Warn("skill-guard execution failed",
					logger.String("error", err.Error()),
					logger.String("stderr", stderr.String()))
				return &ScanResult{Passed: true, Summary: "scan skipped: " + err.Error()}, nil
			}
		} else {
			logger.Warn("skill-guard execution failed",
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
	var sgOutput skillGuardOutput
	if err := json.Unmarshal([]byte(output), &sgOutput); err != nil {
		return &ScanResult{Passed: true, Summary: "parse failed: " + err.Error()}, nil
	}

	issues := make([]ScanIssue, 0, len(sgOutput.Results))
	highCount := 0
	mediumCount := 0
	lowCount := 0
	criticalCount := 0

	for _, r := range sgOutput.Results {
		issues = append(issues, ScanIssue{
			Severity: r.Severity,
			Message:  r.LineContent,
			Path:     r.FilePath,
			Line:     r.LineNumber,
			RuleID:   r.RuleID,
		})

		switch strings.ToLower(r.Severity) {
		case "critical":
			criticalCount++
		case "high":
			highCount++
		case "medium":
			mediumCount++
		default:
			lowCount++
		}
	}

	totalFound := criticalCount + highCount + mediumCount + lowCount
	if totalFound == 0 {
		return &ScanResult{Passed: true, Summary: "no issues found"}, nil
	}

	summary := fmt.Sprintf("found %d issues: %d critical, %d high, %d medium, %d low",
		totalFound, criticalCount, highCount, mediumCount, lowCount)

	passed := highCount == 0 && criticalCount == 0
	if !passed {
		summary = fmt.Sprintf("FAILED - %s", summary)
	}

	return &ScanResult{
		Passed:  passed,
		Summary: summary,
		Issues:  issues,
	}, nil
}
