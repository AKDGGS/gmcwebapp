package filestore

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"os"
)

type S3 struct {
	client *minio.Client
	bucket string
}

func newS3(cfg map[string]interface{}) (*S3, error) {
	endpoint, ok := cfg["endpoint"].(string)
	if !ok {
		return nil, fmt.Errorf("s3 endpoint must exist and be a string")
	}

	accesskeyid, ok := cfg["accesskeyid"].(string)
	if !ok {
		return nil, fmt.Errorf("s3 accesskeyid must exist and be a string")
	}

	secretaccesskey, ok := cfg["secretaccesskey"].(string)
	if !ok {
		return nil, fmt.Errorf("s3 secretaccesskey must exist and be a string")
	}

	bucket, ok := cfg["bucket"].(string)
	if !ok {
		return nil, fmt.Errorf("s3 bucket must exist and be a string")
	}

	secure, ok := cfg["secure"].(bool)
	if !ok {
		secure = false
	}

	var err error
	s3 := &S3{bucket: bucket}
	// Setup S3 Connection
	s3.client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accesskeyid, secretaccesskey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, err
	}

	return s3, nil
}

func (s3 *S3) GetFile(name string) (*File, error) {
	obj, err := s3.client.GetObject(
		context.Background(),
		s3.bucket, name,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, err
	}

	// Verify file exists so a proper (and checkable)
	// error can be returned
	stat, err := obj.Stat()
	if err != nil {
		switch v := err.(type) {
		case minio.ErrorResponse:
			if v.Code == "NoSuchKey" {
				return nil, &os.PathError{Op: "stat", Path: name, Err: err}
			}
		}
		return nil, err
	}

	return &File{
		Name:         stat.Key,
		ETag:         stat.ETag,
		LastModified: stat.LastModified,
		Size:         stat.Size,
		ContentType:  stat.ContentType,
		Content:      obj,
	}, nil
}

func (s3 *S3) Shutdown() {
	// Do nothing
}
