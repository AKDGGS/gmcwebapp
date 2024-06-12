package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"gmc/assets"
	dbf "gmc/db/flag"
)

func (srv *Server) ServeBorehole(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("authentication error: %s", err.Error()),
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
		http.Error(w, "invalid borehole id", http.StatusBadRequest)
		return
	}

	borehole, err := srv.DB.GetBorehole(id, flags)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	// If no details are returned, throw a 404
	if borehole == nil {
		http.Error(w, "borehole not found", http.StatusNotFound)
		return
	}

	boreholeParams := map[string]interface{}{
		"borehole": borehole,
		"user":     user,
	}

	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/borehole.html", &buf, boreholeParams); err != nil {
		http.Error(
			w, fmt.Sprintf("parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	params := map[string]interface{}{
		"title":   "Borehole Detail",
		"content": template.HTML(buf.String()),
		"stylesheets": []string{
			"../ol/ol.css", "../ol/ol-layerswitcher.min.css",
			"../css/map-defaults.css", "../css/view.css",
			"../css/filedrop.css",
		},
		"scripts": []string{
			"../ol/ol.js", "../ol/ol-layerswitcher.min.js",
			"../js/mustache.js", "../js/map-defaults.js", "../js/filedrop.js",
			"../js/view.js",
		},
		"redirect": fmt.Sprintf("borehole/%d", id),
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

	w.Header().Set("Content-Length", fmt.Sprintf("%d", tbuf.Len()))
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
