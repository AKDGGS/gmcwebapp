package web

import (
	"bytes"
	"fmt"
	"gmc/assets"
	dbf "gmc/db/flag"
	"html/template"
	"net/http"
)

func (srv *Server) ServeProspect(id int, w http.ResponseWriter, r *http.Request) {
	user, err := srv.AuthOptional(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Authentication error: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	flags := dbf.ALL
	if user == nil {
		flags = dbf.ALL_NOPRIVATE
	}

	prospect, err := srv.DB.GetProspect(id, flags)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	// If no details were returned, throw a 404
	if prospect == nil {
		http.Error(w, "Prospect not found", http.StatusNotFound)
		return
	}

	pbuf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/prospect.html", &pbuf, prospect); err != nil {
		http.Error(
			w, fmt.Sprintf("Parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	params := map[string]interface{}{
		"title":   "Prospect Detail",
		"content": template.HTML(pbuf.String()),
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
