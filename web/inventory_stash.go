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
			w,
			fmt.Sprintf("authentication error: %s", err),
			http.StatusBadRequest,
		)
		return
	}
	if user == nil {
		http.Error(
			w,
			"access denied",
			http.StatusForbidden,
		)
		return
	}
	q := r.URL.Query()
	id, err := strconv.Atoi(q.Get("id"))
	if err != nil {
		http.Error(
			w,
			"invalid inventory id",
			http.StatusBadRequest,
		)
		return
	}
	stash, err := srv.DB.GetInventoryStash(id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("get inventory stash error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	out, err := compressWriter(r.Header.Get("Accept-Encoding"), w)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("compression error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	defer out.Close()
	jsenc := json.NewEncoder(out)
	w.Header().Set("Content-Type", "application/json")
	if err := jsenc.Encode(stash); err != nil {
		fmt.Fprintf(out, "\n\n%s", err)
	}
}
