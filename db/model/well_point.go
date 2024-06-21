package model

type WellPoint struct {
	WellID int32   `json:"well_id"`
	Name   *string `json:"name"`
	Geog   *string `json:"geog"`
}
