package web

import (
	"fmt"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gmc/db/model"
	fsutil "gmc/filestore/util"
)

func (srv *Server) ServeUpload(w http.ResponseWriter, r *http.Request) int32 {
	var file_id int32
	err := r.ParseMultipartForm(33554432) // 32MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 0
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		http.Error(w, "Multipart form of file is nil", http.StatusBadRequest)
		return 0
	}

	for _, fh := range r.MultipartForm.File["file"] {
		file, err := fh.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 0
		}
		defer file.Close()

		// temporary code until we decide what to do with the MD5.
		rand.Seed(time.Now().UnixNano())
		MD5 := strconv.FormatInt(rand.Int63(), 10)

		mt := mime.TypeByExtension(filepath.Ext(fh.Filename))
		if mt == "" {
			mt = "application/octet-stream"
		}

		f := model.File{
			Name: fh.Filename,
			Size: fh.Size,
			MD5:  MD5,
			Type: mt,
		}

		if _, ok := r.Form["borehole_id"]; ok {

			borehole_id := r.FormValue("borehole_id")
			borehole_id_int, err := strconv.Atoi(borehole_id)
			if err != nil {
				os.Exit(1)
			}
			f.BoreholeIDs = append(f.BoreholeIDs, borehole_id_int)
		}

		if _, ok := r.Form["inventory_id"]; ok {
			inventory_id := r.FormValue("inventory_id")
			inventory_id_int, err := strconv.Atoi(inventory_id)
			if err != nil {
				os.Exit(1)
			}
			f.InventoryIDs = append(f.InventoryIDs, inventory_id_int)
		}

		if _, ok := r.Form["outcrop_id"]; ok {
			outcrop_id := r.FormValue("outcrop_id")
			outcrop_id_int, err := strconv.Atoi(outcrop_id)
			if err != nil {
				os.Exit(1)
			}
			f.OutcropIDs = append(f.OutcropIDs, outcrop_id_int)
		}

		if _, ok := r.Form["prospect_id"]; ok {
			prospect_id := r.FormValue("prospect_id")
			prospect_id_int, err := strconv.Atoi(prospect_id)
			if err != nil {
				os.Exit(1)
			}
			f.ProspectIDs = append(f.ProspectIDs, prospect_id_int)
		}

		if _, ok := r.Form["well_id"]; ok {
			well_id := r.FormValue("well_id")
			well_id_int, err := strconv.Atoi(well_id)
			if err != nil {
				os.Exit(1)
			}
			f.WellIDs = append(f.WellIDs, well_id_int)
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
				return fmt.Errorf("error putting file in filestore: %w", err)
			}
			return nil
		})
		if err != nil {
			http.Error(w, "Error putting file in database or filestore: "+
				err.Error(), http.StatusInternalServerError)
			return 0
		}
		file_id = f.ID
	}
	return file_id
}
