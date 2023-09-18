package model

import (
	"encoding/json"
	"time"
)

type Note struct {
	ID       int32     `json:"id,omitempty"`
	Note     string    `json:"note,omitempty"`
	Date     time.Time `json:"date,omitempty"`
	Public   bool      `json:"is_public"`
	Username string    `json:"username,omitempty"`
	NoteType NoteType  `json:"note_type,omitempty"`
}

func (n *Note) MarshalJSON() ([]byte, error) {
	type Alias Note
	var date string
	if !n.Date.IsZero() {
		date = n.Date.Format("01-02-2006")
	}
	return json.Marshal(&struct {
		Date string `json:"date,omitempty"`
		*Alias
	}{
		Date:  date,
		Alias: (*Alias)(n),
	})
}
