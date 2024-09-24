package web

import (
	"fmt"
	"net/http"
	"strings"
)

func (srv *Server) ServeAPIRecodeInventoryAndContainer(w http.ResponseWriter, r *http.Request) {
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
	q := r.URL.Query()
	old_barcode := strings.TrimSpace(q.Get("old"))
	if old_barcode == "" {
		http.Error(
			w,
			"old barcode cannot be empty",
			http.StatusBadRequest,
		)
		return
	}
	new_barcode := strings.TrimSpace(q.Get("new"))
	if new_barcode == "" {
		http.Error(
			w,
			"new barcode cannot be empty",
			http.StatusBadRequest,
		)
		return
	}
	err = srv.DB.RecodeInventoryAndContainer(old_barcode, new_barcode)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("recode inventory and container error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	js := []byte(`{"success":true}`)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
