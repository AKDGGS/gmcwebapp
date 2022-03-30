package web

import (
	"bytes"
	"fmt"
	"gmc/assets"
	dbf "gmc/db/flag"
	"net/http"
)

func (srv *Server) ServeStash(id int, w http.ResponseWriter, r *http.Request) {
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

	stash, err := srv.DB.GetStash(id, flags)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	// If no details were returned, throw a 404
	if stash == nil {
		http.Error(w, "Inventory stash not found", http.StatusNotFound)
		return
	}

	stash_params := map[string]interface{}{
		"stash": stash,
		"user":  user,
	}

	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/stash.html", &buf, stash_params); err != nil {
		http.Error(
			w, fmt.Sprintf("Parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// params := map[string]interface{}{
	// 	"title":   "Borehole Detail",
	// 	"content": template.HTML(buf.String()),
	// 	"stylesheets": []string{
	// 		"ol/ol.css", "ol/ol-layerswitcher.min.css",
	// 		"css/view.css",
	// 	},
	// 	"scripts": []string{
	// 		"ol/ol.js", "ol/ol-layerswitcher.min.js",
	// 		"js/mustache.js", "js/view.js",
	// 	},
	// 	"redirect": fmt.Sprintf("borehole/%d", id),
	// 	"user":     user,
	// }
	//
	// tbuf := bytes.Buffer{}
	// if err := assets.ExecuteTemplate("tmpl/template.html", &tbuf, params); err != nil {
	// 	http.Error(
	// 		w, fmt.Sprintf("Parse error: %s", err.Error()),
	// 		http.StatusInternalServerError,
	// 	)
	// 	return
	// }

	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}
