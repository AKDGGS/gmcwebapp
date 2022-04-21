package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (srv *Server) ServeWellsPointlist(w http.ResponseWriter, r *http.Request) {
	pl, err := srv.DB.GetWellsPointList()
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	// If no details were returned, throw a 404
	if pl == nil {
		http.Error(w, "Point list not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(pl)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("JSON error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	w.Write(js)
}
