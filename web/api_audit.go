package web

import (
	"fmt"
	"net/http"

	dbe "gmc/db/errors"
)

func (srv *Server) ServeAPIAudit(w http.ResponseWriter, r *http.Request) {
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
	err = srv.DB.AddAudit(q.Get("remark"), q["c"])
	if err != nil {
		if err == dbe.ErrAuditParamsEmpty {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err == dbe.ErrAuditInsertFailed {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(
				w, fmt.Sprintf("Error: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
		return
	}
}
