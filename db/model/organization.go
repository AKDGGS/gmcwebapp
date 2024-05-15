package model

type Organization struct {
	ID      int32            `json:"id,omitempty"`
	Name    string           `json:"name,omitempty"`
	Remark  *string          `json:"remark,omitempty"`
	Current bool             `json:"is_current"`
	Type    OrganizationType `json:"type,omitempty"`
}
