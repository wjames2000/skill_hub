package llm

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
	baseURL     string
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
	httpCli     *http.Client
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

type chatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}

func New(baseURL, apiKey, model string, maxTokens int, temperature float64) *Client {
	baseURL = strings.TrimRight(baseURL, "/")
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL += "/v1"
	}
	return &Client{
		baseURL:     baseURL,
		apiKey:      apiKey,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
		httpCli:     &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *Client) Chat(systemPrompt, userPrompt string) (string, error) {
	messages := []chatMessage{}
	if systemPrompt != "" {
		messages = append(messages, chatMessage{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, chatMessage{Role: "user", Content: userPrompt})

	req := chatRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   c.maxTokens,
		Temperature: c.temperature,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("llm marshal: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("llm request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpCli.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("llm call: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("llm read: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("llm api error: status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("llm unmarshal: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("llm model error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("llm: no choices returned")
	}

	content := chatResp.Choices[0].Message.Content

	if chatResp.Usage != nil {
		logger.Debug("llm chat completed",
			logger.String("model", c.model),
			logger.Int("prompt_tokens", chatResp.Usage.PromptTokens),
			logger.Int("completion_tokens", chatResp.Usage.CompletionTokens))
	}

	return content, nil
}

func (c *Client) Model() string {
	return c.model
}
