package model

import (
	"encoding/json"
	"time"
)

type Note struct {
	ID       int32     `db:"id" json:"id,omitempty"`
	Note     string    `json:"note,omitempty"`
	Date     time.Time `db:"date" json:"date,omitempty"`
	Public   bool      `db:"is_public" json:"is_public"`
	Username string    `json:"username,omitempty"`
	NoteType struct {
		ID          int32  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"note_type,omitempty"`
}

func (n *Note) MarshalJSON() ([]byte, error) {
	type Alias Note
	return json.Marshal(&struct {
		Date string `json:"date,omitempty"`
		*Alias
	}{
		Date:  n.Date.Format("01-02-2006"),
		Alias: (*Alias)(n),
	})
}
