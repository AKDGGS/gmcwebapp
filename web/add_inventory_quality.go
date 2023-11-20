package web

import (
	"fmt"
	"net/http"
)

func (srv *Server) ServeAddInventoryQuality(w http.ResponseWriter, r *http.Request) {
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
	barcode := q.Get("barcode")
	remark := q.Get("remark")
	issues := q["i"]
	err = srv.DB.AddInventoryQuality(barcode, remark, issues, user.Username)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
}
