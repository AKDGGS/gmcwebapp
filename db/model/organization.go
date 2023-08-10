package model

type Organization struct {
	ID      int32  `json:"organization_id,omitempty"`
	Name    string `json:"name,omitempty"`
	Remark  string `json:"remark,omitempty"`
	Type    string `json:"organization_type,omitempty"`
	Current bool   `json:"is_current"`
}
