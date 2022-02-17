package web

import (
	"gmc/assets"
	"net/http"
	"strconv"
	"strings"
)

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, srv.Config.BasePath)
	switch path {
	case "favicon.ico":
		assets.ServeStatic("img/favicon.ico", w, r)
		return

	case "ol/ol.css":
		assets.ServeStatic("ol/ol-v6.9.0.css", w, r)
		return

	case "ol/ol.js":
		assets.ServeStatic("ol/ol-v6.9.0.js", w, r)
		return

	case "js/mustache.js":
		assets.ServeStatic("js/mustache-v4.2.0.js", w, r)
		return

	case "css/template.css", "css/view.css", "js/view.js",
		"ol/ol-layerswitcher.min.css", "ol/ol-layerswitcher.min.js":
		assets.ServeStatic(path, w, r)
		return
	}

	sidx := strings.Index(path, "/")
	if sidx == -1 {
		sidx = len(path)
	}
	action := path[:sidx]

	switch action {
	case "file":
		fn := strings.TrimPrefix(strings.TrimPrefix(path, "file"), "/")
		srv.ServeFile(fn, w, r)

	case "prospect":
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "prospect"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid Prospect ID", http.StatusBadRequest)
			return
		}
		srv.ServeProspect(id, w)

	case "borehole":
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "borehole"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid Borehole ID", http.StatusBadRequest)
			return
		}
		srv.ServeBorehole(id, w)

	case "outcrop":
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "outcrop"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid Outcrop ID", http.StatusBadRequest)
			return
		}
		srv.ServeOutcrop(id, w)

	case "well":
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "well"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid Well ID", http.StatusBadRequest)
			return
		}
		srv.ServeWell(id, w)

	case "shotline":
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "shotline"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid Shotline ID", http.StatusBadRequest)
			return
		}
		srv.ServeShotline(id, w)

	default:
		http.Error(w, "File not found", http.StatusNotFound)
	}
}
