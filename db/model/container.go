package model

type Container struct {
	ID         int32  `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	PathCache  string `json:"path_cache,omitempty"`
	Remark     string `json:"remark,omitempty"`
	Barcode    string `json:"barcode,omitempty"`
	AltBarcode string `json:"alt_barcode,omitempty"`
}
