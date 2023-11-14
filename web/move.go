package web

import (
	"fmt"
	"net/http"
)

func (srv *Server) ServeMove(w http.ResponseWriter, r *http.Request) {
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
	q := r.URL.Query()
	dest := q.Get("d")
	container_list := q["c"]
	err = srv.DB.MoveByBarcode(dest, container_list, user.Username)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
}
