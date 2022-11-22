package model

type Organization struct {
	ID               int32  `json:"organization_id"`
	Name             string `json:"name"`
	OrganizationType string `json:"organization_type"`
}
