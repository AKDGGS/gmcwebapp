package model

import "time"

type ContainerLog struct {
	ID          int32      `json:"container_log_id"`
	Destination string     `json:"destination"`
	Date        *time.Time `json:"log_date"`
}
