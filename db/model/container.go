package model

type Container struct {
	ID         int32   `db:"id" json:"id,omitempty"`
	Name       string  `json:"name,omitempty"`
	PathCache  *string `db:"path_cache" json:"path_cache,omitempty"`
	Remark     *string `json:"remark,omitempty"`
	Barcode    *string `json:"barcode,omitempty"`
	AltBarcode *string `db:"alt_barcode" json:"alt_barcode,omitempty"`
}
