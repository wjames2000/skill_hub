package syncer

import (
	"context"
	"fmt"
	"strings"
	"time"

	githubclient "github.com/hpds/skill-hub/internal/client/github"
	"github.com/hpds/skill-hub/pkg/logger"
)

type DiscoveredRepo struct {
	Owner    string
	Name     string
	FullName string
	Stars    int
	Source   string
}

type DiscoveryStrategy interface {
	Name() string
	Discover(ctx context.Context, since time.Time) ([]DiscoveredRepo, error)
}

type TopicDiscovery struct {
	client     *githubclient.Client
	topics     []string
	maxPerPage int
}

func NewTopicDiscovery(client *githubclient.Client, topics []string) *TopicDiscovery {
	return &TopicDiscovery{
		client:     client,
		topics:     topics,
		maxPerPage: 100,
	}
}

func (d *TopicDiscovery) Name() string { return "topic" }

func (d *TopicDiscovery) Discover(ctx context.Context, since time.Time) ([]DiscoveredRepo, error) {
	var allRepos []DiscoveredRepo
	seen := make(map[string]bool)

	for _, topic := range d.topics {
		query := fmt.Sprintf("topic:%s sort:stars-desc", topic)
		page := 1
		const maxGitHubPages = 10
		for page <= maxGitHubPages {
			select {
			case <-ctx.Done():
				return allRepos, ctx.Err()
			default:
			}

			repos, total, err := d.client.SearchRepos(ctx, query, page)
			if err != nil {
				logger.Warn("topic search failed",
					logger.String("topic", topic),
					logger.Int("page", page),
					logger.String("error", err.Error()))
				break
			}

			if len(repos) == 0 {
				break
			}

			for _, repo := range repos {
				key := repo.FullName
				if seen[key] {
					continue
				}
				seen[key] = true

				if !since.IsZero() {
					t, err := time.Parse(time.RFC3339, repo.UpdatedAt)
					if err == nil && t.Before(since) {
						continue
					}
				}

				allRepos = append(allRepos, DiscoveredRepo{
					Owner:    repo.Owner,
					Name:     repo.Name,
					FullName: repo.FullName,
					Stars:    repo.Stars,
					Source:   fmt.Sprintf("topic:%s", topic),
				})
			}

			if page*d.maxPerPage >= total || len(repos) < d.maxPerPage {
				break
			}
			page++
		}
	}

	return allRepos, nil
}

type PathDiscovery struct {
	client     *githubclient.Client
	repos      []string
	maxPerPage int
}

func NewPathDiscovery(client *githubclient.Client, repos []string) *PathDiscovery {
	return &PathDiscovery{
		client:     client,
		repos:      repos,
		maxPerPage: 100,
	}
}

func (d *PathDiscovery) Name() string { return "path" }

func (d *PathDiscovery) Discover(ctx context.Context, since time.Time) ([]DiscoveredRepo, error) {
	var allRepos []DiscoveredRepo
	seen := make(map[string]bool)

	for _, path := range d.repos {
		parts := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)
		if len(parts) < 1 {
			continue
		}

		query := fmt.Sprintf("path:%s sort:stars-desc", path)
		page := 1
		for {
			select {
			case <-ctx.Done():
				return allRepos, ctx.Err()
			default:
			}

			repos, total, err := d.client.SearchRepos(ctx, query, page)
			if err != nil {
				logger.Warn("path search failed",
					logger.String("path", path),
					logger.Int("page", page),
					logger.String("error", err.Error()))
				break
			}

			if len(repos) == 0 {
				break
			}

			for _, repo := range repos {
				key := repo.FullName
				if seen[key] {
					continue
				}
				seen[key] = true

				if !since.IsZero() {
					t, err := time.Parse(time.RFC3339, repo.UpdatedAt)
					if err == nil && t.Before(since) {
						continue
					}
				}

				allRepos = append(allRepos, DiscoveredRepo{
					Owner:    repo.Owner,
					Name:     repo.Name,
					FullName: repo.FullName,
					Stars:    repo.Stars,
					Source:   fmt.Sprintf("path:%s", path),
				})
			}

			if page*d.maxPerPage >= total || len(repos) < d.maxPerPage {
				break
			}
			page++
		}
	}

	return allRepos, nil
}

type AwesomeDiscovery struct {
	client *githubclient.Client
}

func NewAwesomeDiscovery(client *githubclient.Client) *AwesomeDiscovery {
	return &AwesomeDiscovery{
		client: client,
	}
}

