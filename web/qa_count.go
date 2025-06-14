package web

import (
	"fmt"
	"net/http"
	"strconv"
)

func (srv *Server) ServeQACount(w http.ResponseWriter, r *http.Request) {
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
			"invalid report id",
			http.StatusBadRequest)
		return
	}
	c, err := srv.DB.CountQAReport(id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("count qa report error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	// Numbers are technically JSON
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%d", c)
}
