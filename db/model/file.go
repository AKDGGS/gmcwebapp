package model

import "fmt"

type File struct {
	ID          int    `json:"file_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int64  `json:"size"`
	Type        string `json:"mimetype"`
	MD5         string `json:"content_md5"`
	BoreholeIDs []int  `json:"borehole_ids"`
	WellIDs     []int  `json:"well_ids"`
}

func (f *File) FormattedSize() string {
	if f.Size < 1024 {
		return fmt.Sprintf("%d B", f.Size)
	} else if f.Size < 1048576 {
		return fmt.Sprintf("%.1f KB", float64(f.Size)/1024)
	} else {
		return fmt.Sprintf("%.1f MB", float64(f.Size)/(1048576))
	}
}
