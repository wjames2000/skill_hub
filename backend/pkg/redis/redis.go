package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	*goredis.Client
}

func New(addr, password string, db int) (*Client, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     20,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &Client{rdb}, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.Client.Del(ctx, keys...).Err()
}

func (c *Client) Exists(ctx context.Context, keys ...string) (bool, error) {
	n, err := c.Client.Exists(ctx, keys...).Result()
	return n > 0, err
}

func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.Client.Incr(ctx, key).Result()
}

func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.Client.Expire(ctx, key, expiration).Err()
}

func (c *Client) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return c.Client.Eval(ctx, script, keys, args...).Result()
}
