package model

type URL struct {
	ID          int32   `json:"id,omitempty"`
	URL         *string `json:"url,omitempty"`
	Description *string `json:"description,omitempty"`
	Type        *string `json:"type"`
}
