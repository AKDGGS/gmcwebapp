package model

import (
	"time"
)

type Quality struct {
	ID       int32     `db:"id" json:"id"`
	Remark   *string   `json:"remark,omitempty"`
	Date     time.Time `db:"date" json:"date"`
	Username string    `json:"username"`
	Issues   []string  `json:"issues,omitempty"`
}
