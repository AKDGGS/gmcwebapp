package web

import (
	"fmt"
	"net/http"

	dbe "gmc/db/errors"
)

func (srv *Server) ServeAPIMoveInventoryAndContainersContents(w http.ResponseWriter, r *http.Request) {
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
	src := q.Get("src")
	dest := q.Get("dest")
	err = srv.DB.MoveInventoryAndContainersContents(src, dest)
	if err != nil {
		switch err {
		case dbe.ErrDestinationBarcodeEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case dbe.ErrDestinationBarcodeEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case dbe.ErrDestinationNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case dbe.ErrSourceNotValid:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		case dbe.ErrDestinationMultipleContainers:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		case dbe.ErrNothingMoved:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
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
