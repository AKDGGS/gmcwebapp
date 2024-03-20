package elastic

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gmc/search/util"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/operator"
)

func (es *Elastic) SearchInventory(params *util.InventoryParams) (*util.InventoryResults, error) {
	sea := es.client.Search().
		Index("inventory").
		From(params.From).
		Size(params.Size).
		TrackTotalHits(true).
		StoredFields([]string{"*"}...).
		Source_(nil)

	qry := &types.BoolQuery{}

	if params.Query != "" {
		qry.Must = append(qry.Filter, types.Query{
			QueryString: &types.QueryStringQuery{
				DefaultOperator: &operator.And,
				Query:           params.Query,
			},
		})
	}

	r, err := sea.Query(&types.Query{Bool: qry}).Do(nil)
	if err != nil {
		return nil, err
	}

	res := &util.InventoryResults{
		From:  params.From,
		Total: r.Hits.Total.Value,
		Time:  time.Duration(time.Duration(r.Took) * time.Millisecond),
	}

	for _, hit := range r.Hits.Hits {
		id, err := strconv.Atoi(hit.Id_)
		if err != nil {
			return nil, err
		}

		ih := util.InventoryHit{ID: id}
		if err := hfStr(hit, "collection", &ih.Collection); err != nil {
			return nil, err
		}
		res.Hits = append(res.Hits, ih)
	}

	return res, nil
}

// Returns a hit's field n as string in ptr
func hfStr(hit types.Hit, n string, ptr *string) error {
	f, ok := hit.Fields[n]
	if !ok {
		return fmt.Errorf("unknown field: %s", n)
	}

	var arr []string
	if err := json.Unmarshal(f, &arr); err != nil {
		return err
	}

	if len(arr) < 1 {
		return nil
	}

	*ptr = arr[0]
	return nil
}
