package reranker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	apiKey  string
	model   string
	httpCli *http.Client
}

type rerankRequest struct {
	Model     string   `json:"model"`
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
	TopN      int      `json:"top_n,omitempty"`
}

type rerankResponse struct {
	Model   string         `json:"model"`
	Results []rerankResult `json:"results"`
	Usage   *struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage,omitempty"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type rerankResult struct {
	Index          int     `json:"index"`
	RelevanceScore float64 `json:"relevance_score"`
	Document       string  `json:"document,omitempty"`
}

type Result struct {
	Index    int
	Score    float64
	Document string
}

func New(baseURL, apiKey, model string) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL += "/v1"
	}
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		httpCli: &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *Client) Rerank(query string, documents []string, topN int) ([]Result, error) {
	if topN <= 0 || topN > len(documents) {
		topN = len(documents)
	}

	req := rerankRequest{
		Model:     c.model,
		Query:     query,
		Documents: documents,
		TopN:      topN,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("rerank marshal: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/rerank", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("rerank request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpCli.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("rerank call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("rerank read: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rerank api error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var rerankResp rerankResponse
	if err := json.Unmarshal(respBody, &rerankResp); err != nil {
		return nil, fmt.Errorf("rerank unmarshal: %w", err)
	}

	if rerankResp.Error != nil {
		return nil, fmt.Errorf("rerank model error: %s", rerankResp.Error.Message)
	}

	results := make([]Result, len(rerankResp.Results))
	for i, r := range rerankResp.Results {
		results[i] = Result{
			Index:    r.Index,
			Score:    r.RelevanceScore,
			Document: r.Document,
		}
	}

	return results, nil
}

func (c *Client) Model() string {
	return c.model
}
