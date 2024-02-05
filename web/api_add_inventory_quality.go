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
		switch err {
		case dbe.ErrNotFoundInInventory:
			http.Error(w, err.Error(), http.StatusNotFound)
		case dbe.ErrMultipleIDs:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		case dbe.ErrInventoryQualityInsertFailed:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			http.Error(
				w, fmt.Sprintf("Error: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
		return
	}
}
