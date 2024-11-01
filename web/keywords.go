package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gmc/cache"
)

func (srv *Server) ServeKeywords(w http.ResponseWriter, r *http.Request) {
	e := cache.Get("keywords.json")
	if e == nil {
		keywords, err := srv.DB.ListKeywords()
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("list keywords error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
		js, err := json.Marshal(keywords)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("json marshal error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
		e = cache.NewEntry(&js)
		cache.Put("keywords.json", e)
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
