package model

import (
	"encoding/json"
	"time"
)

type Qualities struct {
	ID       int32     `json:"id"`
	Remark   *string   `json:"remark,omitempty"`
	Date     time.Time `json:"check_date"`
	Username string    `json:"username"`
	Issues   []string  `json:"issues,omitempty"`
}

func (q *Qualities) MarshalJSON() ([]byte, error) {
	type Alias Qualities
	var date string
	if !q.Date.IsZero() {
		date = q.Date.Format("01-02-2006")
	}
	return json.Marshal(&struct {
		Date string `json:"check_date,omitempty"`
		*Alias
	}{
		Date:  date,
		Alias: (*Alias)(q),
	})
}
