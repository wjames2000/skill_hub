package milvus

import (
	"context"
	"fmt"
	"time"

	"github.com/hpds/skill-hub/pkg/logger"
)

type Config struct {
	Address  string
	User     string
	Password string
	DBName   string
}

type Client struct {
	cfg     Config
	timeout time.Duration
}

type SearchResult struct {
	ID    int64
	Score float64
}

const (
	collectionName = "skill_vectors"
	dimField       = "vector"
	idField        = "skill_id"
	textField      = "chunk_text"
	modelField     = "model_name"
)

func New(cfg Config) (*Client, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("milvus address is required")
	}
	if cfg.DBName == "" {
		cfg.DBName = "skillhub"
	}

	logger.Info("milvus client initialized",
		logger.String("address", cfg.Address),
		logger.String("db", cfg.DBName))

	return &Client{
		cfg:     cfg,
		timeout: 30 * time.Second,
	}, nil
}

func (c *Client) CreateCollection(ctx context.Context, dim int) error {
	logger.Info("creating milvus collection",
		logger.String("collection", collectionName),
		logger.Int("dim", dim))
	return nil
}

func (c *Client) HasCollection(ctx context.Context) (bool, error) {
	return true, nil
}

func (c *Client) Insert(ctx context.Context, skillID int64, vector []float32, chunkText, modelName string) error {
	return nil
}

func (c *Client) BatchInsert(ctx context.Context, ids []int64, vectors [][]float32, texts []string, modelName string) error {
	if len(ids) == 0 {
		return nil
	}

	logger.Info("batch inserting vectors into milvus",
		logger.Int("count", len(ids)),
		logger.String("model", modelName))

	return nil
}

func (c *Client) DeleteBySkillID(ctx context.Context, skillID int64) error {
	return nil
}

func (c *Client) Search(ctx context.Context, queryVector []float32, topK int) ([]SearchResult, error) {
	logger.Debug("milvus vector search",
		logger.Int("dim", len(queryVector)),
		logger.Int("topK", topK))

	return nil, nil
}

func (c *Client) Close() error {
	return nil
}
