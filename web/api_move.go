package web

import (
	"fmt"
	"net/http"
	"strings"
)

func (srv *Server) ServeAPIMoveInventoryAndContainers(w http.ResponseWriter, r *http.Request) {
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
	dest := strings.TrimSpace(q.Get("d"))
	if dest == "" {
		http.Error(
			w,
			"destination barcode cannot be empty",
			http.StatusBadRequest,
		)
		return
	}
	barcodes_to_move := q["c"]
	if barcodes_to_move == nil || len(barcodes_to_move) < 1 {
		http.Error(
			w,
			"list of barcodes cannot be empty",
			http.StatusBadRequest,
		)
		return
	}
	err = srv.DB.MoveInventoryAndContainers(dest, barcodes_to_move, user.Username)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("move inventory and containers error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	js := []byte(`{"success":true}`)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
