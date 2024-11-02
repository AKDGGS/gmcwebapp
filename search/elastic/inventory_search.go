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
	src_filter := &types.SourceFilter{
		Includes: []string{"*"},
		Excludes: []string{
			"barcode",
			"collection_id",
			"container_id",
			"core_name",
			"core_unit",
			"description",
			"note",
			"remark",
			"outcrop.year",
			"shotline.year",
			"borehole.name",
			"borehole.prospect.ardf",
			"borehole.prospect.name",
			"well.name",
			"publication.number",
			"publication.description",
			"publication.series",
		},
	}

	sea := es.client.Search().
		Index("inventory").
		From(params.From).
		Size(params.Size).
		TrackTotalHits(true).
		Source_(src_filter)

	if len(params.Sort) > 0 {
		var sort []types.SortCombinations
		for _, v := range params.Sort {
			if v[0] == "_score" {
				sort = append(sort, map[string]string{
					v[0]: v[1],
				})
			} else {
				sort = append(sort, map[string]map[string]string{
					v[0]: map[string]string{
						"order":   v[1],
						"missing": "_last",
					},
				})
			}
		}
		sea = sea.Sort(sort...)
	}

	qry := &types.BoolQuery{}

	if params.Query != "" {
		qry.Must = append(qry.Must, types.Query{
			QueryString: &types.QueryStringQuery{
				DefaultOperator: &operator.And,
				Query:           params.Query,
			},
		})
	}

	if len(params.Keywords) > 0 {
		bq := &types.BoolQuery{}
		for _, kw := range params.Keywords {
			bq.Must = append(bq.Must, types.Query{
				Term: map[string]types.TermQuery{
					"keyword": types.TermQuery{
						Value: kw,
					},
				},
			})
		}
		qry.Filter = append(qry.Filter, types.Query{
			Bool: bq,
		})
	}

	if !params.Private {
		src_filter.Excludes = append(
			src_filter.Excludes,
			"path_cache",
			"display_barcode",
			"can_publish",
			"issue",
		)

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
		From:    params.From,
		Total:   r.Hits.Total.Value,
		Time:    time.Duration(time.Duration(r.Took) * time.Millisecond),
		Private: params.Private,
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
