package web

import (
	"bytes"
	"fmt"
	"gmc/assets"
	"html/template"
	"net/http"
)

func (srv *Server) ServeWells(w http.ResponseWriter, r *http.Request) {

	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/wells.html", &buf, nil); err != nil {
		http.Error(
			w, fmt.Sprintf("Parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	params := map[string]interface{}{
		"content": template.HTML(buf.String()),
		"stylesheets": []string{
			"ol/ol.css", "ol/ol-layerswitcher.min.css",
			"css/wellspage.css",
		},
		"scripts": []string{
			"ol/ol.js", "ol/ol-layerswitcher.min.js",
			"js/mustache.js", "js/wells.js",
		},
		"redirect": fmt.Sprint("wells"),
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
