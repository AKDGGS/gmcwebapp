package model

type Collection struct {
	ID             int32  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	OrganizationID int32  `json:"organization_id,omitempty"`
}
