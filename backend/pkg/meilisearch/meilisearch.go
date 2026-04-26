package meilisearch

import (
	"fmt"

	"github.com/hpds/skill-hub/pkg/logger"
	"github.com/meilisearch/meilisearch-go"
)

type Client struct {
	meilisearch.ServiceManager
}

func New(host, apiKey string) (*Client, error) {
	cli := meilisearch.New(host, meilisearch.WithAPIKey(apiKey))
	_, err := cli.Health()
	if err != nil {
		return nil, fmt.Errorf("meilisearch health: %w", err)
	}
	return &Client{cli}, nil
}

func (c *Client) CreateIndex(uid, primaryKey string) (*meilisearch.TaskInfo, error) {
	task, err := c.ServiceManager.CreateIndex(&meilisearch.IndexConfig{
		Uid:        uid,
		PrimaryKey: primaryKey,
	})
	if err != nil {
		return task, err
	}
	// Configure filterable attributes for the skills index
	attrs := []interface{}{"category", "scan_passed", "status", "tags"}
	_, err = c.Index(uid).UpdateFilterableAttributes(&attrs)
	return task, err
}

func (c *Client) EnsureIndex(uid, primaryKey string) error {
	_, err := c.ServiceManager.GetIndex(uid)
	if err == nil {
		// Index exists, ensure filterable attributes are configured
		attrs := []interface{}{"category", "scan_passed", "status", "tags"}
		_, setErr := c.Index(uid).UpdateFilterableAttributes(&attrs)
		if setErr != nil {
			logger.Warn("failed to update filterable attributes", logger.String("error", setErr.Error()))
		}
		return nil
	}
	_, createErr := c.CreateIndex(uid, primaryKey)
	return createErr
}

func (c *Client) AddDocuments(uid string, docs interface{}) (*meilisearch.TaskInfo, error) {
	return c.Index(uid).AddDocuments(docs, nil)
}

func (c *Client) UpdateDocuments(uid string, docs interface{}) (*meilisearch.TaskInfo, error) {
	return c.Index(uid).UpdateDocuments(docs, nil)
}

func (c *Client) DeleteDocuments(uid string, ids []string) (*meilisearch.TaskInfo, error) {
	return c.Index(uid).DeleteDocuments(ids, nil)
}

func (c *Client) Search(uid string, query string, limit int64, filter ...string) (*meilisearch.SearchResponse, error) {
	req := &meilisearch.SearchRequest{
		Limit: limit,
	}
	if len(filter) > 0 && filter[0] != "" {
		req.Filter = filter[0]
	}
	return c.Index(uid).Search(query, req)
}

func (c *Client) DeleteIndex(uid string) (*meilisearch.TaskInfo, error) {
	return c.ServiceManager.DeleteIndex(uid)
}
