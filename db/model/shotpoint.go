package model

type Shotpoint struct {
	ID         int32    `json:"shotpoint_id,omitempty"`
	Number     float64  `json:"shotpoint_number,omitempty"`
	ShotlineID int32    `json:"shotline_id,omitempty"`
	Shotline   Shotline `json:"shotline,omitempty"`
}
