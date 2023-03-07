package model

type Organization struct {
	ID      int32  `json:"organization_id"`
	Name    string `json:"name"`
	Remark  string `json:"remark"`
	Type    string `json:"organization_type"`
	Current bool   `json:"is_current"`
}
