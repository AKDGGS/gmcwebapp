package web

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func (srv *Server) ServeFile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			"invalid file id",
			http.StatusBadRequest,
		)
		return
	}
	// Fetch the file details from the database
	db_file, err := srv.DB.GetFile(id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("get file error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	if db_file == nil {
		http.Error(
			w,
			"file does not exist",
			http.StatusNotFound,
		)
		return
	}
	// Fetch the file from filestore
	fsfile, err := srv.FileStores.GetFile(fmt.Sprintf("%d/%s",
		db_file.ID, db_file.Name,
	))
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			http.Error(
				w,
				"get file error: file not found in filestore",
				http.StatusNotFound,
			)
		} else {
			http.Error(
				w,
				fmt.Sprintf("get file fetching error: %s", err),
				http.StatusInternalServerError,
			)
		}
		return
	}
	defer fsfile.Close()
	// Suggest filename to the browser
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("inline; filename=\"%s\"", fsfile.Name),
	)
	// Set the ETag if available
	if fsfile.ETag != "" {
		w.Header().Set("ETag", fsfile.ETag)
	}
	http.ServeContent(w, r, fsfile.Name, fsfile.LastModified, fsfile)
	return
}
