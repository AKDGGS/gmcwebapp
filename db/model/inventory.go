package model

import (
	"encoding/json"
	"time"
)

type Inventory struct {
	ID                    int32  `json:"id"`
	ContainerPathCache    string `json:"container_path_cache"`
	Barcode               string `json:"barcode"`
	AltBarcode            string `json:"alt_barcode"`
	CollectionID          int32  `json:"collection_id"`
	CollectionName        string `json:"collection_name"`
	CollectionDescription string `json:"collection_description"`

	ContainerID         int32  `json:"container_id"`
	ContainerName       string `json:"container_name"`
	ContainerRemark     string `json:"container_remark"`
	ContainerBarcode    string `json:"container_barcode"`
	ContainerAltBarcode string `json:"container_alt_barcode"`

	SampleID                 int64    `json:"dggs_sample_id"`
	SampleNumber             string   `json:"sample_number"`
	SampleNumberPrefix       string   `json:"sample_number_prefix"`
	AltSampleNumber          string   `json:"alt_sample_number"`
	PublishedSampleNumber    string   `json:"published_sample_number"`
	PublishedNumberHasSuffix bool     `json:"published_number_has_suffix"`
	StateNumber              string   `json:"state_number"`
	BoxNumber                string   `json:"box_number"`
	SetNumber                string   `json:"set_number"`
	SplitNumber              string   `json:"split_number"`
	SlideNumber              string   `json:"slide_number"`
	SlipNumber               int32    `json:"slip_number"`
	LabNumber                string   `json:"lab_number"`
	LabReportID              string   `json:"lab_report_id"`
	MapNumber                string   `json:"map_number"`
	Description              string   `json:"description"`
	Remark                   string   `json:"remark"`
	Tray                     int16    `json:"tray"`
	IntervalTop              float64  `json:"interval_top"`
	IntervalBottom           float64  `json:"interval_bottom"`
	Keywords                 []string `json:"keywords"`
	IntervalUnit             string   `json:"interval_unit"`
	CoreNumber               string   `json:"core_number"`

	CoreDiameterId   int32   `json:"core_diameter_id"`
	CoreDiameterName string  `json:"core_diameter_name"`
	CoreDiameter     float64 `json:"core_diameter"`
	CoreDiameterUnit string  `json:"core_diameter_unit"`

	Weight          float64                `json:"weight"`
	WeightUnit      string                 `json:"weight_unit"`
	SampleFrequency string                 `json:"sample_frequency"`
	Recovery        string                 `json:"recovery"`
	CanPublish      bool                   `json:"can_publish"`
	RadiationMSVH   float32                `json:"radiation_msvh"`
	ReceivedDate    *time.Time             `json:"received_date"`
	EnteredDate     *time.Time             `json:"entered_date"`
	ModifiedDate    *time.Time             `json:"modified_date"`
	ModifiedUser    string                 `json:"modified_user"`
	Active          bool                   `json:"active"`
	Stash           map[string]interface{} `json:"stash"`
	GeoJSON         map[string]interface{} `json:"geojson"`
	Boreholes       []Borehole             `json:"boreholes"`
	Outcrops        []Outcrop              `json:"outcrops"`
	Shotlines       []Shotline             `json:"shotlines"`
	Wells           []Well                 `json:"wells"`
	Organizations   []Organization         `json:"organizations"`
	Notes           []Note                 `json:"notes"`
	URLs            []URL                  `json:"urls"`
	Files           []File                 `json:"files"`
	Publications    []Publication          `json:"publications"`
	ContainerLog    []ContainerLog         `json:"container_log"`
	Qualities       []Qualities            `json:"qualities"`
}

func (i *Inventory) MarshalJSON() ([]byte, error) {
	type Alias Inventory
	return json.Marshal(&struct {
		ReceivedDate string `json:"received_date"`
		EnteredDate  string `json:"entered_date"`
		ModifiedDate string `json:"modified_date"`
		*Alias
	}{
		ReceivedDate: i.ReceivedDate.Format("01-02-2006"),
		EnteredDate:  i.EnteredDate.Format("01-02-2006"),
		ModifiedDate: i.ModifiedDate.Format("01-02-2006"),
		Alias:        (*Alias)(i),
	})
}
