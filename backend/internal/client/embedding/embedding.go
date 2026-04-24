package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hpds/skill-hub/pkg/logger"
)

type Client struct {
	baseURL string
	apiKey  string
	model   string
	dims    int
	httpCli *http.Client
}

type embeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
	Dims  int      `json:"dimensions,omitempty"`
}

type embeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage *struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage,omitempty"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}

func New(baseURL, apiKey, model string, dims int) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL += "/v1"
	}
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
		dims:    dims,
		httpCli: &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *Client) Embed(text string) ([]float32, error) {
	result, err := c.BatchEmbed([]string{text})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("embedding: empty result")
	}
	return result[0], nil
}

func (c *Client) BatchEmbed(texts []string) ([][]float32, error) {
	req := embeddingRequest{
		Model: c.model,
		Input: texts,
	}
	if c.dims > 0 {
		req.Dims = c.dims
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("embedding marshal: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("embedding request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpCli.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("embedding call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("embedding read: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding api error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var embResp embeddingResponse
	if err := json.Unmarshal(respBody, &embResp); err != nil {
		return nil, fmt.Errorf("embedding unmarshal: %w", err)
	}

	if embResp.Error != nil {
		return nil, fmt.Errorf("embedding model error: %s", embResp.Error.Message)
	}

	result := make([][]float32, len(embResp.Data))
	for i, d := range embResp.Data {
		if d.Index < len(result) {
			result[d.Index] = d.Embedding
		} else {
			result[i] = d.Embedding
		}
	}

	logger.Debug("embedding completed",
		logger.Int("batch_size", len(texts)),
		logger.Int("dims", len(result[0])),
		logger.String("model", c.model))

	return result, nil
}

func (c *Client) Dims() int {
	return c.dims
}

func (c *Client) Model() string {
	return c.model
}
