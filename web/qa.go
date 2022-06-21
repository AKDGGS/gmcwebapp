package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"gmc/assets"
)

func (srv *Server) ServeQA(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Authentication error: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}
	if user == nil {
		http.Redirect(w, r, "login?redirect=qa", http.StatusFound)
		return
	}

	qar, err := srv.DB.ListQAReports()
	if err != nil {
		http.Error(
			w, fmt.Sprintf("QA report error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/qa.html", &buf, qar); err != nil {
		http.Error(
			w, fmt.Sprintf("Parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	params := map[string]interface{}{
		"title":   "Quality Assurance",
		"scripts": []string{"js/qa.js"},
		"content": template.HTML(buf.String()),
		"user":    user,
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
