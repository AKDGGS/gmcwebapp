package model

type Shotpoint struct {
	ShotpointID int32                  `json:"shotpoint_id"`
	ID          int32                  `json:"shotline_id"`
	Number      float64                `json:"shotpoint_number"`
	Stash       map[string]interface{} `json:"stash"`
}
