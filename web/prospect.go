package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"gmc/assets"
	"gmc/cache"
	dbf "gmc/db/flag"
)

func (srv *Server) ServeProspects(w http.ResponseWriter, r *http.Request) {
	e := cache.Get("prospects.json")
	if e == nil {
		prospects, err := srv.DB.ListProspects()
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("list prospects error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
		js, err := json.Marshal(prospects)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("json marshal error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
		e = cache.NewEntry(&js)
		cache.Put("prospects.json", e)
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

func (srv *Server) ServeProspect(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("authentication error: %s", err),
			http.StatusBadRequest,
		)
		return
	}
	flags := dbf.ALL
	if user == nil {
		flags = dbf.ALL_NOPRIVATE
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			"invalid prospect id",
			http.StatusBadRequest,
		)
		return
	}
	prospect, err := srv.DB.GetProspect(id, flags)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("get prospect error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	// If no details are returned, throw a 404
	if prospect == nil {
		http.Error(
			w,
			"prospect not found",
			http.StatusNotFound,
		)
		return
	}
	prospectParams := map[string]interface{}{
		"prospect": prospect,
		"user":     user,
	}
	pbuf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/prospect.html", &pbuf, prospectParams); err != nil {
		http.Error(
			w,
			fmt.Sprintf("parse error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	params := map[string]interface{}{
		"title":   "Prospect Detail",
		"content": template.HTML(pbuf.String()),
		"stylesheets": []string{
			"../ol/ol.css", "../ol/ol-layerswitcher.min.css",
			"../css/map-defaults.css", "../css/view.css",
			"../css/filedrop.css",
		},
		"scripts": []string{
			"../ol/ol.js", "../ol/ol-layerswitcher.min.js",
			"../js/mustache.js", "../js/map-defaults.js",
			"../js/filedrop.js", "../js/view.js",
		},
		"redirect": fmt.Sprintf("prospect/%d", id),
		"user":     user,
	}
	tbuf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/template.html", &tbuf, params); err != nil {
		http.Error(
			w,
			fmt.Sprintf("parse error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", tbuf.Len()))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	out, err := compressWriter(r.Header.Get("Accept-Encoding"), w)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("compression error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	defer out.Close()
	out.Write(tbuf.Bytes())
}
