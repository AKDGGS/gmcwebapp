package model

import (
	"encoding/json"
	"time"
)

type FormattedDate struct {
	time.Time
}

func (fd FormattedDate) MarshalJSON() ([]byte, error) {
	if fd.Time.IsZero() {
		return []byte("null"), nil
	}
	formattedDate := fd.Time.Format("01-02-2006")
	return json.Marshal(formattedDate)
}
