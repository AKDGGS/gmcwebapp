package model

type Publication struct {
	ID                int32   `json:"id,omitempty"`
	Title             string  `json:"title,omitempty"`
	Description       *string `json:"description,omitempty"`
	Year              *int16  `json:"year,omitempty"`
	Type              *string `json:"publication_type,omitempty"`
	PublicationNumber *string `json:"publication_number,omitempty"`
	PublicationSeries *string `json:"publication_series,omitempty"`
	CanPublish        bool    `json:"can_publish"`
}
