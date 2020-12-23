package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"strings"
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

func (u *gcpUploader) Do(ctx context.Context, content, fileName string) error {

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := u.client.Bucket(u.bucket).Object(fileName).NewWriter(ctx)

	if _, err := io.Copy(wc, strings.NewReader(content)); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	fmt.Printf("\nBlob %v uploaded.\n", fileName)
	return nil
}
