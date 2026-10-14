package web

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"gmc/db/model"
	fsutil "gmc/filestore/util"
)

type UploadStatus struct {
	ID    int32  `json:"file_id,omitempty"`
	Name  string `json:"name,omitempty"`
	Size  int64  `json:"size,omitempty"`
	Error string `json:"error,omitempty"`
}

func (srv *Server) ServeUpload(w http.ResponseWriter, r *http.Request) {
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

	q := r.URL.Query()
	inventory_ids, err := strsToInts(q["inventory_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	borehole_ids, err := strsToInts(q["borehole_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	outcrop_ids, err := strsToInts(q["outcrop_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	prospect_ids, err := strsToInts(q["prospect_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	well_ids, err := strsToInts(q["well_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	barcode := q.Get("barcode")
	var description *string
	if d := q.Get("description"); d != "" {
		description = &d
	}

	var results []UploadStatus

	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "no files uploaded", http.StatusInternalServerError)
		return
	}
	for {
		pt, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("multipart error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}

		if fn := pt.FileName(); fn != "" {
			stat := UploadStatus{Name: fn}

			// Determine mime type, set default on failure
			mt := mime.TypeByExtension(filepath.Ext(fn))
			if mt == "" {
				mt = "application/octet-stream"
			}

			// Fetch a new ID
			id, err := srv.DB.GetFileNextID()
			if err != nil {
				stat.Error = err.Error()
				results = append(results, stat)
				continue
			}

			// Add file to file store
			fsfile := &fsutil.File{
				Name:         fmt.Sprintf("%d/%s", id, fn),
				Size:         -1,
				LastModified: time.Now(),
				ContentType:  mt,
			}
			fsfile.SetContent(pt, true)
			if err := srv.FileStores.PutFile(fsfile); err != nil {
				stat.Error = err.Error()
				results = append(results, stat)
				continue
			}

			hash := fmt.Sprintf("%x", fsfile.GetHash())

			f, err := srv.DB.GetFileByHash(hash)
			if err != nil {
				// Try to cleanup file from filestore
				srv.FileStores.DeleteFile(fmt.Sprintf("%d/%s", id, fn))
				stat.Error = err.Error()
				results = append(results, stat)
				continue
			}

			if f == nil {
				// If there's no existing file, add the file
				f = &model.File{
					ID:           id,
					Name:         fn,
					Description:  description,
					Size:         fsfile.GetReadSize(),
					Hash:         &hash,
					MimeType:     mt,
					Barcode:      barcode,
					InventoryIDs: inventory_ids,
					BoreholeIDs:  borehole_ids,
					OutcropIDs:   outcrop_ids,
					ProspectIDs:  prospect_ids,
					WellIDs:      well_ids,
				}
				if err := srv.DB.AddFile(f); err != nil {
					// Try to cleanup file from filestore
					srv.FileStores.DeleteFile(
						fmt.Sprintf("%d/%s", stat.ID, stat.Name),
					)
					stat.Error = err.Error()
					results = append(results, stat)
					continue
				}
			} else {
				// If the file already exists, delete the new upload
				srv.FileStores.DeleteFile(fmt.Sprintf("%d/%s", id, fn))

				// Update the existing file
				f.Description = description
				f.Barcode = barcode
				f.InventoryIDs = inventory_ids
				f.BoreholeIDs = borehole_ids
				f.OutcropIDs = outcrop_ids
				f.ProspectIDs = prospect_ids
				f.WellIDs = well_ids
				if err := srv.DB.UpdateFile(f); err != nil {
					// Do not try to clean up the file here, it already exists
					stat.Error = err.Error()
					results = append(results, stat)
					continue
				}
			}

			stat.ID = f.ID
			stat.Size = f.Size
			results = append(results, stat)
		}
	}

	jsn, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsn)
}

func strsToInts(arr []string) ([]int, error) {
	if arr == nil {
		return []int{}, nil
	}
	ns := make([]int, len(arr))
	for i, s := range arr {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %s", s)
		}
		ns[i] = n
	}
	return ns, nil
}
