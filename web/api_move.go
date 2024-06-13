package web

import (
	"fmt"
	"net/http"
)

func (srv *Server) ServeAPIMoveInventoryAndContainers(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("authentication error: %s", err),
			http.StatusBadRequest,
		)
		return
	}
	if user == nil {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}
	q := r.URL.Query()
	dest := q.Get("d")
	barcodes_to_move := q["c"]
	err = srv.DB.MoveInventoryAndContainers(dest, barcodes_to_move, user.Username)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("move inventory and containers error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	js := []byte(`{"success":true}`)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
