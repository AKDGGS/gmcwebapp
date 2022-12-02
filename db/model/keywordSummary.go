package model

type KeywordSummary struct {
	Keywords []string `json:"keywords"`
	Count    int64    `json:"count"`
}
