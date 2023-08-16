package model

type Note struct {
	ID       int32          `json:"note_id,omitempty"`
	Note     string         `json:"note,omitempty"`
	Date     *FormattedDate `json:"note_date,omitempty"`
	Public   bool           `json:"is_public"`
	Username string         `json:"username,omitempty"`
}
