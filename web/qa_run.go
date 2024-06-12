package web

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (srv *Server) ServeQARun(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("authentication error: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}
	if user == nil {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}
	q := r.URL.Query()
	id, err := strconv.Atoi(q.Get("id"))
	if err != nil {
		http.Error(w, "invalid Report ID", http.StatusBadRequest)
		return
	}
	re, err := srv.DB.RunQAReport(id)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	js, err := json.Marshal(re)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("JSON error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	content := &js
	// 860 bytes is Akamai's recommended minimum for gzip
	// so only bother to gzip files greater than 860 bytes
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && len(js) > 860 {
		var buf bytes.Buffer
		gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
		if err != nil {
			http.Error(
				w, fmt.Sprintf("gzip error: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
		defer gz.Close()

		if _, err := gz.Write(js); err != nil {
			http.Error(
				w, fmt.Sprintf("gz write error: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		if err := gz.Flush(); err != nil {
			http.Error(
				w, fmt.Sprintf("gz write error: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		// Only accept gzip if it's less than the original in size
		if buf.Len() > 0 && buf.Len() < len(js) {
			gzc := buf.Bytes()
			content = &gzc
			w.Header().Set("Content-Encoding", "gzip")
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(*content)))
	w.Write(*content)
}
