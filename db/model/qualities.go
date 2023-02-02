package model

import (
	"encoding/json"
	"time"
)

type Qualities struct {
	ID       int32      `json:"inventory_quality_id"`
	Remark   string     `json:"remark"`
	Date     *time.Time `json:"check_date"`
	Username string     `json:"username"`
	Issues   []string   `json:"issues"`
}

func (q *Qualities) MarshalJSON() ([]byte, error) {
	type Alias Qualities
	return json.Marshal(&struct {
		Date string `json:"check_date"`
		*Alias
	}{
		Date:  q.Date.Format("01-02-2006"),
		Alias: (*Alias)(q),
	})
}
