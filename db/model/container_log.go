package model

import (
	"encoding/json"
	"time"
)

type ContainerLog struct {
	ID          int32     `json:"id,omitempty"`
	Date        time.Time `db:"date" json:"date,omitempty"`
	Destination string    `json:"destination,omitempty"`
}

func (cl *ContainerLog) MarshalJSON() ([]byte, error) {
	type Alias ContainerLog
	var date string
	if !cl.Date.IsZero() {
		date = cl.Date.Format("01-02-2006")
	}
	return json.Marshal(&struct {
		Date string `db:"date" json:"date,omitempty"`
		*Alias
	}{
		Date:  date,
		Alias: (*Alias)(cl),
	})
}
