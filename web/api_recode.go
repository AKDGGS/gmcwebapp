package web

import (
	"fmt"
	"net/http"
)

func (srv *Server) ServeAPIRecodeInventoryAndContainer(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Authentication error: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}
	if user == nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	q := r.URL.Query()
	old_barcode := q.Get("old")
	new_barcode := q.Get("new")
	err = srv.DB.RecodeInventoryAndContainer(old_barcode, new_barcode)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
}
