package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (srv *Server) ServeInventoryStash(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Authentication error: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}
	if user == nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	q := r.URL.Query()
	id, err := strconv.Atoi(q.Get("id"))
	if err != nil {
		http.Error(w, "Invalid Inventory ID", http.StatusBadRequest)
		return
	}
	stash, err := srv.DB.GetInventoryStash(id)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	out, err := compressWriter(r.Header.Get("Accept-Encoding"), w)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Compression error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	defer out.Close()

	jsenc := json.NewEncoder(out)
	w.Header().Set("Content-Type", "application/json")
	if err := jsenc.Encode(stash); err != nil {
		fmt.Fprintf(out, "\n\n%s", err.Error())
	}
}
