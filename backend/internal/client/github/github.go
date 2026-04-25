package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hpds/skill-hub/pkg/logger"
)

type Client struct {
	tokens       []string
	currentIdx   int
	mu           sync.Mutex
	baseURL      string
	httpClient   *http.Client
	maxPerPage   int
	requestDelay time.Duration
}

type RepoInfo struct {
	ID            int64    `json:"id"`
	FullName      string   `json:"full_name"`
	Owner         string   `json:"owner"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	DefaultBranch string   `json:"default_branch"`
	Homepage      string   `json:"homepage"`
	Language      string   `json:"language"`
	Topics        []string `json:"topics"`
	Stars         int      `json:"stargazers_count"`
	Forks         int      `json:"forks_count"`
	OpenIssues    int      `json:"open_issues_count"`
	License       string   `json:"license"`
	Archived      bool     `json:"archived"`
	AvatarURL     string   `json:"avatar_url"`
	CloneURL      string   `json:"clone_url"`
	UpdatedAt     string   `json:"updated_at"`
}

type FileContent struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	SHA     string `json:"sha"`
	Content string `json:"content"`
	Size    int    `json:"size"`
	Type    string `json:"type"`
	HTMLURL string `json:"html_url"`
}

type RateLimitInfo struct {
	Remaining int
	ResetTime time.Time
	Limit     int
}

func New(tokens []string, maxPerPage int, requestDelay int) *Client {
	if len(tokens) == 0 {
		tokens = []string{""}
	}
	if maxPerPage <= 0 {
		maxPerPage = 100
	}
	return &Client{
		tokens:       tokens,
		currentIdx:   0,
		baseURL:      "https://api.github.com",
		httpClient:   &http.Client{Timeout: 120 * time.Second},
		maxPerPage:   maxPerPage,
		requestDelay: time.Duration(requestDelay) * time.Millisecond,
	}
}

func (c *Client) nextToken() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	token := c.tokens[c.currentIdx]
	c.currentIdx = (c.currentIdx + 1) % len(c.tokens)
	return token
}

func (c *Client) rotateOnRateLimit() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.currentIdx = (c.currentIdx + 1) % len(c.tokens)
}

func (c *Client) doRequest(ctx context.Context, method, path string, params url.Values, body io.Reader) (*http.Response, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}

	var resp *http.Response
	var lastErr error

	token := c.nextToken()

	maxAttempts := len(c.tokens) * 2
	if maxAttempts < 5 {
		maxAttempts = 5
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
		if err != nil {
			return nil, fmt.Errorf("new request: %w", err)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("User-Agent", "SkillHub/1.0")

		time.Sleep(c.requestDelay)

		resp, err = c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			backoff := time.Duration(1<<uint(attempt)) * time.Second
			logger.Warn("github request failed, retrying with backoff",
				logger.String("path", path),
				logger.Int("attempt", attempt+1),
				logger.Duration("backoff", backoff),
				logger.String("error", err.Error()))
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, fmt.Errorf("request cancelled during backoff: %w", ctx.Err())
			}
			token = c.nextToken()
			continue
		}

		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
			remaining := resp.Header.Get("X-RateLimit-Remaining")
			if remaining == "0" || resp.StatusCode == http.StatusUnauthorized {
				resp.Body.Close()
				logger.Warn("github rate limit or unauthorized, rotating token",
					logger.String("remaining", remaining),
					logger.Int("attempt", attempt))
				token = c.nextToken()
				if attempt >= len(c.tokens)-1 {
					time.Sleep(10 * time.Second)
				}
				continue
			}
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}

		resp.Body.Close()
		lastErr = fmt.Errorf("unexpected status: %d", resp.StatusCode)
		token = c.nextToken()
	}

	return nil, fmt.Errorf("all tokens exhausted: %w", lastErr)
}

func parsePagination(resp *http.Response) (nextPage int, lastPage int) {
	link := resp.Header.Get("Link")
	if link == "" {
		return 0, 0
	}
	for _, part := range strings.Split(link, ",") {
		part = strings.TrimSpace(part)
		if strings.Contains(part, `rel="next"`) {
			nextPage = extractPageNum(part)
		}
		if strings.Contains(part, `rel="last"`) {
			lastPage = extractPageNum(part)
		}
	}
	return nextPage, lastPage
}

func extractPageNum(linkPart string) int {
	start := strings.Index(linkPart, "page=")
	if start < 0 {
		return 0
	}
	start += 5
	end := strings.IndexAny(linkPart[start:], ">;&")
	if end < 0 {
		end = len(linkPart) - start
	}
	num, _ := strconv.Atoi(linkPart[start : start+end])
	return num
}

func (c *Client) SearchRepos(ctx context.Context, query string, page int) ([]RepoInfo, int, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("per_page", strconv.Itoa(c.maxPerPage))
	params.Set("page", strconv.Itoa(page))
	params.Set("sort", "stars")
	params.Set("order", "desc")

	resp, err := c.doRequest(ctx, "GET", "/search/repositories", params, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("search repos: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read body: %w", err)
	}

	var searchResult struct {
		TotalCount int `json:"total_count"`
		Items      []struct {
			ID            int64    `json:"id"`
			FullName      string   `json:"full_name"`
			Name          string   `json:"name"`
			Description   string   `json:"description"`
			DefaultBranch string   `json:"default_branch"`
			Homepage      string   `json:"homepage"`
			Language      string   `json:"language"`
			Topics        []string `json:"topics"`
			Stars         int      `json:"stargazers_count"`
			Forks         int      `json:"forks_count"`
			OpenIssues    int      `json:"open_issues_count"`
			License       *struct {
				SPDXID string `json:"spdx_id"`
			} `json:"license"`
			Archived bool `json:"archived"`
			Owner    *struct {
				AvatarURL string `json:"avatar_url"`
				Login     string `json:"login"`
			} `json:"owner"`
			CloneURL  string `json:"clone_url"`
			UpdatedAt string `json:"updated_at"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &searchResult); err != nil {
		return nil, 0, fmt.Errorf("unmarshal: %w", err)
	}

	nextPage, _ := parsePagination(resp)
	_ = nextPage

	repos := make([]RepoInfo, 0, len(searchResult.Items))
	for _, item := range searchResult.Items {
		license := ""
		if item.License != nil {
			license = item.License.SPDXID
		}
		avatarURL := ""
		if item.Owner != nil {
			avatarURL = item.Owner.AvatarURL
		}
		repos = append(repos, RepoInfo{
			ID:            item.ID,
			FullName:      item.FullName,
			Owner:         item.Owner.Login,
			Name:          item.Name,
			Description:   item.Description,
			DefaultBranch: item.DefaultBranch,
			Homepage:      item.Homepage,
			Language:      item.Language,
			Topics:        item.Topics,
			Stars:         item.Stars,
			Forks:         item.Forks,
			OpenIssues:    item.OpenIssues,
			License:       license,
			Archived:      item.Archived,
			AvatarURL:     avatarURL,
			CloneURL:      item.CloneURL,
			UpdatedAt:     item.UpdatedAt,
		})
	}

	return repos, searchResult.TotalCount, nil
}

