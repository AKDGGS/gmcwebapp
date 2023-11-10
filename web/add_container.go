package web

import (
	"fmt"
	"net/http"
)

func (srv *Server) ServeAddContainer(barcode string, alt_barcode string, name string, remark string, w http.ResponseWriter, r *http.Request) {
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
	err = srv.DB.AddContainer(barcode, alt_barcode, name, remark, user)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
}
