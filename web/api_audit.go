package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (srv *Server) ServeAPIAudit(w http.ResponseWriter, r *http.Request) {
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
	}
	q := r.URL.Query()
	barcode := strings.TrimSpace(q.Get("barcode"))
	if barcode == "" {
		http.Error(
			w,
			"parameter error: barcode can't be empty",
			http.StatusBadRequest,
		)
		return
	} else {
		results, err := srv.DB.Audit(barcode)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("add audit error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
		var js []byte
		// If no details are returned, return an empty JSON object
		if results == nil {
			js = []byte("{}")
		} else {
			js, err = json.Marshal(results)
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
}
