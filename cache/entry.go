package cache

import (
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"sync"
	"time"

	"github.com/andybalholm/brotli"
)

func NewEntry(content *[]byte) *Entry {
	e := Entry{
		ModTime: time.Now(),
		Expires: time.Now().Add(time.Duration(60) * time.Minute),
	}
	e.SetContent(content)
	if e.pl_content == nil {
		return nil
	}
	return &e
}

func NewEntryFull(content *[]byte, modtime time.Time, expires int) *Entry {
	e := Entry{ModTime: modtime}
	if expires > 0 {
		e.Expires = time.Now().Add(time.Duration(expires) * time.Minute)
	}
	e.SetContent(content)
	if e.pl_content == nil {
		return nil
	}
	return &e
}

type Entry struct {
	ModTime time.Time
	Expires time.Time

	pl_content *[]byte
	pl_etag    string
	br_content *[]byte
	br_etag    string
	gz_content *[]byte
	gz_etag    string
}

func (e *Entry) Content(accept string) (string, string, *[]byte) {
	enc := parseEncoding(accept)
	if ((enc & BROTLI) != 0) && e.br_content != nil {
		return "br", e.br_etag, e.br_content
	}
	if ((enc & GZIP) != 0) && e.gz_content != nil {
		return "gzip", e.gz_etag, e.gz_content
	}
	return "", e.pl_etag, e.pl_content
}

func (e *Entry) SetContent(content *[]byte) {
	// Refuse to setup nil content
	if content == nil {
		return
	}
	e.pl_content = content

	var wg sync.WaitGroup
	// Background md5 calculation
	wg.Add(1)
	go func() {
		defer wg.Done()
		e.pl_etag = fmt.Sprintf("%x", md5.Sum(*content))
	}()

	// Background gzip
	wg.Add(1)
	go func() {
		defer wg.Done()

		buf := NewBuffer()
		gz, err := gzip.NewWriterLevel(&buf, gzip.DefaultCompression)
		if err != nil {
			return
		}
		defer gz.Close()

		if _, err := gz.Write(*content); err != nil {
			return
		}

		if err := gz.Flush(); err != nil {
			return
		}

		// Only accept gzip if it's less than the original in size
		if buf.Len() > 0 && buf.Len() < len(*content) {
			gzc := buf.Bytes()
			e.gz_content = &gzc
			e.gz_etag = buf.Hash()
		}
	}()

	// Background brotli
	wg.Add(1)
	go func() {
		defer wg.Done()

		buf := NewBuffer()
		br := brotli.NewWriterLevel(&buf, brotli.DefaultCompression)
		defer br.Close()

		if _, err := br.Write(*content); err != nil {
			return
		}

		if err := br.Flush(); err != nil {
			return
		}

		// Only accept brotli if it's less than the original in size
		if buf.Len() > 0 && buf.Len() < len(*content) {
			brc := buf.Bytes()
			e.br_content = &brc
			e.br_etag = buf.Hash()
		}
	}()

	wg.Wait()
}

func (e *Entry) SetExpires(sec int64) {
	if sec == 0 {
		e.Expires = time.Time{}
	} else {
		e.Expires = time.Now().Add(time.Duration(sec) * time.Second)
	}
}
