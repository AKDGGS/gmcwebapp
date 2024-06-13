package web

import (
	"fmt"
	"net/http"
)

func (srv *Server) ServeAPIAudit(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("authentication error: %s", err),
			http.StatusBadRequest,
		)
		return
	}
	if user == nil {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}
	q := r.URL.Query()
	err = srv.DB.AddAudit(q.Get("remark"), q["c"])
	if err != nil {
		http.Error(
			w, fmt.Sprintf("add audit error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
}
