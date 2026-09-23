package model

import (
	"fmt"
	"math"
)

type File struct {
	ID            int32   `db:"id" json:"id,omitempty"`
	Name          string  `json:"name,omitempty"`
	Description   *string `json:"description,omitempty"`
	MimeType      string  `db:"mimetype" json:"mimetype"`
	FileType      string  `db:"file_type" json:"file_type,omitempty"`
	FileTypeColor *int32  `db:"file_type_int"`
	Size          int64   `json:"size"`
	Hash          *string `db:"file_hash" json:"file_hash,omitempty"`

	// Transient fields that don't directly exist in the file table
	Barcode      string `json:"barcode,omitempty"`
	BoreholeIDs  []int  `json:"borehole_ids,omitempty"`
	InventoryIDs []int  `json:"inventory_ids,omitempty"`
	OutcropIDs   []int  `json:"outcrop_ids,omitempty"`
	ProspectIDs  []int  `json:"prospect_ids,omitempty"`
	WellIDs      []int  `json:"well_ids,omitempty"`
}

func (f *File) FormattedSize() string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	if f.Size == 0 {
		return "0 B"
	}

	idx := min(int(math.Log2(float64(f.Size))/10), len(units)-1)
	value := float64(f.Size) / math.Pow(1024, float64(idx))
	return fmt.Sprintf("%.2f %s", value, units[idx])
}
