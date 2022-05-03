package web

import (
	"encoding/json"
	"fmt"
	dbf "gmc/db/flag"
	"net/http"
	"time"
)

func (srv *Server) ServeWellJSON(id int, w http.ResponseWriter, r *http.Request) {
	flags := dbf.ALL_NOPRIVATE
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

	cd, ok := welljson["completion_date"].(time.Time)
	if !ok {
		delete(welljson, "completion_date")
	} else {
		fcd := cd.Format("01-02-2006")
		welljson["completion_date"] = &fcd
	}

	sd, ok := welljson["spud_date"].(time.Time)
	if !ok {
		delete(welljson, "spud_date")
	} else {
		fsd := sd.Format("01-02-2006")
		welljson["spud_date"] = &fsd
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
