package model

import (
	"encoding/json"
	"time"
)

type Inventory struct {
	ID           int32      `json:"id"`
	Barcode      string     `json:"barcode,omitempty"`
	AltBarcode   string     `json:"alt_barcode,omitempty"`
	CollectionID int32      `json:"collection_id,omitempty"`
	Collection   Collection `json:"collection,omitempty"`

	ContainerID int32     `json:"container_id,omitempty"`
	Container   Container `json:"container,omitempty"`

	SampleID                 int64    `json:"dggs_sample_id,omitempty"`
	SampleNumber             string   `json:"sample_number,omitempty"`
	SampleNumberPrefix       string   `json:"sample_number_prefix,omitempty"`
	AltSampleNumber          string   `json:"alt_sample_number,omitempty"`
	PublishedSampleNumber    string   `json:"published_sample_number,omitempty"`
	PublishedNumberHasSuffix bool     `json:"published_number_has_suffix"`
	StateNumber              string   `json:"state_number,omitempty"`
	BoxNumber                string   `json:"box_number,omitempty"`
	SetNumber                string   `json:"set_number,omitempty"`
	SplitNumber              string   `json:"split_number,omitempty"`
	SlideNumber              string   `json:"slide_number,omitempty"`
	SlipNumber               int32    `json:"slip_number,omitempty"`
	LabNumber                string   `json:"lab_number,omitempty"`
	LabReportID              string   `json:"lab_report_id,omitempty"`
	MapNumber                string   `json:"map_number,omitempty"`
	Description              string   `json:"description,omitempty"`
	Remark                   string   `json:"remark,omitempty"`
	Tray                     int16    `json:"tray,omitempty"`
	IntervalTop              float64  `json:"interval_top"`
	IntervalBottom           float64  `json:"interval_bottom,omitempty"`
	Keywords                 []string `json:"keywords,omitempty"`
	IntervalUnit             string   `json:"interval_unit,omitempty"`
	CoreNumber               string   `json:"core_number,omitempty"`

	CoreDiameterID int32        `json:"core_diameter_id,omitempty"`
	CoreDiameter   CoreDiameter `json:"core_diameter,omitempty"`

	Weight          float64                `json:"weight,omitempty"`
	WeightUnit      string                 `json:"weight_unit,omitempty"`
	SampleFrequency string                 `json:"sample_frequency,omitempty"`
	Recovery        string                 `json:"recovery,omitempty"`
	CanPublish      bool                   `json:"can_publish"`
	RadiationMSVH   float32                `json:"radiation_msvh,omitempty"`
	ReceivedDate    *time.Time             `json:"received_date,omitempty"`
	EnteredDate     *time.Time             `json:"entered_date,omitempty"`
	ModifiedDate    *time.Time             `json:"modified_date,omitempty"`
	ModifiedUser    string                 `json:"modified_user,omitempty"`
	Active          bool                   `json:"active"`
	Stash           map[string]interface{} `json:"stash,omitempty"`
	GeoJSON         map[string]interface{} `json:"geojson,omitempty"`
	Boreholes       []Borehole             `json:"boreholes,omitempty"`
	Outcrops        []Outcrop              `json:"outcrops,omitempty"`
	Shotpoints      []Shotpoint            `json:"shotpoints,omitempty"`
	Wells           []Well                 `json:"wells,omitempty"`
	Organizations   []Organization         `json:"organizations,omitempty"`
	Notes           []Note                 `json:"notes,omitempty"`
	URLs            []URL                  `json:"urls,omitempty"`
	Files           []File                 `json:"files,omitempty"`
	Publications    []Publication          `json:"publications,omitempty"`
	ContainerLog    []ContainerLog         `json:"container_log,omitempty"`
	Qualities       []Qualities            `json:"qualities,omitempty"`
}

func (i *Inventory) MarshalJSON() ([]byte, error) {
	type Alias Inventory
	receivedDate := ""
	if i.ReceivedDate != nil {
		receivedDate = i.ReceivedDate.Format("01-02-2006")
	}
	enteredDate := ""
	if i.EnteredDate != nil {
		enteredDate = i.EnteredDate.Format("01-02-2006")
	}
	modifiedDate := ""
	if i.ModifiedDate != nil {
		modifiedDate = i.ModifiedDate.Format("01-02-2006")
	}
	return json.Marshal(&struct {
		ReceivedDate string `json:"received_date,omitempty"`
		EnteredDate  string `json:"entered_date,omitempty"`
		ModifiedDate string `json:"modified_date,omitempty"`
		*Alias
	}{
		ReceivedDate: receivedDate,
		EnteredDate:  enteredDate,
		ModifiedDate: modifiedDate,
		Alias:        (*Alias)(i),
	})
}
