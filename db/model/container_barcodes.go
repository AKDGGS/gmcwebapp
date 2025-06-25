package model

type ContainerBarcodes struct {
	Barcode    *string `json:"barcode,omitempty"`
	AltBarcode *string `db:"alt_barcode" json:"alt_barcode,omitempty"`
}
