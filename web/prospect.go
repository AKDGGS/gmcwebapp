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

func (srv *Server) ServeProspect(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "invalid prospect id", http.StatusBadRequest)
		return
	}

	prospect, err := srv.DB.GetProspect(id, flags)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	// If no details are returned, throw a 404
	if prospect == nil {
		http.Error(w, "Prospect not found", http.StatusNotFound)
		return
	}

	prospectParams := map[string]interface{}{
		"prospect": prospect,
		"user":     user,
	}

	pbuf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/prospect.html", &pbuf, prospectParams); err != nil {
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
			w, fmt.Sprintf("Parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", tbuf.Len()))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(tbuf.Bytes())
}
