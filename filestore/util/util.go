package util

import (
	"crypto/md5"
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
	content      io.ReadSeekCloser
}

func (f *File) Read(p []byte) (int, error) {
	i, e := f.content.Read(p)
	if f.hasher != nil {
		if n, err := f.hasher.Write(p[:i]); err != nil {
			return n, err
		}
	}
	return i, e
}

func (f *File) Seek(o int64, w int) (int64, error) {
	return f.content.Seek(o, w)
}

func (f *File) Close() error {
	return f.content.Close()
}

func (f *File) SetContent(c io.ReadSeekCloser, hash bool) {
	if hash {
		f.hasher = md5.New()
	}
	f.content = c
}

func (f *File) GetHash() []byte {
	if f.hasher != nil {
		return f.hasher.Sum(nil)
	}
	return nil
}
