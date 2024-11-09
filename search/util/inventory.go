package util

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gmc/db/model"
)

type InventoryIndex interface {
	Add(*model.FlatInventory) error
	Count() int
	Commit() error
	Flush() error
	Rollback() error
}

type InventoryParams struct {
	Query          string
	Keywords       []string
	ProspectIDs    []int
	CollectionIDs  []int
	IntervalTop    *float64
	IntervalBottom *float64
	From           int
	Size           int
	Private        bool
	Sort           [][2]string
}

func (ip *InventoryParams) ParseQuery(q url.Values) {
	var err error
	ip.Size, err = strconv.Atoi(q.Get("size"))
	if err != nil {
		ip.Size = 25
	}

	ip.From, err = strconv.Atoi(q.Get("from"))
	if err != nil || ip.From < 0 {
		ip.From = 0
	}

	ip.Query = q.Get("q")

	if t := q.Get("top"); t != "" {
		if n, err := strconv.ParseFloat(t, 64); err == nil {
			ip.IntervalTop = &n
		}
	}

	if t := q.Get("bottom"); t != "" {
		if n, err := strconv.ParseFloat(t, 64); err == nil {
			ip.IntervalBottom = &n
		}
	}

	// Save people from flipping top and bottom values
	if ip.IntervalTop != nil && ip.IntervalBottom != nil {
		if *ip.IntervalTop > *ip.IntervalBottom {
			t := ip.IntervalTop
			ip.IntervalTop = ip.IntervalBottom
			ip.IntervalBottom = t
		}
	}

	if keywords, ok := q["keyword"]; ok {
		ip.Keywords = keywords
	}

	if ids, ok := q["prospect_id"]; ok {
		for _, sid := range ids {
			if id, err := strconv.Atoi(sid); err == nil {
				ip.ProspectIDs = append(ip.ProspectIDs, id)
			}
		}
	}

	if ids, ok := q["collection_id"]; ok {
		for _, sid := range ids {
			if id, err := strconv.Atoi(sid); err == nil {
				ip.CollectionIDs = append(ip.CollectionIDs, id)
			}
		}
	}

	if sorts, ok := q["sort"]; ok {
		dirs, _ := q["dir"]
		for i, v := range sorts {
			dir := "asc"
			if i < len(dirs) && strings.ToLower(dirs[i]) == "desc" {
				dir = "desc"
			}
			ip.Sort = append(ip.Sort, [2]string{v, dir})
		}
	}
}

type InventoryResults struct {
	Hits    []model.FlatInventory `json:"hits,omitempty"`
	From    int                   `json:"from"`
	Total   int64                 `json:"total"`
	Time    time.Duration         `json:"time"`
	Private bool                  `json:"private,omitempty"`
}

func (ir *InventoryResults) MarshalJSON() ([]byte, error) {
	type Alias InventoryResults

	time := ir.Time.String()
	return json.Marshal(&struct {
		Time string `json:"time"`
		*Alias
	}{
		Time:  time,
		Alias: (*Alias)(ir),
	})
}
