package web

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (srv *Server) ServeWellPoints(name string, w http.ResponseWriter, r *http.Request) {
	pts, err := srv.DB.GetWellPoints()
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// If no details were returned, return empty json
	if pts == nil {
		pts = make([]map[string]interface{}, 0)
	}

	js, err := json.Marshal(pts)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("JSON error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	content := &js
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		// 860 bytes is Akamai's recommended minimum for gzip
		// so only bother to gzip files greater than 860 bytes
		if len(js) > 860 {
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

			var gzc []byte
			// Only accept gzip if it's less than the original in size
			if buf.Len() > 0 && buf.Len() < len(js) {
				gzc = buf.Bytes()
				content = &gzc
			}
			w.Header().Set("Content-Encoding", "gzip")
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(*content)))
	w.Write(*content)
}
