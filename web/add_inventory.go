package web

import (
	"fmt"
	"net/http"
)

func (srv *Server) ServeAddInventory(w http.ResponseWriter, r *http.Request) {
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
	q := r.URL.Query()
	var container_id *int32
	issues, ok := q["i"]
	if !ok {
		http.Error(w, "Issues list cannot be empty", http.StatusInternalServerError)
		return
	}
	err = srv.DB.AddInventory(q.Get("barcode"), q.Get("remark"), container_id, issues, user.Username)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
}
