package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gmc/assets"
)

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, srv.Config.BasePath)
	switch path {
	case "favicon.ico":
		assets.ServeStatic("img/favicon.ico", w, r)
		return

	case "ol/ol.css":
		assets.ServeStatic("ol/ol-v7.1.0.css", w, r)
		return

	case "ol/ol.js":
		assets.ServeStatic("ol/ol-v7.1.0.js", w, r)
		return

	case "js/mustache.js":
		assets.ServeStatic("js/mustache-v4.2.0.js", w, r)
		return

	case "css/template.css", "css/map-defaults.css", "css/view.css",
		"css/wells.css", "css/qa.css", "js/map-defaults.js",
		"js/view.js", "js/stash.js", "js/wells.js", "js/qa.js",
		"img/loader.gif", "ol/ol-layerswitcher.min.css",
		"ol/ol-layerswitcher.min.js", "css/filedrop.css", "js/filedrop.js":
		assets.ServeStatic(path, w, r)
		return

	case "qa/":
		srv.ServeQA(w, r)
		return

	case "login":
		err := srv.Auths.CheckForm(w, r)
		if err != nil {
			http.Error(
				w, fmt.Sprintf("error: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
		return

	case "logout":
		err := srv.Auths.Logout(w, r)
		if err != nil {
			http.Error(
				w, fmt.Sprintf("error: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
		return

	case "stash.json":
		q := r.URL.Query()
		id, err := strconv.Atoi(q.Get("id"))
		if err != nil {
			http.Error(w, "Invalid Inventory ID", http.StatusBadRequest)
			return
		}
		srv.ServeStash(id, w, r)
		return

	case "qa/qa_count.json":
		q := r.URL.Query()
		id, err := strconv.Atoi(q.Get("id"))
		if err != nil {
			http.Error(w, "Invalid Report ID", http.StatusBadRequest)
			return
		}
		srv.ServeQACount(id, w, r)
		return

	case "qa/qa_run.json":
		q := r.URL.Query()
		id, err := strconv.Atoi(q.Get("id"))
		if err != nil {
			http.Error(w, "Invalid Report ID", http.StatusBadRequest)
			return
		}
		srv.ServeQARun(id, w, r)
		return

	case "wells/":
		srv.ServeWells(w, r)
		return

	case "wells/points.json":
		srv.ServeWellsPointsJSON(w, r)
		return

	case "wells/detail.json":
		q := r.URL.Query()
		id, err := strconv.Atoi(q.Get("id"))
		if err != nil {
			http.Error(w, "Invalid Well ID", http.StatusBadRequest)
			return
		}
		srv.ServeWellsDetailJSON(id, w, r)
		return

	case "upload":
		file_id := srv.ServeUpload(w, r)
		w.Header().Set("file_id", strconv.Itoa(int(file_id)))
		return
	}

	sidx := strings.Index(path, "/")
	if sidx == -1 {
		sidx = len(path)
	}
	action := path[:sidx]

	switch action {
	case "file":
		fid := strings.TrimPrefix(strings.TrimPrefix(path, "file"), "/")
		id, err := strconv.Atoi(fid)
		if err != nil {
			http.Error(w, "Invalid File ID", http.StatusBadRequest)
			return
		}
		srv.ServeFile(id, w, r)

	case "prospect":
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "prospect"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid Prospect ID", http.StatusBadRequest)
			return
		}
		srv.ServeProspect(id, w, r)

	case "borehole":
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "borehole"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid Borehole ID", http.StatusBadRequest)
			return
		}
		srv.ServeBorehole(id, w, r)

	case "outcrop":
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "outcrop"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid Outcrop ID", http.StatusBadRequest)
			return
		}
		srv.ServeOutcrop(id, w, r)

	case "well":
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "well"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid Well ID", http.StatusBadRequest)
			return
		}
		srv.ServeWell(id, w, r)

	case "shotline":
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "shotline"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid Shotline ID", http.StatusBadRequest)
			return
		}
		srv.ServeShotline(id, w, r)

	case "inventory":
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "inventory"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid Inventory ID", http.StatusBadRequest)
			return
		}
		srv.ServeInventory(id, w, r)

	default:
		http.Error(w, "File not found", http.StatusNotFound)
	}
}
