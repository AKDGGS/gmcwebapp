package web

import (
	"fmt"
	"net/http"

	dbe "gmc/db/errors"
)

func (srv *Server) ServeAPIRecodeInventoryAndContainer(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("authentication error: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}
	if user == nil {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}
	q := r.URL.Query()
	old_barcode := q.Get("old")
	new_barcode := q.Get("new")
	err = srv.DB.RecodeInventoryAndContainer(old_barcode, new_barcode)
	if err != nil {
		switch err {
		case dbe.ErrOldBarcodeCannotBeEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case dbe.ErrNewBarcodeCannotBeEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case dbe.ErrNothingRecoded:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(
				w, fmt.Sprintf("error: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
		return
	}
}
