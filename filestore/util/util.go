package util

import (
	"io"
	"time"
)

type File struct {
	Name         string
	ETag         string
	Size         int64
	LastModified time.Time
	ContentType  string
	Content      io.ReadSeekCloser
}
