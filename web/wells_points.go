package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gmc/cache"
)

func (srv *Server) ServeWellsPointsJSON(w http.ResponseWriter, r *http.Request) {
	e := cache.Get("wells/points.json")
	if e == nil {
		pts, err := srv.DB.GetWellPoints()
		if err != nil {
			http.Error(
				w, fmt.Sprintf("error: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		// If no details are returned, return an empty json
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

		e = cache.NewEntry(&js)
		cache.Put("wells/points.json", e)
	}

	enc, etag, content := e.Content(r.Header.Get("Accept-Encoding"))
	// Ignore requests for the same content
	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	if enc != "" {
		w.Header().Set("Content-Encoding", enc)
	}
	w.Header().Set("ETag", etag)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(*content)))
	w.Write(*content)
}
