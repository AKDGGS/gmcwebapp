package web

import (
	"encoding/csv"
	"fmt"
	"net/http"

	"gmc/db/model"
	sutil "gmc/search/util"
)

func (srv *Server) ServeSearchInventoryCSV(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("authentication error: %s", err),
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

	cwr := csv.NewWriter(out)
	defer cwr.Flush()

	w.Header().Set("Content-Type", "text/csv")

	cwr.Write(model.FlatInventoryFields(user != nil))

	params := &sutil.InventoryParams{}
	params.ParseQuery(r.URL.Query(), (user != nil))
	params.IncludeDescription = true
	params.From = 0
	params.Size = 10000

	for {
		result, err := srv.Search.SearchInventory(params)
		if err != nil {
			fmt.Fprintf(out, "\r\n\r\nsearch error: %s", err)
			return
		}

		for _, h := range result.Hits {
			cwr.Write(h.AsStringArray(user != nil))
		}

		params.From = (result.From + len(result.Hits))
		if int64(params.From) >= result.Total {
			return
		}
	}
}
