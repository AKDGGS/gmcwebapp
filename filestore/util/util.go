package util

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"time"
)

type File struct {
	Name         string
	ETag         string
	Size         int64
	LastModified time.Time
	ContentType  string
	hasher       hash.Hash
	readsize     int64
	content      io.ReadCloser
}

type FileStat struct {
	Name         string
	Size         int64
	LastModified time.Time
	Hash         string
}

// The way this calculates a hash is problematic if downstream
// the upload process uses seeking.
// Currently (2/9/26) all of our uploads use io.ReadCloser (no seek)
func (f *File) SetContent(c io.ReadCloser, hash bool) {
	if hash {
		f.hasher = md5.New()
	}
	f.content = c
}

func (f *File) Read(p []byte) (int, error) {
	i, e := f.content.Read(p)
	f.readsize += int64(i)
	if f.hasher != nil {
		if n, err := f.hasher.Write(p[:i]); err != nil {
			return n, err
		}
	}
	return i, e
}

func (f *File) Seek(o int64, w int) (int64, error) {
	if rsc, ok := f.content.(io.ReadSeekCloser); ok {
		return rsc.Seek(o, w)
	}
	return 0, fmt.Errorf("seek unsupported")
}

func (f *File) Close() error {
	return f.content.Close()
}

func (f *File) GetHash() []byte {
	if f.hasher != nil {
		return f.hasher.Sum(nil)
	}
	return nil
}

func (f *File) GetReadSize() int64 {
	return f.readsize
}
