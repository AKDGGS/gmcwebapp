package web

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	sutil "gmc/search/util"
)

func (srv *Server) ServeSearchInventoryCSV(w http.ResponseWriter, r *http.Request) {
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

	cwr := csv.NewWriter(out)
	defer cwr.Flush()

	w.Header().Set("Content-Type", "text/csv")

	cwr.Write([]string{
		"ID",
		"Collection ID",
		"Collection",
		"Sample Number",
		"Slide Number",
		"Box Number",
		"Set Number",
		"Core Number",
		"Core Diameter",
		"Core Name",
		"Core Unit",
		"Interval Top",
		"Interval Bottom",
		"Interval Unit",
		"Keywords",
		"Barcode",
		"Container ID",
		"Container",
		"Project ID",
		"Project",
	})

	params := &sutil.InventoryParams{}
	params.ParseQuery(r.URL.Query(), (user != nil))
	params.From = 0
	params.Size = 1000

	for {
		result, err := srv.Search.SearchInventory(params)
		if err != nil {
			fmt.Fprintf(out, "\r\n\r\nsearch error: %s", err)
			return
		}

		for _, h := range result.Hits {
			cwr.Write([]string{
				qfmt(h.ID),
				qfmt(h.CollectionID),
				qfmt(h.Collection),
				qfmt(h.SampleNumber),
				qfmt(h.SlideNumber),
				qfmt(h.BoxNumber),
				qfmt(h.SetNumber),
				qfmt(h.CoreNumber),
				qfmt(h.CoreDiameter),
				qfmt(h.CoreName),
				qfmt(h.CoreUnit),
				qfmt(h.IntervalTop),
				qfmt(h.IntervalBottom),
				qfmt(h.IntervalUnit),
				qfmt(h.Keyword),
				qfmt(h.DisplayBarcode),
				qfmt(h.ContainerID),
				qfmt(h.ContainerPath),
				qfmt(h.ProjectID),
				qfmt(h.Project),
			})
		}

		params.From = (result.From + len(result.Hits))
		if int64(params.From) >= result.Total {
			return
		}
	}
}

func qfmt(v interface{}) string {
	switch t := v.(type) {
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case *int32:
		if t == nil {
			return ""
		}
		return strconv.FormatInt(int64(*t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case *int64:
		if t == nil {
			return ""
		}
		return strconv.FormatInt(*t, 10)
	case float64:
		return strconv.FormatFloat(t, 'f', 2, 64)
	case *float64:
		if t == nil {
			return ""
		}
		return strconv.FormatFloat(*t, 'f', 2, 64)
	case float32:
		return strconv.FormatFloat(float64(t), 'f', 2, 32)
	case *float32:
		if t == nil {
			return ""
		}
		return strconv.FormatFloat(float64(*t), 'f', 2, 32)
	case string:
		return t
	case *string:
		if t == nil {
			return ""
		}
		return *t
	case []string:
		return strings.Join(t, "; ")
	default:
		return ""
	}
}
