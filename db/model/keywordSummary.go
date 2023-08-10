package model

type KeywordSummary struct {
	Keywords []string `json:"keywords,omitempty"`
	Count    int64    `json:"count,omitempty"`
}
