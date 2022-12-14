package model

type Prospect struct {
	ID         int32  `json:"prospect_id"`
	Name       string `json:"prospect_name"`
	AltNames   string `json:"alt_names"`
	ARDFNumber string `json:"ardf_number"`
}
