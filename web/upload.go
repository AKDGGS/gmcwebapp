package web

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"gmc/db/model"
	fsutil "gmc/filestore/util"
)

func (srv *Server) ServeUpload(w http.ResponseWriter, r *http.Request) {
	if srv.FileStore == nil {
		http.Error(
			w,
			"filestore not configured",
			http.StatusInternalServerError,
		)
		return
	}
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
			http.StatusUnauthorized,
		)
		return
	}
	err = r.ParseMultipartForm(33554432) // 32MB
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("parse error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		http.Error(
			w,
			"multipart form of file is nil",
			http.StatusBadRequest,
		)
		return
	}
	for _, fh := range r.MultipartForm.File["file"] {
		file, err := fh.Open()
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("failed to open file: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
		defer file.Close()
		mt := mime.TypeByExtension(filepath.Ext(fh.Filename))
		if mt == "" {
			mt = "application/octet-stream"
		}
		f := model.File{
			Name: fh.Filename,
			Size: fh.Size,
			Type: mt,
		}
		if id := r.FormValue("borehole_id"); id != "" {
			borehole_id, err := strconv.Atoi(id)
			if err != nil {
				http.Error(
					w,
					"invalid borehole_id",
					http.StatusBadRequest,
				)
				return
			}
			f.BoreholeIDs = append(f.BoreholeIDs, borehole_id)
		}
		if id := r.FormValue("inventory_id"); id != "" {
			inventory_id, err := strconv.Atoi(id)
			if err != nil {
				http.Error(
					w,
					"invalid inventory_id",
					http.StatusBadRequest,
				)
				return
			}
			f.InventoryIDs = append(f.InventoryIDs, inventory_id)
		}
		if id := r.FormValue("outcrop_id"); id != "" {
			outcrop_id, err := strconv.Atoi(id)
			if err != nil {
				http.Error(
					w,
					"invalid outcrop_id",
					http.StatusBadRequest,
				)
				return
			}
			f.OutcropIDs = append(f.OutcropIDs, outcrop_id)
		}
		if id := r.FormValue("prospect_id"); id != "" {
			prospect_id, err := strconv.Atoi(id)
			if err != nil {
				http.Error(
					w,
					"invalid prospect_id",
					http.StatusBadRequest,
				)
				return
			}
			f.ProspectIDs = append(f.ProspectIDs, prospect_id)
		}
		if id := r.FormValue("well_id"); id != "" {
			well_id, err := strconv.Atoi(id)
			if err != nil {
				http.Error(
					w,
					"invalid well_id",
					http.StatusBadRequest,
				)
				return
			}
			f.WellIDs = append(f.WellIDs, well_id)
		}
		err = srv.DB.PutFile(&f, func() error {
			//Add the file to the filestore
			err := srv.FileStore.PutFile(&fsutil.File{
				Name:         fmt.Sprintf("%d/%s", f.ID, fh.Filename),
				Size:         fh.Size,
				LastModified: time.Now(),
				ContentType:  mt,
				Content:      file,
			})
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("error putting file in database or filestore: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
		w.Header().Set("file_id", strconv.Itoa(int(f.ID)))
	}
}
