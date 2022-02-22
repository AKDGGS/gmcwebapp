package web

import (
	"bytes"
	"fmt"
	"gmc/assets"
	"net/http"
)

func (srv *Server) ServeLogin(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Authentication error: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}

	if user != nil {
		http.Redirect(w, r, ".", http.StatusFound)
		return
	}

	buf := bytes.Buffer{}
	if err := assets.ExecuteTemplate("tmpl/login.html", &buf, nil); err != nil {
		http.Error(
			w, fmt.Sprintf("Parse error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(buf.Bytes())
}
