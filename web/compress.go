package web

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/andybalholm/brotli"
)

type ResponseWrapper struct{ http.ResponseWriter }

func (rw ResponseWrapper) Close() error {
	return nil
}

func compressWriter(accept string, w http.ResponseWriter) (io.WriteCloser, error) {
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

	return &ResponseWrapper{w}, nil
}
