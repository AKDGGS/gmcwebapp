package model

type Collection struct {
	ID             int32  `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrganizationID int32  `json:"organization_id"`
}
