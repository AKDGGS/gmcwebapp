package s3

import (
	"context"
	"fmt"
	"os"
	"path"

	fsutil "gmc/filestore/util"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3 struct {
	client *minio.Client
	bucket string
}

func New(cfg map[string]interface{}) (*S3, error) {
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

func (s3 *S3) GetFile(name string) (*fsutil.File, error) {
	objectCh := s3.client.ListObjects(context.Background(), s3.bucket,
		minio.ListObjectsOptions{Prefix: fmt.Sprintf("%s/", name)})

	var objectKey string
	objectCount := 0
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objectKey = object.Key
		objectCount++
		if objectCount > 1 {
			return nil, fmt.Errorf("more than one object found")
		}
	}

	obj, err := s3.client.GetObject(
		context.Background(),
		s3.bucket, objectKey,
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
	return &fsutil.File{
		Name:         path.Base(stat.Key),
		ETag:         stat.ETag,
		LastModified: stat.LastModified,
		Size:         stat.Size,
		ContentType:  stat.ContentType,
		Content:      obj,
	}, nil
}

func (s3 *S3) PutFile(file *fsutil.File) error {
	// Check if bucket exists
	exists, err := s3.client.BucketExists(context.Background(), s3.bucket)
	if err != nil {
		return err
	}

	// Create bucket if it doesn't exist
	if !exists {
		err := s3.client.MakeBucket(context.Background(), s3.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	_, err = s3.client.PutObject(context.Background(), s3.bucket, file.Name,
		file.Content, file.Size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (s3 *S3) Shutdown() {
	// Do nothing
}
