package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"gmc/db/model"
	sutil "gmc/search/util"

	"github.com/parquet-go/parquet-go"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkb"
)

type InventoryItem struct {
	model.FlatInventory
	Geometry []byte `parquet:"geometry,uncompressed"`
}

func (srv *Server) ServeSearchInventoryParquet(w http.ResponseWriter, r *http.Request) {
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

	params := &sutil.InventoryParams{}
	params.ParseQuery(r.URL.Query(), (user != nil))
	params.IncludeDescription = true
	params.IncludeLatLon = true
	params.From = 0
	params.Size = 10000

	writer := parquet.NewGenericWriter[InventoryItem](
		out,
		parquet.Compression(&parquet.Uncompressed),
	)
	defer writer.Close()

	geometrytypes := []string{}

	w.Header().Set("Content-Type", "application/vnd.apache.parquet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "search.parquet"))

	for {
		result, err := srv.Search.SearchInventory(params)
		if err != nil {
			fmt.Fprintf(out, "\r\n\r\nsearch error: %s", err)
			return
		}

		itemsbatch := make([]InventoryItem, 0)

		for _, h := range result.Hits {
			var wkbbytes []byte
			if len(h.Geometries) > 0 {
				var gt geom.T
				if err = geojson.Unmarshal(h.Geometries[0], &gt); err != nil {
					fmt.Fprintf(out, "\r\n\r\ngeometry unmarshal error: %s", err)
					return
				}
				switch gt.(type) {
				case *geom.Point:
					if !slices.Contains(geometrytypes, "Point") {
						geometrytypes = append(geometrytypes, "Point")
					}
				case *geom.MultiPoint:
					if !slices.Contains(geometrytypes, "MultiPoint") {
						geometrytypes = append(geometrytypes, "MultiPoint")
					}
				case *geom.Polygon:
					if !slices.Contains(geometrytypes, "Polygon") {
						geometrytypes = append(geometrytypes, "Polygon")
					}
				case *geom.MultiPolygon:
					if !slices.Contains(geometrytypes, "MultiPolygon") {
						geometrytypes = append(geometrytypes, "MultiPolygon")
					}
				case *geom.LineString:
					if !slices.Contains(geometrytypes, "LineString") {
						geometrytypes = append(geometrytypes, "LineString")
					}
				case *geom.MultiLineString:
					if !slices.Contains(geometrytypes, "MultiLineString") {
						geometrytypes = append(geometrytypes, "MultiLineString")
					}
				}
				wkbbytes, err = wkb.Marshal(gt, wkb.NDR)
				if err != nil {
					fmt.Fprintf(out, "\r\n\r\ngemometry marshal error: %s", err)
					return
				}
			} else {
				point := geom.NewPoint(geom.XY).MustSetCoords([]float64{0, 0})
				wkbbytes, err = wkb.Marshal(point, wkb.NDR)
				if err != nil {
					fmt.Fprintf(out, "\r\n\r\ngeometry marshal error: %s", err)
					return
				}
				if !slices.Contains(geometrytypes, "Point") {
					geometrytypes = append(geometrytypes, "Point")
				}
			}
			item := InventoryItem{
				FlatInventory: h,
				Geometry:      wkbbytes,
			}
			item.Geometry = wkbbytes
			if user == nil {
				item.Issue = make([]string, 0)
				item.CanPublish = nil
			}
			itemsbatch = append(itemsbatch, item)
		}

		_, err = writer.Write(itemsbatch)
		if err != nil {
			fmt.Fprintf(out, "\r\n\r\nwrite error: %s", err)
			return
		}
		if err = writer.Flush(); err != nil {
			fmt.Fprintf(out, "\r\n\r\nflush error: %s", err)
			return
		}
		params.From = (result.From + len(result.Hits))

		if int64(params.From) >= result.Total {
			meta := struct {
				Version       string `json:"version"`
				PrimaryColumn string `json:"primary_column"`
				Columns       map[string]struct {
					Encoding      string   `json:"encoding"`
					GeometryTypes []string `json:"geometry_types"`
				} `json:"columns"`
			}{
				Version:       "1.1.0",
				PrimaryColumn: "geometry",
				Columns: map[string]struct {
					Encoding      string   `json:"encoding"`
					GeometryTypes []string `json:"geometry_types"`
				}{
					"geometry": {
						Encoding:      "WKB",
						GeometryTypes: geometrytypes,
					},
				},
			}
			metajson, err := json.Marshal(meta)
			if err != nil {
				fmt.Fprintf(out, "\r\n\r\nheader creation error: %s", err)
				return
			}
			writer.SetKeyValueMetadata("geo", string(metajson))
			return
		}
	}
}
