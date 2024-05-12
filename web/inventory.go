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

func (srv *Server) ServeInventory(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "invalid inventory id", http.StatusBadRequest)
		return
	}

	inventory, err := srv.DB.GetInventory(id, flags)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// If no details are returned, throw a 404
	if inventory == nil {
		http.Error(w, "Inventory not found", http.StatusNotFound)
		return
	}

	// If can_publish is false, throw a 403
	if user == nil && inventory.CanPublish == false {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	inventory_params := map[string]interface{}{
		"inventory": inventory,
		"user":      user,
	}

	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/inventory.html", &buf, inventory_params); err != nil {
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
			"../ol/ol.css", "../ol/ol-layerswitcher.min.css",
			"../css/map-defaults.css", "../css/view.css",
			"../css/filedrop.css",
		},
		"scripts": []string{
			"../ol/ol.js", "../ol/ol-layerswitcher.min.js",
			"../js/mustache.js", "../js/map-defaults.js",
			"../js/filedrop.js", "../js/view.js",
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	out, err := compressWriter(r.Header.Get("Accept-Encoding"), w)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Compression error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	defer out.Close()
	out.Write(tbuf.Bytes())
}
