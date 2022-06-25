package model

type Table struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}
