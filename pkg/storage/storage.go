// Copyright 2025 anza-labs contributors
// SPDX-License-Identifier: MIT

// import "example.com/pkg/storage"
package storage

import (
	"context"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/anza-labs/cosi-sample-app/pkg/storage/azure"
	"github.com/anza-labs/cosi-sample-app/pkg/storage/s3"
)

type Config struct {
	Spec Spec `json:"spec"`
}

type Spec struct {
	BucketName         string             `json:"bucketName"`
	AuthenticationType string             `json:"authenticationType"`
	Protocols          []string           `json:"protocols"`
	SecretS3           *s3.SecretS3       `json:"secretS3,omitempty"`
	SecretAzure        *azure.SecretAzure `json:"secretAzure,omitempty"`
}

type Storage interface {
	Delete(ctx context.Context, key string) error
	Get(ctx context.Context, key string, wr io.Writer) error
	Put(ctx context.Context, key string, data io.Reader, size int64) error
}

func New(config Config, ssl bool) (Storage, error) {
	if slices.ContainsFunc(config.Spec.Protocols, func(s string) bool { return strings.EqualFold(s, "s3") }) {
		if !strings.EqualFold(config.Spec.AuthenticationType, "key") {
			return nil, fmt.Errorf("invalid authentication type for s3")
		}

		s3secret := config.Spec.SecretS3
		if s3secret == nil {
			return nil, fmt.Errorf("s3 secret missing")
		}

		return s3.New(config.Spec.BucketName, *s3secret, ssl)
	}

	if slices.ContainsFunc(config.Spec.Protocols, func(s string) bool { return strings.EqualFold(s, "azure") }) {
		if !strings.EqualFold(config.Spec.AuthenticationType, "key") {
			return nil, fmt.Errorf("invalid authentication type for azure")
		}

		azureSecret := config.Spec.SecretAzure
		if azureSecret == nil {
			return nil, fmt.Errorf("azure secret missing")
		}

		return azure.New(config.Spec.BucketName, *azureSecret)
	}

	return nil, fmt.Errorf("invalid protocol (%v)", config.Spec.Protocols)
}
