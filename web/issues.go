package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (srv *Server) ServeIssues(w http.ResponseWriter, r *http.Request) {
	list, err := srv.DB.ListIssues()
	if err != nil {
		http.Error(
			w, fmt.Sprintf("error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
