package web

import (
	"fmt"
	"net/http"
)

func (srv *Server) ServeAPIAddInventoryQuality(w http.ResponseWriter, r *http.Request) {
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
	err = srv.DB.AddInventoryQuality(q.Get("barcode"), q.Get("remark"), q["i"], user.Username)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("add inventory quality error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
}
