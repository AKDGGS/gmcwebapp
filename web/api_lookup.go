package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	dbf "gmc/db/flag"
)

func (srv *Server) ServeAPILookup(w http.ResponseWriter, r *http.Request) {
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
	barcode := q.Get("barcode")
	inventory_count, container_count, err := srv.DB.GetBarcodeCounts(barcode)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("get barcode counts error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	var js []byte
	if container_count == 0 && inventory_count == 1 {
		invs, err := srv.DB.GetInventoryByBarcode(barcode, flags)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("get inventory by barcode error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
		// If no details are returned, return an empty JSON object
		if invs == nil {
			js = []byte("{}")
		} else {
			js, err = json.Marshal(invs)
		}
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("json marshal error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
	} else if container_count >= 1 || inventory_count > 1 {
		sum, err := srv.DB.GetSummaryByBarcode(barcode, flags)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("get summary by barcode error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
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
	}
	if js == nil {
		js = []byte("{}")
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
