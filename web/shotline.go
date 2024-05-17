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

func (srv *Server) ServeShotline(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
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

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid shotline id", http.StatusBadRequest)
		return
	}

	shotline, err := srv.DB.GetShotline(id, flags)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	// If no details are returned, throw a 404
	if shotline == nil {
		http.Error(w, "Shotline not found", http.StatusNotFound)
		return
	}

	shotlineParams := map[string]interface{}{
		"shotline": shotline,
		"user":     user,
	}

	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/shotline.html", &buf, shotlineParams); err != nil {
		http.Error(
			w, fmt.Sprintf("Parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	params := map[string]interface{}{
		"title":   "Shotline Detail",
		"content": template.HTML(buf.String()),
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
		"redirect": fmt.Sprintf("shotline/%d", id),
		"user":     user,
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
