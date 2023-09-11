package model

type NoteType struct {
	ID          int32  `json:"type_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
