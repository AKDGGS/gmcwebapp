package web

import (
	"fmt"
	"net/http"
	"strings"
)

func (srv *Server) ServeAPIAddInventory(w http.ResponseWriter, r *http.Request) {
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
	barcode := strings.TrimSpace(q.Get("barcode"))
	if barcode == "" {
		http.Error(
			w,
			"parameter error: barcode cannot be empty",
			http.StatusBadRequest,
		)
		return
	}
	var container_id *int32
	issues := q["i"]
	err = srv.DB.AddInventory(barcode, q.Get("remark"), container_id, issues, user.Username)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("add inventory error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	js := []byte(`{"success":true}`)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
