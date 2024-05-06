package web

import (
	"fmt"
	"net/http"

	dbe "gmc/db/errors"
)

func (srv *Server) ServeAPIAddInventory(w http.ResponseWriter, r *http.Request) {
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
	issues := q["i"]
	err = srv.DB.AddInventory(q.Get("barcode"), q.Get("remark"), container_id, issues, user.Username)
	if err != nil {
		if err == dbe.ErrInventoryInsertFailed {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if err == dbe.ErrInventoryQualityInsertFailed {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(
				w, fmt.Sprintf("Error: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
		return
	}
	js := []byte(`{"success":true}`)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
