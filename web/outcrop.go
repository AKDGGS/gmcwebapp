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

func (srv *Server) ServeOutcrop(w http.ResponseWriter, r *http.Request) {
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
			"invalid outcrop id",
			http.StatusBadRequest,
		)
		return
	}
	outcrop, err := srv.DB.GetOutcrop(id, flags)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("get outcrop error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	// If no details are returned, throw a 404
	if outcrop == nil {
		http.Error(
			w,
			"outcrop not found",
			http.StatusNotFound)
		return
	}
	outcropParams := map[string]interface{}{
		"outcrop": outcrop,
		"user":    user,
	}
	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/outcrop.html", &buf, outcropParams); err != nil {
		http.Error(
			w,
			fmt.Sprintf("parse error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	params := map[string]interface{}{
		"title":   "Outcrop Detail",
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
		"redirect": fmt.Sprintf("outcrop/%d", id),
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
