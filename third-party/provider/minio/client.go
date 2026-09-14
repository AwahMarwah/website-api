package minio

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client    *minio.Client
	bucket    string
	endpoint  string
	useSSL    bool
}

func NewClient() (*Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	bucket := os.Getenv("MINIO_BUCKET")
	if bucket == "" {
		bucket = "website"
	}
	useSSL := strings.TrimSpace(os.Getenv("MINIO_USE_SSL")) == "true"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init minio client: %w", err)
	}

	c := &Client{
		client:   minioClient,
		bucket:   bucket,
		endpoint: endpoint,
		useSSL:   useSSL,
	}

	if err := c.ensureBucket(); err != nil {
		log.Printf("[minio] ensure bucket failed (continuing): %v", err)
	}

	return c, nil
}

func (c *Client) ensureBucket() error {
	ctx := context.Background()
	exists, err := c.client.BucketExists(ctx, c.bucket)
	if err != nil {
		return err
	}
	if !exists {
		if err := c.client.MakeBucket(ctx, c.bucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}
	// bucket public-read supaya URL permanen bisa diakses
	return c.client.SetBucketPolicy(ctx, c.bucket, publicReadPolicy(c.bucket))
}

func (c *Client) UploadFile(objectKey string, body []byte, contentType string) (string, error) {
	reader := bytes.NewReader(body)
	_, err := c.client.PutObject(
		context.Background(),
		c.bucket,
		objectKey,
		reader,
		int64(len(body)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return "", fmt.Errorf("minio upload failed: %w", err)
	}
	return c.PublicURL(objectKey), nil
}

func (c *Client) PublicURL(objectKey string) string {
	scheme := "http"
	if c.useSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, c.endpoint, c.bucket, objectKey)
}

func publicReadPolicy(bucket string) string {
	return fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {"AWS": ["*"]},
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::%s/*"]
    }
  ]
}`, bucket)
}