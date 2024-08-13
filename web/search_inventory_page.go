package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"gmc/assets"
)

func (srv *Server) ServeSearchInventoryPage(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("authentication error: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	sparams := map[string]interface{}{
		"user": user,
	}

	sbuf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/search.html", &sbuf, sparams); err != nil {
		http.Error(
			w, fmt.Sprintf("parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	if err := assets.ExecuteTemplate("tmpl/inventory_search.html", &sbuf, sparams); err != nil {
		http.Error(
			w, fmt.Sprintf("parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	params := map[string]interface{}{
		"title":   "Inventory Search",
		"content": template.HTML(sbuf.String()),
		"stylesheets": []string{
			"../ol/ol.css",
			"../ol/ol-layerswitcher.min.css",
			"../ol/ol-drawbox-control.css",
			"../ol/ol-search-control.css",
			"../css/map-defaults.css",
			"../css/search.css",
		},
		"scripts": []string{
			"../ol/ol.js",
			"../ol/ol-layerswitcher.min.js",
			"../ol/ol-drawbox-control.js",
			"../ol/ol-search-control.js",
			"../js/mustache.js",
			"../js/map-defaults.js",
			"../js/search.js",
		},
		"redirect": "inventory/search",
		"user":     user,
	}

	tbuf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/template.html", &tbuf, params); err != nil {
		http.Error(
			w, fmt.Sprintf("parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	out, err := compressWriter(r.Header.Get("Accept-Encoding"), w)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("compression error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	defer out.Close()
	out.Write(tbuf.Bytes())
}
