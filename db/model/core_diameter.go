package model

type CoreDiameter struct {
	ID           int32   `json:"id,omitempty"`
	Name         *string `json:"name,omitempty"`
	CoreDiameter float64 `json:"core_diameter"`
	Unit         *string `json:"unit,omitempty"`
}
