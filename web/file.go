package web

import (
	"fmt"
	"net/http"
	"os"
)

func (srv *Server) ServeFile(fp string, w http.ResponseWriter, r *http.Request) {
	// Fetch the file from S3
	file, err := srv.FileStore.GetFile(fp)
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
	defer file.Content.Close()

	// Suggest filename to the browser
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("inline; filename=\"%s\"", file.Name),
	)
	// Set the ETag if available
	if file.ETag != "" {
		w.Header().Set("ETag", file.ETag)
	}
	http.ServeContent(w, r, file.Name, file.LastModified, file.Content)
	return
}
