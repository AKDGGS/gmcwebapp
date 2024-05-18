package model

import (
	"encoding/json"
	"time"
)

type Inventory struct {
	ID                       int32       `json:"id"`
	Barcode                  *string     `json:"barcode,omitempty"`
	AltBarcode               *string     `db:"alt_barcode" json:"alt_barcode,omitempty"`
	SampleID                 *int64      `db:"sample_id" json:"sample_id,omitempty"`
	SampleNumber             *string     `db:"sample_number" json:"sample_number,omitempty"`
	SampleNumberPrefix       *string     `db:"sample_number_prefix" json:"sample_number_prefix,omitempty"`
	AltSampleNumber          *string     `db:"alt_sample_number" json:"alt_sample_number,omitempty"`
	PublishedSampleNumber    *string     `db:"published_sample_number" json:"published_sample_number,omitempty"`
	PublishedNumberHasSuffix bool        `db:"published_number_has_suffix" json:"published_number_has_suffix"`
	StateNumber              *string     `db:"state_number" json:"state_number,omitempty"`
	BoxNumber                *string     `db:"box_number" json:"box_number,omitempty"`
	SetNumber                *string     `db:"set_number" json:"set_number,omitempty"`
	SplitNumber              *string     `db:"split_number" json:"split_number,omitempty"`
	SlideNumber              *string     `db:"slide_number" json:"slide_number,omitempty"`
	SlipNumber               *int32      `db:"slip_number" json:"slip_number,omitempty"`
	LabNumber                *string     `db:"lab_number" json:"lab_number,omitempty"`
	LabReportID              *string     `db:"lab_report_id" json:"lab_report_id,omitempty"`
	MapNumber                *string     `db:"map_number" json:"map_number,omitempty"`
	Description              *string     `json:"description,omitempty"`
	Remark                   *string     `json:"remark,omitempty"`
	Tray                     *int16      `json:"tray,omitempty"`
	IntervalTop              *float64    `db:"interval_top" json:"interval_top"`
	IntervalBottom           *float64    `db:"interval_bottom" json:"interval_bottom,omitempty"`
	IntervalUnit             *string     `db:"interval_unit" json:"interval_unit,omitempty"`
	Keywords                 []string    `json:"keywords,omitempty"`
	CoreNumber               *string     `db:"core_number" json:"core_number,omitempty"`
	Weight                   *float64    `json:"weight,omitempty"`
	WeightUnit               *string     `db:"weight_unit" json:"weight_unit,omitempty"`
	SampleFrequency          *string     `db:"sample_frequency" json:"sample_frequency,omitempty"`
	Recovery                 *string     `json:"recovery,omitempty"`
	CanPublish               bool        `db:"can_publish" json:"can_publish"`
	RadiationMSVH            *float64    `db:"radiation_msvh" json:"radiation_msvh,omitempty"`
	ReceivedDate             *time.Time  `db:"received_date" json:"received_date,omitempty"`
	EnteredDate              *time.Time  `db:"entered_date" json:"entered_date,omitempty"`
	ModifiedDate             *time.Time  `db:"modified_date" json:"modified_date,omitempty"`
	ModifiedUser             *string     `db:"modified_user" json:"modified_user,omitempty"`
	Active                   bool        `json:"active"`
	Stash                    interface{} `json:"stash,omitempty"`

	CoreDiameter  *CoreDiameter  `db:"core_diameter" json:"core_diameter,omitempty"`
	Collection    *Collection    `json:"collection,omitempty"`
	Container     *Container     `json:"container,omitempty"`
	Boreholes     []Borehole     `json:"boreholes,omitempty"`
	Outcrops      []Outcrop      `json:"outcrops,omitempty"`
	Shotpoints    []Shotpoint    `json:"shotpoints,omitempty"`
	Wells         []Well         `json:"wells,omitempty"`
	Organizations []Organization `json:"organizations,omitempty"`
	Notes         []Note         `json:"notes,omitempty"`
	URLs          []URL          `json:"urls,omitempty"`
	Files         []File         `json:"files,omitempty"`
	Publications  []Publication  `json:"publications,omitempty"`
	ContainerLog  []ContainerLog `db:"container_logs" json:"container_logs,omitempty"`
	Qualities     []Quality      `json:"qualities,omitempty"`

	//transient fields that are generated on-the-fly
	//these fields don't exist in the database
	GeoJSON interface{} `json:"geojson,omitempty"`
}

func (i *Inventory) MarshalJSON() ([]byte, error) {
	type Alias Inventory
	var receivedDate string
	if i.ReceivedDate != nil && !i.ReceivedDate.IsZero() {
		receivedDate = i.ReceivedDate.Format("01-02-2006")
	}
	var enteredDate string
	if i.EnteredDate != nil && !i.EnteredDate.IsZero() {
		enteredDate = i.EnteredDate.Format("01-02-2006")
	}
	var modifiedDate string
	if i.ModifiedDate != nil && !i.ModifiedDate.IsZero() {
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
