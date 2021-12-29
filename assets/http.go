package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
)

var assetETags map[string]string
var gzAssets map[string][]byte
var gzAssetETags map[string]string

func initStatic() error {
	assetETags = make(map[string]string)
	gzAssets = make(map[string][]byte)
	gzAssetETags = make(map[string]string)

	patterns := []string{"img/*", "css/*", "ol/*", "js/*"}
	for _, pattern := range patterns {
		files, err := fs.Glob(assets, pattern)
		if err != nil {
			return err
		}

		for _, fn := range files {
			f := ReadBytes(fn)

			// Throw an error if asked to process zero length file
			if len(f) == 0 {
				return fmt.Errorf("cannot process empty file: %s", fn)
			}

			assetETags[fn] = fmt.Sprintf("%x", md5.Sum(f))

			// 860 bytes is Akamai's recommended minimum for gzip
			// so don't bother to gzip files less than 860 bytes
			if len(f) < 860 {
				continue
			}

			var buf bytes.Buffer
			gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
			if err != nil {
				return err
			}
			defer gz.Close()

			if _, err := gz.Write(f); err != nil {
				return err
			}

			if err := gz.Flush(); err != nil {
				return err
			}

			// Only accept gzip if it's less than the original in size
			if buf.Len() > 0 && buf.Len() < len(f) {
				gzAssets[fn] = buf.Bytes()
				gzAssetETags[fn] = fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
			}
		}
	}
	return nil
}

func ServeStatic(name string, w http.ResponseWriter, r *http.Request) {
	// Is it okay to enable gzip?
	gzok := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

	var content []byte
	var localETag string
	if gzok {
		content, _ = gzAssets[name]
		localETag, _ = gzAssetETags[name]
	}
	if len(content) == 0 {
		content = ReadBytes(name)
		localETag, _ = assetETags[name]
		gzok = false
	}

	if len(content) == 0 {
		http.Error(w, "File not found (Static)", http.StatusNotFound)
		return
	}

	// ETag if file hasn't changed
	remoteETag := r.Header.Get("If-None-Match")
	if localETag == remoteETag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("ETag", localETag)
	contenttype := mime.TypeByExtension(path.Ext(name))
	if contenttype == "" {
		contenttype = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contenttype)

	if gzok {
		w.Header().Set("Content-Encoding", "gzip")
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
	w.Write(content)
}
