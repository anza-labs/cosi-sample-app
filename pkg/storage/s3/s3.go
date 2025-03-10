// Copyright 2025 anza-labs contributors
// SPDX-License-Identifier: MIT

// import "example.com/pkg/storage/s3"
package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type SecretS3 struct {
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyID"`
	AccessSecretKey string `json:"accessSecretKey"`
}

type Client struct {
	s3cli      *minio.Client
	bucketName string
}

func New(bucketName string, s3secret SecretS3, ssl bool) (*Client, error) {
	s3cli, err := minio.New(s3secret.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3secret.AccessKeyID, s3secret.AccessSecretKey, ""),
		Region: s3secret.Region,
		Secure: ssl,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %w", err)
	}

	return &Client{
		s3cli:      s3cli,
		bucketName: bucketName,
	}, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.s3cli.RemoveObject(ctx, c.bucketName, key, minio.RemoveObjectOptions{})
}

func (c *Client) Get(ctx context.Context, key string, wr io.Writer) error {
	obj, err := c.s3cli.GetObject(ctx, c.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	_, err = io.Copy(wr, obj)
	return err
}

func (c *Client) Put(ctx context.Context, key string, data io.Reader, size int64) error {
	_, err := c.s3cli.PutObject(ctx, c.bucketName, key, data, size, minio.PutObjectOptions{})
	return err
}
