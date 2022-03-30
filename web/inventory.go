package web

import (
	"bytes"
	"fmt"
	"gmc/assets"
	dbf "gmc/db/flag"
	"html/template"
	"net/http"
)

func (srv *Server) Callme(fn, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You called me!")
}

func (srv *Server) ServeInventory(id int, w http.ResponseWriter, r *http.Request) {
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

	inventory, err := srv.DB.GetInventory(id, flags)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	// If no details were returned, throw a 404
	if inventory == nil {
		http.Error(w, "Inventory not found", http.StatusNotFound)
		return
	}

	// If can_publish is false, throw a 403
	if user == nil && inventory["can_publish"] == false {
		http.Error(w, "Access denied.", http.StatusForbidden)
		return
	}

	inventory["_user"] = user

	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/inventory.html", &buf, inventory); err != nil {
		http.Error(
			w, fmt.Sprintf("Parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	params := map[string]interface{}{
		"title":   "Inventory Detail",
		"content": template.HTML(buf.String()),
		"stylesheets": []string{
			"ol/ol.css", "ol/ol-layerswitcher.min.css",
			"css/view.css",
		},
		"scripts": []string{
			"ol/ol.js", "ol/ol-layerswitcher.min.js",
			"js/mustache.js", "js/view.js",
		},
		"redirect": fmt.Sprintf("inventory/%d", id),
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
