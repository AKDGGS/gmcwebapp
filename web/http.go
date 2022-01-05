package web

import (
	"bytes"
	"errors"
	"fmt"
	"gmc/assets"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
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

	case "css/template.css", "css/prospect.css", "js/prospect.js":
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
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "file"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// Fetch the file details from the database
		aid, fname, ftime, err := srv.DB.GetFile(id, false)
		if err != nil {
			http.Error(
				w, fmt.Sprintf("Query error: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		// Fetch the file from S3
		file, err := srv.FileStore.GetFile(fmt.Sprintf("%d/%s", aid, fname))
		if err != nil {
			if _, ok := err.(*os.PathError); ok {
				http.Error(w, "File not found (FileStore)", http.StatusNotFound)
			} else {
				http.Error(
					w, fmt.Sprintf("file fetch error: %s", err.Error()),
					http.StatusInternalServerError,
				)
			}
			return
		}
		defer file.Close()

		// Suggest filename to the browser
		w.Header().Set(
			"Content-Disposition",
			fmt.Sprintf("inline; filename=\"%s\"", fname),
		)
		http.ServeContent(w, r, fname, ftime, file)

	case "prospect":
		sid := strings.TrimPrefix(strings.TrimPrefix(path, "prospect"), "/")
		id, err := strconv.Atoi(sid)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		prospect, err := srv.DB.GetProspect(id)
		if err != nil {
			http.Error(
				w, fmt.Sprintf("Query error: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
		// If no details were returned, throw a 404
		if prospect == nil {
			http.Error(w, "Record not found", http.StatusNotFound)
			return
		}

		buf := bytes.Buffer{}
		err = assets.ExecuteTemplate("tmpl/prospect.html", &buf, prospect)
		if err != nil {
			http.Error(
				w, fmt.Sprintf("Parse error: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

		params := map[string]interface{}{
			"title":   "Prospect Detail",
			"content": template.HTML(buf.String()),
			"stylesheets": []string{
				"css/prospect.css", "ol/ol.css",
			},
			"scripts": []string{
				"ol/ol.js", "js/mustache.js",
				"js/prospect.js",
			},
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = assets.ExecuteTemplate("tmpl/template.html", w, params)
		// Ignore broken pipe errors
		if err != nil && !errors.Is(err, syscall.EPIPE) {
			http.Error(
				w, fmt.Sprintf("Parse error: %s", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}

	default:
		http.Error(w, "File not found", http.StatusNotFound)
	}
}
