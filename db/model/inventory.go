package model

import (
	"encoding/json"
	"time"
)

type Inventory struct {
	ID         int32       `json:"id"`
	Barcode    *string     `json:"barcode,omitempty"`
	AltBarcode *string     `json:"alt_barcode,omitempty"`
	Collection *Collection `json:"collection,omitempty"`
	Container  *Container  `json:"container,omitempty"`

	SampleID                 *int64   `json:"dggs_sample_id,omitempty"`
	SampleNumber             *string  `json:"sample_number,omitempty"`
	SampleNumberPrefix       *string  `json:"sample_number_prefix,omitempty"`
	AltSampleNumber          *string  `json:"alt_sample_number,omitempty"`
	PublishedSampleNumber    *string  `json:"published_sample_number,omitempty"`
	PublishedNumberHasSuffix bool     `json:"published_number_has_suffix"`
	StateNumber              *string  `json:"state_number,omitempty"`
	BoxNumber                *string  `json:"box_number,omitempty"`
	SetNumber                *string  `json:"set_number,omitempty"`
	SplitNumber              *string  `json:"split_number,omitempty"`
	SlideNumber              *string  `json:"slide_number,omitempty"`
	SlipNumber               *int32   `json:"slip_number,omitempty"`
	LabNumber                *string  `json:"lab_number,omitempty"`
	LabReportID              *string  `json:"lab_report_id,omitempty"`
	MapNumber                *string  `json:"map_number,omitempty"`
	Description              *string  `json:"description,omitempty"`
	Remark                   *string  `json:"remark,omitempty"`
	Tray                     *int16   `json:"tray,omitempty"`
	IntervalTop              *float64 `json:"interval_top"`
	IntervalBottom           *float64 `json:"interval_bottom,omitempty"`
	IntervalUnit             *string  `json:"interval_unit,omitempty"`
	Keywords                 []string `json:"keywords,omitempty"`
	CoreNumber               *string  `json:"core_number,omitempty"`

	CoreDiameter *CoreDiameter `json:"core_diameter,omitempty"`

	Weight          *float64               `json:"weight,omitempty"`
	WeightUnit      *string                `json:"weight_unit,omitempty"`
	SampleFrequency *string                `json:"sample_frequency,omitempty"`
	Recovery        *string                `json:"recovery,omitempty"`
	CanPublish      bool                   `json:"can_publish"`
	RadiationMSVH   *float64               `json:"radiation_msvh,omitempty"`
	ReceivedDate    *time.Time             `json:"received_date,omitempty"`
	EnteredDate     *time.Time             `json:"entered_date,omitempty"`
	ModifiedDate    *time.Time             `json:"modified_date,omitempty"`
	ModifiedUser    *string                `json:"modified_user,omitempty"`
	Active          bool                   `json:"active"`
	Stash           map[string]interface{} `json:"stash,omitempty"`

	Boreholes     []Borehole     `json:"boreholes,omitempty"`
	Outcrops      []Outcrop      `json:"outcrops,omitempty"`
	Shotpoints    []Shotpoint    `json:"shotpoints,omitempty"`
	Wells         []Well         `json:"wells,omitempty"`
	Organizations []Organization `json:"organizations,omitempty"`
	Notes         []Note         `json:"notes,omitempty"`
	URLs          []URL          `json:"urls,omitempty"`
	Files         []File         `json:"files,omitempty"`
	Publications  []Publication  `json:"publications,omitempty"`
	ContainerLog  []ContainerLog `json:"container_logs,omitempty"`
	Qualities     []Quality      `json:"qualities,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	GeoJSON interface{} `json:"geojson,omitempty"`
}

func (i *Inventory) MarshalJSON() ([]byte, error) {
	type Alias Inventory
	var receivedDate string
	if !i.ReceivedDate.IsZero() {
		receivedDate = i.ReceivedDate.Format("01-02-2006")
	}
	var enteredDate string
	if !i.EnteredDate.IsZero() {
		enteredDate = i.EnteredDate.Format("01-02-2006")
	}
	var modifiedDate string
	if !i.ModifiedDate.IsZero() {
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

func (i *Inventory) Shotlines() []*Shotline {
	var shotlines []*Shotline
	sp_map := make(map[int32][]Shotpoint)
	for _, sp := range i.Shotpoints {
		shotline := sp.Shotline
		if _, ok := sp_map[shotline.ID]; !ok {
			shotlines = append(shotlines, shotline)
			sp_map[shotline.ID] = []Shotpoint{}
		}
		sp_map[shotline.ID] = append(sp_map[shotline.ID], sp)
	}
	for _, sl := range shotlines {
		sl.Shotpoints = sp_map[sl.ID]
	}
	return shotlines
}
