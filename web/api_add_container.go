package web

import (
	"fmt"
	"net/http"

	dbe "gmc/db/errors"
)

func (srv *Server) ServeAPIAddContainer(w http.ResponseWriter, r *http.Request) {
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
	err = srv.DB.AddContainer(q.Get("barcode"), q.Get("name"), q.Get("remark"))
	if err != nil {
		if err == dbe.ErrBarcodeCannotBeEmpty {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err == dbe.ErrBarcodeExists {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(
				w, fmt.Sprintf("Error: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
		return
	}
	js := []byte(`{"success":true}`)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
