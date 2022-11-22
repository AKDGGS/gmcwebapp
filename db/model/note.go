package model

import "time"

type Note struct {
	ID       int32     `json:"note_id"`
	Note     string    `json:"note"`
	Date     time.Time `json:"note_date"`
	IsPublic bool      `json:"is_public"`
	Username string    `json:"username"`
}
