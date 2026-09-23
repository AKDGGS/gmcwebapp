package model

const (
	BOREHOLE = iota
	INVENTORY
	OUTCROP
	PROSPECT
	WELL
)

type Table struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
}
