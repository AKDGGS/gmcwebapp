package web

import (
	"fmt"
	dbf "gmc/db/flag"
	"net/http"
)

func (srv *Server) ServeMove(dest string, container_list []string, w http.ResponseWriter, r *http.Request) {
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

	flags := dbf.ALL
	if user == nil {
		flags = dbf.ALL_NOPRIVATE
	}

	err = srv.DB.MoveByBarcode(dest, container_list, flags, user)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
}
