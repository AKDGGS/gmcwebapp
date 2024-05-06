package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (srv *Server) ServeAPIQualityList(w http.ResponseWriter, r *http.Request) {
	list, err := srv.DB.ListIssues()
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}