package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	*minio.Client
	Bucket string
}

func New(endpoint, accessKey, secretKey string, useSSL bool, bucket string) (*Client, error) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio new: %w", err)
	}
	ctx := context.Background()
	exists, err := cli.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("minio bucket check: %w", err)
	}
	if !exists {
		if err := cli.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio make bucket: %w", err)
		}
	}
	return &Client{cli, bucket}, nil
}

func (c *Client) Upload(ctx context.Context, objectName, contentType string, reader io.Reader, size int64) (minio.UploadInfo, error) {
	return c.Client.PutObject(ctx, c.Bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
}

func (c *Client) Download(ctx context.Context, objectName string) (*minio.Object, error) {
	return c.Client.GetObject(ctx, c.Bucket, objectName, minio.GetObjectOptions{})
}

func (c *Client) Delete(ctx context.Context, objectName string) error {
	return c.Client.RemoveObject(ctx, c.Bucket, objectName, minio.RemoveObjectOptions{})
}

func (c *Client) List(ctx context.Context, prefix string) <-chan minio.ObjectInfo {
	return c.Client.ListObjects(ctx, c.Bucket, minio.ListObjectsOptions{Prefix: prefix})
}
