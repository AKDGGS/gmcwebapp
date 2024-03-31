package cache

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"hash"
)

type Buffer struct {
	buf  bytes.Buffer
	hash hash.Hash
}

func NewBuffer() Buffer {
	buf := Buffer{hash: md5.New()}
	return buf
}

func (b *Buffer) Write(p []byte) (int, error) {
	b.hash.Write(p)
	return b.buf.Write(p)
}

func (b *Buffer) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *Buffer) Len() int {
	return b.buf.Len()
}

func (b *Buffer) Hash() string {
	return fmt.Sprintf("%x", b.hash.Sum(nil))
}
