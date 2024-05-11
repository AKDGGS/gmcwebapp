package web

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/andybalholm/brotli"
)

type CloseFlusher interface {
	Flush() error
	Close() error
}

func compressWriter(accept string, w http.ResponseWriter) (io.Writer, error) {
	encs := strings.Split(accept, ", ")
	if slices.Contains(encs, "br") {
		br := brotli.NewWriterLevel(w, brotli.DefaultCompression)
		w.Header().Set("Content-Encoding", "br")
		return br, nil
	}
	if slices.Contains(encs, "gzip") {
		gz, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
		if err != nil {
			return nil, err
		}
		w.Header().Set("Content-Encoding", "gzip")
		return gz, nil
	}
	return w, nil
}
