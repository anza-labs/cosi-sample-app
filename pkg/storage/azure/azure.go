// Copyright 2025 anza-labs contributors
// SPDX-License-Identifier: MIT

// import "example.com/pkg/storage/azure"
package azure

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type SecretAzure struct {
	AccessToken     string    `json:"accessToken"`
	ExpiryTimestamp time.Time `json:"expiryTimeStamp"`
}

type Client struct {
	azCli         *azblob.Client
	containerName string
}

func New(containerName string, azureSecret SecretAzure) (*Client, error) {
	azCli, err := azblob.NewClientWithNoCredential(azureSecret.AccessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %w", err)
	}

	return &Client{
		azCli:         azCli,
		containerName: containerName,
	}, nil
}

func (c *Client) Delete(ctx context.Context, blobName string) error {
	_, err := c.azCli.DeleteBlob(ctx, c.containerName, blobName, nil)
	return err
}

func (c *Client) Get(ctx context.Context, blobName string, wr io.Writer) error {
	stream, err := c.azCli.DownloadStream(ctx, c.containerName, blobName, nil)
	if err != nil {
		return fmt.Errorf("unable to get download stream: %w", err)
	}
	_, err = io.Copy(wr, stream.Body)
	return err
}

func (c *Client) Put(ctx context.Context, blobName string, data io.Reader, size int64) error {
	_, err := c.azCli.UploadStream(ctx, c.containerName, blobName, data, nil)
	return err
}
