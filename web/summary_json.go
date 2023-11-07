package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	dbf "gmc/db/flag"
)

func (srv *Server) ServeSummaryDetail(barcode string, w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Authentication error: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}
	if user == nil {
		http.Error(w, "Authentication error", http.StatusForbidden)
		return
	}
	flags := dbf.ALL
	if user == nil {
		flags = dbf.ALL_NOPRIVATE
	}
	sum, err := srv.DB.GetSummaryByBarcode(barcode, flags)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// If no details are returned, throw a 404
	if sum == nil {
		http.Error(w, "Summary not found", http.StatusNotFound)
		return
	}

	js, err := json.Marshal(sum)
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
