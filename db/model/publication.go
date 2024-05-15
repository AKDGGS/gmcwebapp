package model

type Publication struct {
	ID                int32   `json:"id,omitempty"`
	Title             string  `json:"title,omitempty"`
	Description       *string `json:"description,omitempty"`
	Year              *int32  `json:"year,omitempty"`
	Type              *string `db:"publication_type" json:"publication_type,omitempty"`
	PublicationNumber *string `db:"publication_number" json:"publication_number,omitempty"`
	PublicationSeries *string `db:"publication_series" json:"publication_series,omitempty"`
	CanPublish        bool    `db:"can_publish" json:"can_publish"`
}
