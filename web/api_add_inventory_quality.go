package web

import (
	"fmt"
	"net/http"

	dbe "gmc/db/errors"
)

func (srv *Server) ServeAPIAddInventoryQuality(w http.ResponseWriter, r *http.Request) {
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
	err = srv.DB.AddInventoryQuality(q.Get("barcode"), q.Get("remark"), q["i"], user.Username)
	if err != nil {
		if err == dbe.ErrNotFoundInInventory {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err == dbe.ErrMultipleIDs {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if err == dbe.ErrInventoryQualityInsertFailed {
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
