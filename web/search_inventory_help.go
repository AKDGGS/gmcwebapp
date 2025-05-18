package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"gmc/assets"
)

func (srv *Server) ServeSearchInventoryHelp(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("authentication error: %s", err),
			http.StatusBadRequest,
		)
		return
	}

	stmpl := fmt.Sprintf(
		"tmpl/%s/inventory_search_help.html",
		srv.Search.Name(),
	)
	sparams := map[string]interface{}{
		"user": user,
	}
	sbuf := bytes.Buffer{}

	if err := assets.ExecuteTemplate(stmpl, &sbuf, sparams); err != nil {
		http.Error(
			w,
			fmt.Sprintf("parse error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	params := map[string]interface{}{
		"title":       "Inventory Search Help",
		"content":     template.HTML(sbuf.String()),
		"stylesheets": []string{"../css/help.css"},
		"scripts":     []string{},
		"redirect":    "inventory/search-help",
		"user":        user,
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
