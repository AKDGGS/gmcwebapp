package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	dbe "gmc/db/errors"
	dbf "gmc/db/flag"
)

func (srv *Server) ServeAPISummary(w http.ResponseWriter, r *http.Request) {
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
	flags := dbf.ALL
	if user == nil {
		flags = dbf.ALL_NOPRIVATE
	}
	q := r.URL.Query()
	barcode := q.Get("barcode")
	sum, err := srv.DB.GetSummaryByBarcode(barcode, flags)
	if err != nil {
		if err == dbe.ErrBarcodeNotContainer {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(
				w, fmt.Sprintf("Error: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
		return
	}
	var js []byte
	// If no details are returned, return an empty JSON object
	if sum == nil {
		js = []byte("{}")
	} else {
		js, err = json.Marshal(sum)
	}
	if err != nil {
		http.Error(
			w, fmt.Sprintf("JSON error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
