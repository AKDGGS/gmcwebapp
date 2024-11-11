package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	sutil "gmc/search/util"
)

func (srv *Server) ServeSearchInventoryJSON(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("authentication error: %s", err),
			http.StatusBadRequest,
		)
		return
	}

	params := &sutil.InventoryParams{}
	params.ParseQuery(r.URL.Query(), (user != nil))

	result, err := srv.Search.SearchInventory(params)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("search error: %s", err),
			http.StatusBadRequest,
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
	if err := jsenc.Encode(result); err != nil {
		fmt.Fprintf(out, "\n\n%s", err)
	}
}
