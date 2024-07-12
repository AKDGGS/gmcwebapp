package util

import (
	"fmt"
	"net/http"
)

// Similar to http.Redirect() except it allows for relative URLs
func Redirect(w http.ResponseWriter, url string, code int) {
	w.Header().Set("Location", url)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(
		w, "<a href=\"%s\">%s</a>\n",
		url, http.StatusText(code),
	)
}
