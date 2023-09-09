package model

type Collection struct {
	ID           int32        `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	Description  *string      `json:"description,omitempty"`
	Organization Organization `json:"organization,omitempty"`
}
