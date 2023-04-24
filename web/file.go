package web

import (
	"fmt"
	"net/http"
	"os"
)

func (srv *Server) ServeFile(id int, w http.ResponseWriter, r *http.Request) {
	// Fetch the file details from the database
	db_file, err := srv.DB.GetFile(id)
	fmt.Println(db_file)
	if err != nil {
		http.Error(
			w, fmt.Sprintf("Query error: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}

	// Fetch the file from filestore
	fs_file, err := srv.FileStore.GetFile(fmt.Sprintf("%d/%s",
		db_file.ID, db_file.Name,
	))

	fmt.Println(fs_file, fmt.Sprintf("%d/%s",
		db_file.ID, db_file.Name,
	))
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			http.Error(w, "File not found (FileStore)", http.StatusNotFound)
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
