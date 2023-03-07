package model

type Publication struct {
	ID                int32  `json:"publication_id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Year              int16  `json:"year"`
	Type              string `json:"url_type"`
	PublicationNumber string `json:"publication_number"`
	PublicationSeries string `json:"publication_series"`
	CanPublish        bool   `json:"can_publish"`
}
