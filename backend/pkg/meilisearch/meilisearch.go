package meilisearch

import (
	"fmt"

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
	return c.ServiceManager.CreateIndex(&meilisearch.IndexConfig{
		Uid:        uid,
		PrimaryKey: primaryKey,
	})
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

func (c *Client) Search(uid string, query string, limit int64) (*meilisearch.SearchResponse, error) {
	return c.Index(uid).Search(query, &meilisearch.SearchRequest{
		Limit: limit,
	})
}

func (c *Client) DeleteIndex(uid string) (*meilisearch.TaskInfo, error) {
	return c.ServiceManager.DeleteIndex(uid)
}
