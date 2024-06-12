package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	dbe "gmc/db/errors"
	dbf "gmc/db/flag"
)

func (srv *Server) ServeAPIInventoryDetail(w http.ResponseWriter, r *http.Request) {
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
	flags := dbf.ALL
	if user == nil {
		flags = dbf.ALL_NOPRIVATE
	}
	q := r.URL.Query()
	barcode := q.Get("barcode")
	invs, err := srv.DB.GetInventoryByBarcode(barcode, flags)
	if err != nil {
		if err == dbe.ErrBarcodeNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(
				w, fmt.Sprintf("error: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
		return
	}
	var js []byte
	// If no details are returned, return an empty JSON object
	if invs == nil {
		js = []byte("{}")
	} else {
		js, err = json.Marshal(invs)
	}
	if err != nil {
		switch err {
		case dbe.ErrBarcodeCannotBeEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case dbe.ErrInventoryInsertFailed:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		case dbe.ErrInventoryQualityInsertFailed:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(
				w, fmt.Sprintf("error: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
