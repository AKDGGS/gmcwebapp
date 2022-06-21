package pg

import (
	"encoding/json"

	"gmc/assets"
)

func (pg *Postgres) ListQAReports() ([]map[string]string, error) {
	b, err := assets.ReadBytes("pg/qa/reports.json")
	if err != nil {
		return nil, err
	}

	var reports []map[string]string
	err = json.Unmarshal(b, &reports)
	if err != nil {
		return nil, err
	}
	return reports, nil
}
