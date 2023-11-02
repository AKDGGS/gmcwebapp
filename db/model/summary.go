package model

type Summary struct {
	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	Containers struct {
		PathCache      *string `json:"path_cache,omitempty"`
		ContainerTotal *int32  `json:"container_total,omitempty"`
	} `json:"containers,omitempty"`

	Collections []struct {
		Collection      string `json:"collection,omitempty"`
		CollectionTotal *int32 `json:"collection_total,omitempty"`
	} `json:"collections,omitempty"`

	Keywords struct {
		Keywords []string `json:"keywords,omitempty"`
	} `json:"keywords,omitempty"`

	Barcodes struct {
		Barcodes     []string `json:"barcodes,omitempty"`
		BarcodeCount int32    `json:"barcode_count,omitempty"`
	} `json:"barcodes,omitempty"`

	Boreholes []struct {
		Borehole *string `json:"borehole,omitempty"`
		Prospect *string `json:"prospect,omitempty"`
		Total    *int32  `json:"total,omitempty"`
	} `json:"boreholes,omitempty"`

	Outcrops []struct {
		Outcrop *string `json:"outcrop,omitempty"`
		Total   *int32  `json:"total,omitempty"`
	} `json:"outcrops,omitempty"`

	Shotlines []struct {
		Shotline *string `json:"shotline,omitempty"`
		Total    *int32  `json:"total,omitempty"`
	} `json:"shotlines,omitempty"`

	Wells []struct {
		Well  *string `json:"well,omitempty"`
		Total *int32  `json:"total,omitempty"`
	} `json:"wells,omitempty"`
}
