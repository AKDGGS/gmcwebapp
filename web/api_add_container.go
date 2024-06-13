package web

import (
	"fmt"
	"net/http"
)

func (srv *Server) ServeAPIAddContainer(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("authentication error: %s", err),
			http.StatusBadRequest,
		)
		return
	}
	if user == nil {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}
	q := r.URL.Query()
	err = srv.DB.AddContainer(q.Get("barcode"), q.Get("name"), q.Get("remark"))
	if err != nil {
		http.Error(
			w, fmt.Sprintf("add container error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	js := []byte(`{"success":true}`)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
