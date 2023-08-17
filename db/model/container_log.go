package model

import (
	"encoding/json"
	"time"
)

type ContainerLog struct {
	ID          int32     `json:"container_log_id,omitempty"`
	Destination string    `json:"destination,omitempty"`
	Date        time.Time `json:"log_date,omitempty"`
}

func (cl *ContainerLog) MarshalJSON() ([]byte, error) {
	type Alias ContainerLog
	var date string
	if !cl.Date.IsZero() {
		date = cl.Date.Format("01-02-2006")
	}
	return json.Marshal(&struct {
		Date string `json:"date,omitempty"`
		*Alias
	}{
		Date:  date,
		Alias: (*Alias)(cl),
	})
}
