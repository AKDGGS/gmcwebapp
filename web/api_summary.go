package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	dbf "gmc/db/flag"
)

func (srv *Server) ServeAPISummary(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("authentication error: %s", err),
			http.StatusBadRequest,
		)
		return
	}
	if user == nil {
		http.Error(
			w,
			"access denied",
			http.StatusForbidden,
		)
		return
	}
	flags := dbf.ALL
	if user == nil {
		flags = dbf.ALL_NOPRIVATE
	}
	q := r.URL.Query()
	barcode := strings.TrimSpace(q.Get("barcode"))
	if barcode == "" {
		http.Error(
			w,
			"barcode cannot be empty",
			http.StatusBadRequest,
		)
		return
	}
	sum, err := srv.DB.GetSummaryByBarcode(barcode, flags)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("get summary by barcode error: %s", err),
			http.StatusInternalServerError,
		)
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
			w,
			fmt.Sprintf("json marshal error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
