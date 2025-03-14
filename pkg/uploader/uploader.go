// Copyright 2025 anza-labs contributors
// SPDX-License-Identifier: MIT

package uploader

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/anza-labs/cosi-sample-app/pkg/storage"
)

type Uploader struct {
	client   storage.Storage
	interval time.Duration
}

func New(client storage.Storage, intervalSeconds int) *Uploader {
	return &Uploader{
		client:   client,
		interval: time.Second * time.Duration(intervalSeconds),
	}
}

func (u *Uploader) Run(ctx context.Context, filepath string) error {
	ticker := time.NewTicker(u.interval)
	log.Printf("interval for %v is %v", filepath, u.interval)
	for {
		select {
		case <-ctx.Done():
			log.Printf("context done")
			if err := ctx.Err(); err != nil {
				if !errors.Is(err, context.Canceled) {
					log.Printf("bye (with error)")
					return err
				}
			}
			log.Printf("bye")
			return nil

		case <-ticker.C:
			go func() {
				log.Printf("uploading %v", filepath)
				if err := u.upload(ctx, filepath); err != nil {
					log.Printf("[ERROR] %v", err)
				}
			}()
		}
	}
}

func (u *Uploader) upload(
	ctx context.Context,
	filepath string,
) error {
	stat, err := os.Stat(filepath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	fp, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer fp.Close() //nolint:errcheck // best effort call

	key := fmt.Sprintf("%s-%s", time.Now().Format("2006-Jan-02-15-04-05"), stat.Name())

	if err := u.client.Put(ctx, key, fp, stat.Size()); err != nil {
		return fmt.Errorf("put failed: %w", err)
	}

	return nil
}
