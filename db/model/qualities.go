package model

type Qualities struct {
	ID       int32          `json:"inventory_quality_id,omitempty"`
	Remark   string         `json:"remark,omitempty"`
	Date     *FormattedDate `json:"check_date,omitempty"`
	Username string         `json:"username,omitempty"`
	Issues   []string       `json:"issues,omitempty"`
}
