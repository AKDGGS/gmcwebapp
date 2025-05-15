package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	sutil "gmc/search/util"

	"codeberg.org/go-pdf/fpdf"
)

func (srv *Server) ServeSearchInventoryPDF(w http.ResponseWriter, r *http.Request) {
	user, err := srv.Auths.CheckRequest(w, r)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("authentication error: %s", err),
			http.StatusBadRequest,
		)
		return
	}

	params := &sutil.InventoryParams{}
	params.ParseQuery(r.URL.Query(), (user != nil))
	params.From = 0

	if srv.Config.ExportLimit > 0 && srv.Config.ExportLimit < 10000 {
		params.Size = srv.Config.ExportLimit
	} else {
		params.Size = 10000
	}

	tmpl, err := template.ParseFiles("assets/tmpl/inventory_item.txt")
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("pdf creation error: %s", err),
			http.StatusInternalServerError,
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
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetLineWidth(0.5)

	headerheight := 8.0

	header := []string{"Related", "Sample /\nSlide", "Box /\nSet", "Core No /\nDiam", "Top /\nBottom", "Keywords", "Collection"}
	columnwidth := []float64{55.0, 20.0, 20.0, 15.0, 20.0, 35.0, 25.0}
	if user != nil {
		header = append(header, "Barcode", "Location")
		columnwidth = []float64{35.0, 20.0, 15.0, 15.0, 15.0, 22.0, 18.0, 25.0, 25.0}
	}

	nextrecord := []string{}
	nextheight := 0.0
	for {
		if srv.Config.ExportLimit > 0 && params.From+params.Size > srv.Config.ExportLimit {
			params.Size = srv.Config.ExportLimit - params.From
		}
		result, err := srv.Search.SearchInventory(params)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("search error: %s", err),
				http.StatusBadRequest,
			)
			return
		}
		if srv.Config.ExportLimit > 0 && result.Total > int64(srv.Config.ExportLimit) {
			result.Total = int64(srv.Config.ExportLimit)
		}
		i := 0
		for i < len(result.Hits) {

			if pdf.GetY() < 11 {
				if pdf.PageNo() == 1 {
					pdf.SetFont("Arial", "B", 12)
					pdf.CellFormat(50, 5, "SEARCH RESULTS", "0", 0, "LM", false, 0, "")
					pdf.SetFont("Arial", "B", 8)
					pdf.Ln(-1)
					pdf.SetY(pdf.GetY() + 2)
					pdf.SetFont("Arial", "B", 10)
					pdf.MultiCell(20, 3, "Query: ", "0", "LM", false)
					pdf.SetFont("Arial", "", 8)
					pdf.SetTextColor(0, 0, 255)
					pdf.WriteLinkString(5, r.Host+r.URL.String(), r.Host+r.URL.String())
					pdf.SetTextColor(0, 0, 0)
					pdf.SetFont("Arial", "B", 10)
					pdf.Ln(-1)
					pdf.SetY(pdf.GetY() + 2)
				}
				pdf.SetFillColor(256, 256, 256)
				//draw header
				start_x := pdf.GetX()
				start_y := pdf.GetY()
				pdf.SetFont("Arial", "B", 8)
				for j := 0; j < len(header); j++ {
					pdf.SetXY(start_x, start_y)
					pdf.Rect(start_x, start_y, columnwidth[j], headerheight, "F")
					lines := pdf.SplitLines([]byte(header[j]), columnwidth[j])
					if len(lines) == 1 {
						pdf.SetXY(start_x, start_y+(headerheight/2))
					}
					pdf.MultiCell(columnwidth[j], headerheight/2, header[j], "0", "B", false)
					if len(lines) == 1 {
						pdf.SetXY(start_x, start_y-(headerheight/2))
					}
					start_x += columnwidth[j]
				}
				pdf.SetY(start_y + headerheight - 3)
				pdf.Ln(-1)
				pageWidth, _ := pdf.GetPageSize()
				pdf.Line(10, pdf.GetY(), pageWidth-10, pdf.GetY())
				pdf.SetY(pdf.GetY() + 0.3)
				pdf.SetFont("Arial", "", 8)
			}

			if len(nextrecord) == 0 {
				h := result.Hits[i]
				var buf bytes.Buffer
				if err = tmpl.Execute(&buf, h); err != nil {
					http.Error(
						w,
						fmt.Sprintf("data formatting error: %s", err),
						http.StatusInternalServerError,
					)
					return
				}

				for j := 0; j < 7; j++ {
					nextrecord = append(nextrecord, "")
				}
				nextrecord[0] = buf.String()
				nextrecord[1] = ""
				if h.SampleNumber != nil {
					nextrecord[1] = *h.SampleNumber
				}
				if h.SlideNumber != nil {
					nextrecord[1] = "\n" + nextrecord[1] + *h.SlideNumber
				}
				nextrecord[2] = ""
				if h.BoxNumber != nil {
					nextrecord[2] = *h.BoxNumber
				}
				if h.SetNumber != nil {
					nextrecord[2] += "\n" + *h.SetNumber
				}
				nextrecord[3] = ""
				if h.CoreNumber != nil {
					nextrecord[3] = *h.CoreNumber
				}
				if h.CoreDiameter != nil {
					nextrecord[3] += "\n" + strconv.FormatFloat(*h.CoreDiameter, 'f', 2, 64)
				}
				nextrecord[4] = ""
				if h.IntervalTop != nil {
					nextrecord[4] = strconv.FormatFloat(*h.IntervalTop, 'f', 2, 64)
					if h.IntervalUnit != nil {
						nextrecord[4] += " " + *h.IntervalUnit
					}
					nextrecord[4] += "\n"
				}
				if h.IntervalBottom != nil {
					nextrecord[4] += strconv.FormatFloat(*h.IntervalBottom, 'f', 2, 64)
					if h.IntervalUnit != nil {
						nextrecord[4] += " " + *h.IntervalUnit
					}
				}
				nextrecord[5] = strings.Join(h.Keyword, "; ")
				if h.Collection != nil {
					nextrecord[6] = *h.Collection
				}

				if user != nil {
					nextrecord = append(nextrecord, "", "") //barcode, location
					if h.DisplayBarcode != nil {
						nextrecord[7] = *h.DisplayBarcode
					}
					if h.ContainerPath != nil {
						nextrecord[8] = *h.ContainerPath
					}
				}

				_, lineheight := pdf.GetFontSize()

				for j := 0; j < len(nextrecord); j++ {
					lines := pdf.SplitLines([]byte(nextrecord[j]), columnwidth[j])
					currentheight := float64(len(lines))*lineheight + float64(len(lines))
					if nextheight < currentheight {
						nextheight = currentheight
					}
				}
			}

			if pdf.GetY()+nextheight < 273 {
				start_x := pdf.GetX()
				start_y := pdf.GetY()
				r, _, _ := pdf.GetFillColor()
				if r == 255 {
					pdf.SetFillColor(220, 220, 220)
				} else {
					pdf.SetFillColor(256, 256, 256)
				}
				for i := 0; i < len(nextrecord); i++ {
					pdf.SetXY(start_x, start_y)
					pdf.Rect(start_x, start_y, columnwidth[i], nextheight, "F")
					pdf.MultiCell(columnwidth[i], pdf.PointConvert(10), nextrecord[i], "0", "LM", false)
					start_x += columnwidth[i]
				}
				pdf.SetY(start_y + nextheight)
				nextrecord = []string{}
				nextheight = 0.0
				i++
			}

			if pdf.GetY()+nextheight >= 273 || int64(params.From+i) >= result.Total {
				// draw footer
				pdf.SetY(-24)
				pageWidth, _ := pdf.GetPageSize()
				pdf.Line(10, pdf.GetY(), pageWidth-10, pdf.GetY())
				pdf.Ln(-1)
				pdf.SetFont("Arial", "I", 8)
				pdf.CellFormat(0, 0, fmt.Sprintf(time.Now().Format("January 2, 2006, 15:04:05PM -0700 AKT")+", Page %d", pdf.PageNo()), "", 0, "RM", false, 0, "")
			}
			if pdf.GetY()+nextheight >= 273 && int64(params.From+i) < result.Total {
				pdf.AddPage()
			}
		}
		params.From = (result.From + len(result.Hits))
		if int64(params.From) >= result.Total {
			if err = pdf.Output(out); err != nil {
				http.Error(
					w,
					fmt.Sprintf("pdf creation error: %s", err),
					http.StatusInternalServerError,
				)
				return
			}
			w.Header().Set("Content-Type", "application/pdf")
			return
		}
	}
}
