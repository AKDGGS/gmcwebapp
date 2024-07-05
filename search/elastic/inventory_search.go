package elastic

import (
	"encoding/json"
	"time"

	"gmc/db/model"
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
		Source_(&types.SourceFilter{
			Includes: []string{"*"},
			Excludes: []string{"remark", "well.altnames", "description",
				"project", "project_id", "collection_id",
				"core_name", "alt_barcode", "container_id",
				"c_barcode", "c_alt_barcode", "path_cache",
				"core_unit", "interval_unit", "barcode"},
		})

	qry := &types.BoolQuery{}

	if params.Query != "" {
		qry.Must = append(qry.Must, types.Query{
			QueryString: &types.QueryStringQuery{
				DefaultOperator: &operator.And,
				Query:           params.Query,
			},
		})
	}

	if !params.Private {
		qry.Filter = append(qry.Filter, types.Query{
			Match: map[string]types.MatchQuery{
				"can_publish": {Query: "true"},
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
		ih := model.FlatInventory{}
		if err := json.Unmarshal(hit.Source_, &ih); err != nil {
			return nil, err
		}
		res.Hits = append(res.Hits, ih)
	}

	return res, nil
}
