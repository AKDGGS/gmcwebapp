package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	welljson, err := srv.DB.GetWell(id, flags)
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

	if cd, ok := welljson["completion_date"].(time.Time); ok {
		welljson["completion_date"] = cd.Format("01-02-2006")
	} else {
		delete(welljson, "completion_date")
	}

	if sd, ok := welljson["spud_date"].(time.Time); ok {
		welljson["spud_date"] = sd.Format("01-02-2006")
	} else {
		delete(welljson, "spud_datee")
	}

	js, err := json.Marshal(welljson)
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
