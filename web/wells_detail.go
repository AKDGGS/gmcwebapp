package web

import (
	"fmt"
	"net/http"

	dbf "gmc/db/flag"
)

func (srv *Server) ServeWellsDetailJSON(id int, w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Authentication error: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}
	flags := dbf.INVENTORY_SUMMARY
	if user != nil {
		flags |= dbf.PRIVATE
	}
	well, err := srv.DB.GetWell(id, flags)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	// If no details are returned, throw a 404
	if well == nil {
		http.Error(w, "Well ID not found", http.StatusNotFound)
		return
	}
	js, err := well.MarshalJSON()
	if err != nil {
		http.Error(
			w, fmt.Sprintf("JSON error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
