package web

import (
	"encoding/json"
	"fmt"
	dbf "gmc/db/flag"
	"net/http"
)

func (srv *Server) ServeWellJSON(id int, w http.ResponseWriter, r *http.Request) {
	flags := dbf.ALL_NOPRIVATE
	welljson, err := srv.DB.GetWellJSON(id, flags)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// If no details were returned, throw a 404
	if welljson == nil {
		http.Error(w, "Well ID not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	js, err := json.Marshal(welljson)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("JSON error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	w.Write(js)
}
