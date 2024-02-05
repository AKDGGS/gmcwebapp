package web

import (
	"fmt"
	"net/http"

	dbe "gmc/db/errors"
)

func (srv *Server) ServeAPIMoveInventoryAndContainers(w http.ResponseWriter, r *http.Request) {
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
	dest := q.Get("d")
	barcodes_to_move := q["c"]
	err = srv.DB.MoveInventoryAndContainers(dest, barcodes_to_move, user.Username)
	if err != nil {
		switch err {
		case dbe.ErrDestinationBarcodeEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case dbe.ErrListOfBarcodesEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case dbe.ErrAtLeastOneBarcodeNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case dbe.ErrDestinationNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case dbe.ErrDestinationMultipleContainers:
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
