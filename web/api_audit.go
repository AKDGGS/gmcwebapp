package web

import (
	"fmt"
	"net/http"
	"strings"
)

func (srv *Server) ServeAPIAudit(w http.ResponseWriter, r *http.Request) {
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
	}
	q := r.URL.Query()
	remark := strings.TrimSpace(q.Get("remark"))
	var js []byte
	if remark == "" || len(q["c"]) == 0 {
		http.Error(
			w,
			"parameter error: remark and container list can't be empty",
			http.StatusBadRequest,
		)
		return
	} else {
		err = srv.DB.AddAudit(remark, q["c"])
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("add audit error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
		js = []byte(`{"success":true}`)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
		w.Write(js)
	}
}
