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
			w, fmt.Sprintf("authentication error: %s", err.Error()),
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
		switch err {
		case dbe.ErrAuditParamsEmpty:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case dbe.ErrAuditInsertFailed:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(
				w, fmt.Sprintf("error: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
		return
	}
}
