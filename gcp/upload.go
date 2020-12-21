package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type uploader struct {
	client *storage.Client
}

func NewUploader(c *storage.Client) *uploader {
	return &uploader{
		client: c,
	}
}

func (u *uploader) Do(ctx context.Context, sourceFileName, bucket string) error {

	// Open local file.
	f, err := os.Open(sourceFileName)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := u.client.Bucket(bucket).Object(sourceFileName).NewWriter(ctx)

	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	fmt.Printf("\nBlob %v uploaded.\n", sourceFileName)
	return nil
}
