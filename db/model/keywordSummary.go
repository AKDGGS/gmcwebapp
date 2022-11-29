package model

type KeywordSummary struct {
	Keywords []interface{} `json:"keywords"`
	Count    int64         `json:"count"`
}
