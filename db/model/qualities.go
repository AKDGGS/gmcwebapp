package model

import (
	"encoding/json"
	"time"
)

type Qualities struct {
	ID       int32     `json:"inventory_quality_id,omitempty"`
	Remark   string    `json:"remark,omitempty"`
	Date     time.Time `json:"check_date,omitempty"`
	Username string    `json:"username,omitempty"`
	Issues   []string  `json:"issues,omitempty"`
}

func (q *Qualities) MarshalJSON() ([]byte, error) {
	type Alias Qualities
	var date string
	if !q.Date.IsZero() {
		date = q.Date.Format("01-02-2006")
	}
	return json.Marshal(&struct {
		Date string `json:"date,omitempty"`
		*Alias
	}{
		Date:  date,
		Alias: (*Alias)(q),
	})
}
