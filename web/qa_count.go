package web

import (
	"fmt"
	"net/http"
)

func (srv *Server) ServeQACount(id int, w http.ResponseWriter, r *http.Request) {
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

	c, err := srv.DB.CountQAReport(id)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// Numbers are technically JSON
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%d", c)
}
