package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type gcpUploader struct {
	client *storage.Client
	bucket string
}

func NewGcpUploader(c *storage.Client, bucket string) *gcpUploader {
	return &gcpUploader{
		client: c,
		bucket: bucket,
	}
}

func (u *gcpUploader) Do(ctx context.Context, sourceFileName string) error {

	f, err := os.Open(sourceFileName)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := u.client.Bucket(u.bucket).Object(sourceFileName).NewWriter(ctx)

	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	fmt.Printf("\nBlob %v uploaded.\n", sourceFileName)
	return nil
}
