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
	return json.Marshal(&struct {
		Date string `db:"date" json:"date,omitempty"`
		*Alias
	}{
		Date:  cl.Date.Format("01-02-2006"),
		Alias: (*Alias)(cl),
	})
}
