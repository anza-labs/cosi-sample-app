// Copyright 2025 anza-labs contributors
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"

	"github.com/anza-labs/cosi-sample-app/pkg/storage"
	"github.com/anza-labs/cosi-sample-app/pkg/uploader"
)

func main() {
	opts := runOptions{}

	flag.StringVar(&opts.cosiConfig, "bucket-info", "/cosi/BucketInfo.json", "")
	flag.StringArrayVar(&opts.files, "file", []string{}, "")
	flag.IntVar(&opts.interval, "upload-interval", 10, "")
	flag.Parse()

	if err := run(context.Background(), opts); err != nil {
		log.Fatal(err)
	}
}

type runOptions struct {
	cosiConfig string
	files      []string
	interval   int
}

func run(ctx context.Context, opts runOptions) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Printf("reading config from %s", opts.cosiConfig)
	f, err := os.Open(opts.cosiConfig)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close() //nolint:errcheck // best effort call

	var cfg storage.Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return fmt.Errorf("failed to decode file: %w", err)
	}

	client, err := storage.New(cfg, true)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	u := uploader.New(client, opts.interval)

	eg, ctx := errgroup.WithContext(ctx)

	for _, file := range opts.files {
		eg.Go(func() error {
			log.Printf("starting for %s", file)
			if err := u.Run(ctx, file); err != nil {
				return fmt.Errorf("%s: %w", file, err)
			}
			return nil
		})
	}

	return eg.Wait()
}
