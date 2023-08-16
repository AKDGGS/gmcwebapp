package model

type ContainerLog struct {
	ID          int32          `json:"container_log_id,omitempty"`
	Destination string         `json:"destination,omitempty"`
	Date        *FormattedDate `json:"log_date,omitempty"`
}
