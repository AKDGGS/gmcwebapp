package web

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"gmc/assets"
	"gmc/db/model"
)

func (srv *Server) ServeAPIGetFileListPartial(w http.ResponseWriter, r *http.Request) {
	var files []model.File
	q := r.URL.Query()
	var tbl int
	var sid string
	if q.Has("borehole_id") {
		tbl = model.BOREHOLE
		sid = q.Get("borehole_id")
	} else if q.Has("inventory_id") {
		tbl = model.INVENTORY
		sid = q.Get("inventory_id")
	} else if q.Has("outcrop_id") {
		tbl = model.OUTCROP
		sid = q.Get("outcrop_id")
	} else if q.Has("prospect_id") {
		tbl = model.PROSPECT
		sid = q.Get("prospect_id")
	} else if q.Has("well_id") {
		tbl = model.WELL
		sid = q.Get("well_id")
	} else {
		http.Error(
			w, "id required", http.StatusInternalServerError,
		)
		return
	}
	id, err := strconv.Atoi(sid)
	if err != nil {
		http.Error(
			w,
			"invalid id",
			http.StatusInternalServerError,
		)
		return
	}
	files, err = srv.DB.ListFilesById(tbl, id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("fetch partial error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	buf := bytes.Buffer{}
	if files != nil {
		if err := assets.ExecuteTemplate("tmpl/filelist.partial.html", &buf, files); err != nil {
			http.Error(
				w,
				fmt.Sprintf("parse error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
	out.Write(buf.Bytes())
}
