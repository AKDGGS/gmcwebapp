package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	dbf "gmc/db/flag"
)

func (srv *Server) ServeWellsDetail(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("authentication error: %s", err),
			http.StatusBadRequest,
		)
		return
	}
	flags := dbf.INVENTORY_SUMMARY
	if user != nil {
		flags |= dbf.PRIVATE
	}
	q := r.URL.Query()
	id, err := strconv.Atoi(q.Get("id"))
	if err != nil {
		http.Error(
			w,
			"invalid well id",
			http.StatusBadRequest,
		)
		return
	}
	well, err := srv.DB.GetWell(id, flags)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("get well error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	// If no details are returned, throw a 404
	if well == nil {
		http.Error(
			w,
			"well id not found",
			http.StatusNotFound,
		)
		return
	}
	js, err := json.Marshal(well)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("json marshal error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