func (d *AwesomeDiscovery) Name() string { return "awesome" }

func (d *AwesomeDiscovery) Discover(ctx context.Context, since time.Time) ([]DiscoveredRepo, error) {
	awesomeLists := []string{
		"awesome-ai-agents",
		"awesome-ai-tools",
		"awesome-chatgpt-plugins",
		"awesome-copilot",
		"awesome-claude",
		"awesome-gpt",
		"awesome-llm",
	}

	var allRepos []DiscoveredRepo
	seen := make(map[string]bool)

	for _, awesomeName := range awesomeLists {
		logger.Info("searching awesome list", logger.String("name", awesomeName))
		repos, _, err := d.client.SearchRepos(ctx, awesomeName, 1)
		if err != nil {
			logger.Warn("awesome search failed",
				logger.String("name", awesomeName),
				logger.String("error", err.Error()))
			continue
		}

		for _, repo := range repos {
			key := repo.FullName
			if seen[key] {
				continue
			}
			seen[key] = true

			allRepos = append(allRepos, DiscoveredRepo{
				Owner:    repo.Owner,
				Name:     repo.Name,
				FullName: repo.FullName,
				Stars:    repo.Stars,
				Source:   fmt.Sprintf("awesome:%s", awesomeName),
			})
		}
	}

	return allRepos, nil
}

func (d *PathDiscovery) getMaxPerPage() int { return d.maxPerPage }

type KnownRepoDiscovery struct {
	client *githubclient.Client
	repos  []string
}

func NewKnownRepoDiscovery(client *githubclient.Client, repos []string) *KnownRepoDiscovery {
	return &KnownRepoDiscovery{
		client: client,
		repos:  repos,
	}
}

func (d *KnownRepoDiscovery) Name() string { return "known" }

func (d *KnownRepoDiscovery) Discover(ctx context.Context, since time.Time) ([]DiscoveredRepo, error) {
	var allRepos []DiscoveredRepo

	for _, fullName := range d.repos {
		fullName = strings.TrimSpace(fullName)
		if fullName == "" {
			continue
		}

		parts := strings.SplitN(fullName, "/", 2)
		if len(parts) != 2 {
			logger.Warn("invalid known repo format, expected owner/repo",
				logger.String("repo", fullName))
			continue
		}

		owner, name := parts[0], parts[1]

		repoInfo, err := d.client.GetRepo(ctx, owner, name)
		if err != nil {
			logger.Warn("fetch known repo failed",
				logger.String("repo", fullName),
				logger.String("error", err.Error()))
			continue
		}

		if repoInfo.Archived {
			logger.Debug("skipping archived known repo", logger.String("repo", fullName))
			continue
		}

		allRepos = append(allRepos, DiscoveredRepo{
			Owner:    repoInfo.Owner,
			Name:     repoInfo.Name,
			FullName: repoInfo.FullName,
			Stars:    repoInfo.Stars,
			Source:   fmt.Sprintf("known:%s", fullName),
		})

		logger.Info("known repo discovered",
			logger.String("repo", fullName),
			logger.Int("stars", repoInfo.Stars))
	}

	return allRepos, nil
}

type DiscoveryManager struct {
	strategies []DiscoveryStrategy
}

func NewDiscoveryManager(strategies ...DiscoveryStrategy) *DiscoveryManager {
	return &DiscoveryManager{
		strategies: strategies,
	}
}

func (dm *DiscoveryManager) DiscoverAll(ctx context.Context, since time.Time) ([]DiscoveredRepo, error) {
	seen := make(map[string]bool)
	var all []DiscoveredRepo

	for _, strategy := range dm.strategies {
		select {
		case <-ctx.Done():
			return all, ctx.Err()
		default:
		}

		logger.Info("running discovery strategy",
			logger.String("strategy", strategy.Name()))

		repos, err := strategy.Discover(ctx, since)
		if err != nil {
			logger.Warn("discovery strategy errored",
				logger.String("strategy", strategy.Name()),
				logger.String("error", err.Error()))
			continue
		}

		for _, repo := range repos {
			if seen[repo.FullName] {
				continue
			}
			seen[repo.FullName] = true
			all = append(all, repo)
		}

		logger.Info("discovery strategy completed",
			logger.String("strategy", strategy.Name()),
			logger.Int("found", len(repos)),
			logger.Int("total_unique", len(all)))
	}

	if len(all) == 0 {
		logger.Warn("discovery returned zero repos. Check that GitHub tokens are configured in config.yaml")
	}

	return all, nil
}
