package model

import (
	"encoding/json"
	"time"
)

type Note struct {
	ID       int32      `json:"note_id"`
	Note     string     `json:"note"`
	Date     *time.Time `json:"note_date"`
	Public   bool       `json:"is_public"`
	Username string     `json:"username"`
}

func (n *Note) MarshalJSON() ([]byte, error) {
	type Alias Note
	return json.Marshal(&struct {
		Date string `json:"date"`
		*Alias
	}{
		Date:  n.Date.Format("01-02-2006"),
		Alias: (*Alias)(n),
	})
}
