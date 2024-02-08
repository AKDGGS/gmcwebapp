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
		http.Error(w, "invalid file id", http.StatusBadRequest)
		return
	}

	// Fetch the file details from the database
	db_file, err := srv.DB.GetFile(id)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// Fetch the file from filestore
	fs_file, err := srv.FileStore.GetFile(fmt.Sprintf("%d/%s",
		db_file.ID, db_file.Name,
	))
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			http.Error(w, "file not found (filestore)", http.StatusNotFound)
		} else {
			http.Error(
				w, fmt.Sprintf("file fetch error: %s", err.Error()),
				http.StatusInternalServerError,
			)
		}
		return
	}
	defer fs_file.Content.Close()

	// Suggest filename to the browser
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("inline; filename=\"%s\"", fs_file.Name),
	)
	// Set the ETag if available
	if fs_file.ETag != "" {
		w.Header().Set("ETag", fs_file.ETag)
	}
	http.ServeContent(w, r, fs_file.Name, fs_file.LastModified, fs_file.Content)
	return
}
