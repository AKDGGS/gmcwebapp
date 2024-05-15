package model

type Organization struct {
	ID     int32            `db:"id" json:"id,omitempty"`
	Name   string           `json:"name,omitempty"`
	Remark *string          `json:"remark,omitempty"`
	Type   OrganizationType `json:"type,omitempty"`
	// Only used by well_operator
	Current bool `db:"is_current" json:"is_current"`
}
