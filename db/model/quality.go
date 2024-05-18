package model

import (
	"encoding/json"
	"time"
)

type Quality struct {
	ID       int32     `db:"id" json:"id"`
	Remark   *string   `json:"remark,omitempty"`
	Date     time.Time `db:"date" json:"date"`
	Username string    `json:"username"`
	Issues   []string  `json:"issues,omitempty"`
}

func (q *Quality) MarshalJSON() ([]byte, error) {
	type Alias Quality
	return json.Marshal(&struct {
		Date string `json:"date,omitempty"`
		*Alias
	}{
		Date:  q.Date.Format("01-02-2006"),
		Alias: (*Alias)(q),
	})
}
