package model

type Shotpoint struct {
	ID       int32    `json:"id,omitempty"`
	Number   *float64 `json:"number,omitempty"`
	Shotline Shotline `json:"shotline,omitempty"`
}
