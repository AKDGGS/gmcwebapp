package assets

import (
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"path"

	"gmc/cache"
)

func ServeStatic(w http.ResponseWriter, r *http.Request) {
	ServeStaticPath(r.URL.Path[1:], w, r.Header)
}

func ServeStaticPath(name string, w http.ResponseWriter, h http.Header) {
	s, err := Stat(name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(
				w,
				"file not found",
				http.StatusNotFound,
			)
		} else {
			http.Error(
				w,
				fmt.Sprintf("stat error: %s", err),
				http.StatusInternalServerError,
			)
		}
		return
	}

	e := cache.Get(name)
	if e == nil || s.ModTime().After(e.ModTime) {
		b, err := ReadBytes(name)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("read error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}

		e = cache.NewEntryFull(&b, s.ModTime(), 0)
		cache.Put(name, e)
	}

	enc, etag, content := e.Content(h.Get("Accept-Encoding"))
	// Ignore requests for the same content
	if h.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	if enc != "" {
		w.Header().Set("Content-Encoding", enc)
	}
	w.Header().Set("ETag", etag)
	contenttype := mime.TypeByExtension(path.Ext(name))
	if contenttype == "" {
		contenttype = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contenttype)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(*content)))
	w.Write(*content)
}
