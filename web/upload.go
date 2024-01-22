package web

import (
	"fmt"
	"math/rand"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"gmc/db/model"
	fsutil "gmc/filestore/util"
)

func (srv *Server) ServeUpload(w http.ResponseWriter, r *http.Request) error {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Authentication error: %s", err.Error()),
			http.StatusBadRequest,
		)
		return err
	}
	if user == nil {
		http.Error(w, "Access denied", http.StatusUnauthorized)
		return err
	}
	err = r.ParseMultipartForm(33554432) // 32MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		http.Error(w, "Multipart form of file is nil", http.StatusBadRequest)
		return err
	}
	file, fh, err := r.FormFile("content")
	if err != nil {
		return err
	}
	defer file.Close()
	// temporary code until we decide what to do with the MD5.
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	MD5 := strconv.FormatInt(random.Int63(), 10)

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
	if bid := r.FormValue("borehole_id"); bid != "" {
		borehole_id_int, err := strconv.Atoi(bid)
		if err != nil {
			http.Error(w, "Invalid Borehole ID", http.StatusBadRequest)
			return err
		}
		f.BoreholeIDs = append(f.BoreholeIDs, borehole_id_int)
	}
	if iid := r.FormValue("inventory_id"); iid != "" {
		inventory_id_int, err := strconv.Atoi(iid)
		if err != nil {
			http.Error(w, "Invalid Inventory ID", http.StatusBadRequest)
			return err
		}
		f.InventoryIDs = append(f.InventoryIDs, inventory_id_int)
	}
	if oid := r.FormValue("outcrop_id"); oid != "" {
		outcrop_id_int, err := strconv.Atoi(oid)
		if err != nil {
			http.Error(w, "Invalid Outcrop ID", http.StatusBadRequest)
			return err
		}
		f.OutcropIDs = append(f.OutcropIDs, outcrop_id_int)
	}
	if pid := r.FormValue("prospect_id"); pid != "" {
		prospect_id_int, err := strconv.Atoi(pid)
		if err != nil {
			http.Error(w, "Invalid Prospect ID", http.StatusBadRequest)
			return err
		}
		f.ProspectIDs = append(f.ProspectIDs, prospect_id_int)
	}
	if wid := r.FormValue("well_id"); wid != "" {
		well_id_int, err := strconv.Atoi(wid)
		if err != nil {
			http.Error(w, "Invalid Well ID", http.StatusBadRequest)
			return err
		}
		f.WellIDs = append(f.WellIDs, well_id_int)
	}
	if b := r.FormValue("barcode"); b != "" {
		f.Barcodes = append(f.Barcodes, b)
	}
	if d := r.FormValue("description"); d != "" {
		f.Description = &d
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
			return fmt.Errorf("Error putting file in filestore: %w", err)
		}
		return nil
	})
	if err != nil {
		http.Error(w, "Error putting file in database or filestore: "+
			err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Header().Set("file_id", strconv.Itoa(int(f.ID)))
	return nil
}