func (c *Client) GetRepoContent(ctx context.Context, owner, repo, path, branch string) (*FileContent, error) {
	apiPath := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path)
	params := url.Values{}
	if branch != "" {
		params.Set("ref", branch)
	}

	resp, err := c.doRequest(ctx, "GET", apiPath, params, nil)
	if err != nil {
		return nil, fmt.Errorf("get content: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result []FileContent
	if err := json.Unmarshal(body, &result); err != nil {
		var single FileContent
		if err2 := json.Unmarshal(body, &single); err2 != nil {
			return nil, fmt.Errorf("unmarshal content: %w", err)
		}
		if single.Content != "" {
			single.Content = decodeBase64(single.Content)
		}
		return &single, nil
	}

	for _, item := range result {
		if item.Name == "SKILL.md" || item.Name == "skill.md" {
			item.Content = decodeBase64(item.Content)
			return &item, nil
		}
	}

	return nil, nil
}

func (c *Client) GetRepo(ctx context.Context, owner, repo string) (*RepoInfo, error) {
	apiPath := fmt.Sprintf("/repos/%s/%s", owner, repo)

	resp, err := c.doRequest(ctx, "GET", apiPath, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get repo: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var data struct {
		ID            int64    `json:"id"`
		FullName      string   `json:"full_name"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		DefaultBranch string   `json:"default_branch"`
		Homepage      string   `json:"homepage"`
		Language      string   `json:"language"`
		Topics        []string `json:"topics"`
		Stars         int      `json:"stargazers_count"`
		Forks         int      `json:"forks_count"`
		OpenIssues    int      `json:"open_issues_count"`
		Archived      bool     `json:"archived"`
		Owner         *struct {
			AvatarURL string `json:"avatar_url"`
			Login     string `json:"login"`
		} `json:"owner"`
		CloneURL  string `json:"clone_url"`
		UpdatedAt string `json:"updated_at"`
		License   *struct {
			SPDXID string `json:"spdx_id"`
		} `json:"license"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	license := ""
	if data.License != nil {
		license = data.License.SPDXID
	}
	avatarURL := ""
	if data.Owner != nil {
		avatarURL = data.Owner.AvatarURL
	}

	return &RepoInfo{
		ID:            data.ID,
		FullName:      data.FullName,
		Owner:         data.Owner.Login,
		Name:          data.Name,
		Description:   data.Description,
		DefaultBranch: data.DefaultBranch,
		Homepage:      data.Homepage,
		Language:      data.Language,
		Topics:        data.Topics,
		Stars:         data.Stars,
		Forks:         data.Forks,
		OpenIssues:    data.OpenIssues,
		License:       license,
		Archived:      data.Archived,
		AvatarURL:     avatarURL,
		CloneURL:      data.CloneURL,
		UpdatedAt:     data.UpdatedAt,
	}, nil
}

func (c *Client) GetReadme(ctx context.Context, owner, repo, branch string) (string, error) {
	apiPath := fmt.Sprintf("/repos/%s/%s/readme", owner, repo)
	params := url.Values{}
	if branch != "" {
		params.Set("ref", branch)
	}

	resp, err := c.doRequest(ctx, "GET", apiPath, params, nil)
	if err != nil {
		return "", fmt.Errorf("get readme: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var data struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("unmarshal readme: %w", err)
	}

	return decodeBase64(data.Content), nil
}

func (c *Client) ListRepos(ctx context.Context, owner string, page int) ([]RepoInfo, int, error) {
	apiPath := fmt.Sprintf("/users/%s/repos", owner)
	params := url.Values{}
	params.Set("per_page", strconv.Itoa(c.maxPerPage))
	params.Set("page", strconv.Itoa(page))
	params.Set("type", "public")
	params.Set("sort", "updated")

	resp, err := c.doRequest(ctx, "GET", apiPath, params, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("list repos: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read body: %w", err)
	}

	var items []struct {
		ID            int64    `json:"id"`
		FullName      string   `json:"full_name"`
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		DefaultBranch string   `json:"default_branch"`
		Homepage      string   `json:"homepage"`
		Language      string   `json:"language"`
		Topics        []string `json:"topics"`
		Stars         int      `json:"stargazers_count"`
		Forks         int      `json:"forks_count"`
		OpenIssues    int      `json:"open_issues_count"`
		Archived      bool     `json:"archived"`
		Owner         *struct {
			AvatarURL string `json:"avatar_url"`
			Login     string `json:"login"`
		} `json:"owner"`
		CloneURL  string `json:"clone_url"`
		UpdatedAt string `json:"updated_at"`
		License   *struct {
			SPDXID string `json:"spdx_id"`
		} `json:"license"`
	}

	if err := json.Unmarshal(body, &items); err != nil {
		return nil, 0, fmt.Errorf("unmarshal: %w", err)
	}

	nextPage, lastPage := parsePagination(resp)
	totalPages := lastPage
	if totalPages == 0 {
		totalPages = nextPage
	}
	totalCount := len(items)
	if totalPages > page {
		totalCount = totalPages * c.maxPerPage
	}

	repos := make([]RepoInfo, 0, len(items))
	for _, item := range items {
		license := ""
		if item.License != nil {
			license = item.License.SPDXID
		}
		avatarURL := ""
		if item.Owner != nil {
			avatarURL = item.Owner.AvatarURL
		}
		repos = append(repos, RepoInfo{
			ID:            item.ID,
			FullName:      item.FullName,
			Owner:         item.Owner.Login,
			Name:          item.Name,
			Description:   item.Description,
			DefaultBranch: item.DefaultBranch,
			Homepage:      item.Homepage,
			Language:      item.Language,
			Topics:        item.Topics,
			Stars:         item.Stars,
			Forks:         item.Forks,
			OpenIssues:    item.OpenIssues,
			License:       license,
			Archived:      item.Archived,
			AvatarURL:     avatarURL,
			CloneURL:      item.CloneURL,
			UpdatedAt:     item.UpdatedAt,
		})
	}

	return repos, totalCount, nil
}

func (c *Client) GetContents(ctx context.Context, owner, repo, path, branch string) ([]FileContent, error) {
	apiPath := fmt.Sprintf("/repos/%s/%s/contents/%s", owner, repo, path)
	params := url.Values{}
	if branch != "" {
		params.Set("ref", branch)
	}

	resp, err := c.doRequest(ctx, "GET", apiPath, params, nil)
	if err != nil {
		return nil, fmt.Errorf("get contents: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result []FileContent
	if err := json.Unmarshal(body, &result); err != nil {
		var single FileContent
		if err2 := json.Unmarshal(body, &single); err2 != nil {
			return nil, fmt.Errorf("unmarshal contents: %w", err)
		}
		single.Content = decodeBase64(single.Content)
		return []FileContent{single}, nil
	}

	for i := range result {
		result[i].Content = decodeBase64(result[i].Content)
	}

	return result, nil
}
