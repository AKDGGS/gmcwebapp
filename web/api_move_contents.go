package web

import (
	"fmt"
	"net/http"
	"strings"
)

func (srv *Server) ServeAPIMoveInventoryAndContainersContents(w http.ResponseWriter, r *http.Request) {
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
	src := strings.TrimSpace(q.Get("src"))
	if src == "" {
		http.Error(
			w,
			"source barcode cannot be empty",
			http.StatusBadRequest,
		)
		return
	}
	dest := strings.TrimSpace(q.Get("dest"))
	if dest == "" {
		http.Error(
			w,
			"destination barcode cannot be empty",
			http.StatusBadRequest,
		)
		return
	}
	err = srv.DB.MoveInventoryAndContainersContents(src, dest)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("move inventory and container contents error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	js := []byte(`{"success":true}`)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(js)))
	w.Write(js)
}
