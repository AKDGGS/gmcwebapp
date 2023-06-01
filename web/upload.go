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

func (srv *Server) ServeUpload(w http.ResponseWriter, r *http.Request) int {
	var file_id int
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

		borehole_id := r.FormValue("boreholeId")
		borehole_id_int, err := strconv.Atoi(borehole_id)
		if err != nil {
			os.Exit(1)
		}
		f.BoreholeIDs = append(f.BoreholeIDs, borehole_id_int)

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
