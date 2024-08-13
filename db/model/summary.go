package model

type Summary struct {
	Barcodes     []string `json:"barcodes,omitempty"`
	BarcodeTotal int32    `json:"barcode_total,omitempty"`
	Container    *string  `json:"path_cache,omitempty"`
	Total        *int32   `json:"container_total,omitempty"`
	Keywords     []string `json:"keywords,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database

	Collections []struct {
		Collection      string `json:"collection,omitempty"`
		CollectionTotal *int32 `json:"collection_total,omitempty"`
	} `json:"collections,omitempty"`

	Boreholes []struct {
		Borehole      *string `json:"borehole,omitempty"`
		Prospect      *string `json:"prospect,omitempty"`
		BoreholeTotal *int32  `json:"borehole_total,omitempty"`
	} `json:"boreholes,omitempty"`

	Outcrops []struct {
		Outcrop      *string `json:"outcrop,omitempty"`
		OutcropTotal *int32  `json:"outcrop_total,omitempty"`
	} `json:"outcrops,omitempty"`

	Shotlines []struct {
		Shotline      *string `json:"shotline,omitempty"`
		ShotlineTotal *int32  `json:"shotline_total,omitempty"`
	} `json:"shotlines,omitempty"`

	Wells []struct {
		Well      *string `json:"well,omitempty"`
		WellTotal *int32  `json:"well_total,omitempty"`
	} `json:"wells,omitempty"`
}
