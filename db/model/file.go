package model

import "fmt"

type File struct {
	ID          int32   `db:"id" json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Type        string  `db:"mimetype" json:"mimetype"`
	Size        int64   `json:"size"`

	BoreholeIDs  []int    `json:"borehole_ids,omitempty"`
	InventoryIDs []int    `json:"inventory_ids,omitempty"`
	OutcropIDs   []int    `json:"outcrop_ids,omitempty"`
	ProspectIDs  []int    `json:"prospect_ids,omitempty"`
	WellIDs      []int    `json:"well_ids,omitempty"`
	Barcodes     []string `json:"barcodes,omitempty"`
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
