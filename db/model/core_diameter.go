package model

type CoreDiameter struct {
	ID           int32   `json:"core_diameter_id,omitempty"`
	Name         string  `json:"core_diameter_name,omitempty"`
	CoreDiameter float64 `json:"core_diameter,omitempty"`
	Unit         string  `json:"core_diameter_unit,omitempty"`
}
