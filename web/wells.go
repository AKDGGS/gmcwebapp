package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"gmc/assets"
)

func (srv *Server) ServeWells(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("authentication error: %s", err),
			http.StatusBadRequest,
		)
		return
	}
	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/wells.html", &buf, nil); err != nil {
		http.Error(
			w,
			fmt.Sprintf("parse error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	params := map[string]interface{}{
		"title":   "Wells Page",
		"content": template.HTML(buf.String()),
		"stylesheets": []string{
			"../ol/ol.css", "../ol/ol-layerswitcher.min.css",
			"../css/map-defaults.css", "../css/wells.css",
		},
		"scripts": []string{
			"../ol/ol.js", "../ol/ol-layerswitcher.min.js",
			"../js/mustache.js", "../js/map-defaults.js",
			"../js/wells.js",
		},
		"redirect": "wells/",
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
	w.Write(tbuf.Bytes())
}
