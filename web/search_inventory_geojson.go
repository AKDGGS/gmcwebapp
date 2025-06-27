package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gmc/db/model"
	sutil "gmc/search/util"
)

func (srv *Server) ServeSearchInventoryGeoJSON(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("authentication error: %s", err),
			http.StatusBadRequest,
		)
		return
	}

	out, err := compressWriter(r.Header.Get("Accept-Encoding"), w)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("compression error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
	defer out.Close()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "search.geojson"))

	params := &sutil.InventoryParams{}
	params.ParseQuery(r.URL.Query(), (user != nil))
	params.IncludeDescription = true
	params.IncludeLatLon = true
	params.From = 0
	params.Size = 10000

	_, err = out.Write([]byte("{\"type\": \"FeatureCollection\",\"features\":["))
	if err != nil {
		fmt.Fprintf(out, "\r\n\r\nwrite error: %s", err)
		return
	}

	for {
		result, err := srv.Search.SearchInventory(params)
		if err != nil {
			fmt.Fprintf(out, "\r\n\r\nsearch error: %s", err)
			return
		}

		for i, h := range result.Hits {
			var data []byte
			var feature struct {
				Type     string `json:"type"`
				Geometry *struct {
					Type   string          `json:"type"`
					Coords json.RawMessage `json:"coordinates"`
				} `json:"geometry"`
				model.FlatInventory `json:"properties,omitempty"`
			}
			feature.Type = "Feature"
			feature.FlatInventory = h
			if len(h.Geometries) > 0 {
				if err = json.Unmarshal(h.Geometries[0], &feature.Geometry); err != nil {
					fmt.Fprintf(out, "\r\n\r\ngeometry unmarshal error: %s", err)
					return
				}
				data, err = json.Marshal(feature)
				if err != nil {
					fmt.Fprintf(out, "\r\n\r\nfeature marshal error: %s", err)
					return
				}
				_, err = out.Write(data)
				if err != nil {
					fmt.Fprintf(out, "\r\n\r\nfeature write error: %s", err)
					return
				}
			} else {
				data, err = json.Marshal(feature)
				if err != nil {
					fmt.Fprintf(out, "\r\n\r\nfeature marshal error: %s", err)
					return
				}
				_, err = out.Write(data)
				if err != nil {
					fmt.Fprintf(out, "\r\n\r\nfeature write error: %s", err)
					return
				}
			}
			if result.From+i != int(result.Total)-1 {
				_, err = out.Write([]byte(","))
				if err != nil {
					fmt.Fprintf(out, "\r\n\r\nfeature write error: %s", err)
					return
				}
			}
		}

		params.From = (result.From + len(result.Hits))
		if int64(params.From) >= result.Total {
			_, err = out.Write([]byte("]}"))
			if err != nil {
				fmt.Fprintf(out, "\r\n\r\nwrite error: %s", err)
				return
			}
			return
		}
	}
}
