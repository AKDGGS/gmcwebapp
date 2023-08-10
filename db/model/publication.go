package model

type Publication struct {
	ID                int32  `json:"publication_id,omitempty"`
	Title             string `json:"title,omitempty"`
	Description       string `json:"description,omitempty"`
	Year              int16  `json:"year,omitempty"`
	Type              string `json:"url_type,omitempty"`
	PublicationNumber string `json:"publication_number,omitempty"`
	PublicationSeries string `json:"publication_series,omitempty"`
	CanPublish        bool   `json:"can_publish"`
}
