package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

	q := r.URL.Query()

	size, err := strconv.Atoi(q.Get("size"))
	if err != nil {
		size = 25
	} else if user == nil && size > 1000 {
		size = 25
	} else if user != nil && size > 10000 {
		size = 25
	}

	from, err := strconv.Atoi(q.Get("from"))
	if err != nil || from < 0 {
		from = 0
	}

	params := &sutil.InventoryParams{
		Query:   q.Get("q"),
		Size:    size,
		From:    from,
		Private: (user != nil),
	}

	if t := q.Get("top"); t != "" {
		if n, err := strconv.ParseFloat(t, 64); err == nil {
			params.IntervalTop = &n
		}
	}
	if t := q.Get("bottom"); t != "" {
		if n, err := strconv.ParseFloat(t, 64); err == nil {
			params.IntervalBottom = &n
		}
	}
	// Save people from flipping top and bottom values
	if params.IntervalTop != nil && params.IntervalBottom != nil {
		if *params.IntervalTop > *params.IntervalBottom {
			t := params.IntervalTop
			params.IntervalTop = params.IntervalBottom
			params.IntervalBottom = t
		}
	}

	if keywords, ok := q["keyword"]; ok {
		params.Keywords = keywords
	}

	if ids, ok := q["prospect_id"]; ok {
		for _, sid := range ids {
			if id, err := strconv.Atoi(sid); err == nil {
				params.ProspectIDs = append(params.ProspectIDs, id)
			}
		}
	}

	if ids, ok := q["collection_id"]; ok {
		for _, sid := range ids {
			if id, err := strconv.Atoi(sid); err == nil {
				params.CollectionIDs = append(params.CollectionIDs, id)
			}
		}
	}

	if sorts, ok := q["sort"]; ok {
		dirs, _ := q["dir"]
		for i, v := range sorts {
			dir := "asc"
			if i < len(dirs) && strings.ToLower(dirs[i]) == "desc" {
				dir = "desc"
			}
			params.Sort = append(params.Sort, [2]string{v, dir})
		}
	}

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
