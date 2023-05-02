package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"gmc/assets"
)

func (srv *Server) ServeFileDrop(w http.ResponseWriter, r *http.Request) {
	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/filedrop.html", &buf, nil); err != nil {
		http.Error(
			w, fmt.Sprintf("Parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	params := map[string]interface{}{
		"title":   "File Drop Page",
		"content": template.HTML(buf.String()),
		"stylesheets": []string{
			"../css/filedrop.css",
		},
		"scripts": []string{
			"../js/filedrop.js",
		},
		"redirect": "filedrop/",
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
