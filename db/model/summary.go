package model

type Summary struct {
	Barcodes     []string `json:"barcodes,omitempty"`
	BarcodeCount int32    `json:"barcode_count,omitempty"`
	Keywords     []string `json:"keywords,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	ParentContainers []string `json:"parent_containers,omitempty"`
	ChildContainers  []string `json:"child_containers,omitempty"`
	Collections      []string `json:"collections,omitempty"`

	Boreholes []struct {
		Borehole *string `json:"borehole,omitempty"`
		Prospect *string `json:"prospect,omitempty"`
	} `json:"boreholes,omitempty"`

	Outcrops []struct {
		Outcrop *string `json:"outcrop,omitempty"`
	} `json:"outcrops,omitempty"`

	Shotlines []struct {
		Shotline *string `json:"shotline,omitempty"`
	} `json:"shotlines,omitempty"`

	Wells []struct {
		Well *string `json:"well,omitempty"`
	} `json:"wells,omitempty"`
}
