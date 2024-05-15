package model

type CoreDiameter struct {
	ID           int32   `db:"id" json:"id,omitempty"`
	Name         *string `json:"name,omitempty"`
	CoreDiameter float64 `db:"core_diameter" json:"core_diameter"`
	Unit         *string `json:"unit,omitempty"`
}
