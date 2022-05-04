package web

import (
	"encoding/json"
	"fmt"
	dbf "gmc/db/flag"
	"net/http"
	"time"
)

func (srv *Server) ServeWellJSON(id int, w http.ResponseWriter, r *http.Request) {
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
		fcd := cd.Format("01-02-2006")
		welljson["completion_date"] = &fcd
	} else {
		delete(welljson, "completion_date")
	}

	if sd, ok := welljson["spud_date"].(time.Time); ok {
		fsd := sd.Format("01-02-2006")
		welljson["spud_date"] = &fsd
	} else {
		delete(welljson, "spud_datee")
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
