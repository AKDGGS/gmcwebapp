package web

import (
	"fmt"
	"net/http"
)

func (srv *Server) ServeAddInventory(barcode string, remark string, container_id *int32, issues []string, w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Authentication error: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}
	if user == nil {
		http.Error(w, "Access denied.", http.StatusForbidden)
		return
	}
	err = srv.DB.AddInventory(barcode, remark, container_id, issues, user)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
}
