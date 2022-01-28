package web

import (
	"bytes"
	"fmt"
	"gmc/assets"
	"gmc/db"
	"html/template"
	"net/http"
)

func (srv *Server) ServeBorehole(id int, w http.ResponseWriter) {
	borehole, err := srv.DB.GetBorehole(id, db.ALL_NOPRIVATE)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	// If no details were returned, throw a 404
	if borehole == nil {
		http.Error(w, "Borehole not found", http.StatusNotFound)
		return
	}

	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/borehole.html", &buf, borehole); err != nil {
		http.Error(
			w, fmt.Sprintf("Parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	params := map[string]interface{}{
		"title":   "Borehole Detail",
		"content": template.HTML(buf.String()),
		"stylesheets": []string{
			"ol/ol.css", "ol/ol-layerswitcher.min.css",
			"css/view.css",
		},
		"scripts": []string{
			"ol/ol.js", "ol/ol-layerswitcher.min.js",
			"js/mustache.js", "js/view.js",
		},
	}

	tbuf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/template.html", &tbuf, params); err != nil {
		http.Error(
			w, fmt.Sprintf("Parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", tbuf.Len()))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(tbuf.Bytes())
}