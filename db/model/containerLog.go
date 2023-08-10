package model

import "time"

type ContainerLog struct {
	ID          int32      `json:"container_log_id,omitempty"`
	Destination string     `json:"destination,omitempty"`
	Date        *time.Time `json:"log_date,omitempty"`
}
