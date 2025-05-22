package elastic

import (
	"encoding/json"
	"fmt"
	"time"

	"gmc/db/model"
	"gmc/search/util"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/geoshaperelation"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/operator"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
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
			"note",
			"remark",
			"interval",
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

	if !params.IncludeDescription {
		src_filter.Excludes = append(src_filter.Excludes, "description")
	}

	if !params.IncludeLatLon {
		src_filter.Excludes = append(src_filter.Excludes, "latitude")
		src_filter.Excludes = append(src_filter.Excludes, "longitude")
	}

	sea := es.client.Search().
		Index(fmt.Sprintf("%s-inventory", es.index)).
		From(params.From).
		Size(params.Size).
		TrackTotalHits(true).
		Source_(src_filter)

	if len(params.Sort) > 0 {
		var sort []types.SortCombinations
		for _, v := range params.Sort {
			if v[0] == "_score" || v[0] == "" {
				sort = append(sort, map[string]string{
					"_score": v[1],
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
				Lenient:         &yes,
				Type:            &textquerytype.TextQueryType{Name: "cross_fields"},
				DefaultOperator: &operator.And,
				Query:           params.Query,
			},
		})
	}

	if params.GeoJSON != "" {
		qry.Must = append(qry.Must, types.Query{
			GeoShape: &types.GeoShapeQuery{
				GeoShapeQuery: map[string]types.GeoShapeFieldQuery{
					"geometries": types.GeoShapeFieldQuery{
						Relation: &geoshaperelation.Intersects,
						Shape:    json.RawMessage(params.GeoJSON),
					},
				},
			},
		})
	}

	if len(params.Keywords) > 0 {
		bq := &types.BoolQuery{MinimumShouldMatch: 1}
		for _, kw := range params.Keywords {
			bq.Should = append(bq.Should, types.Query{
				MatchPhrase: map[string]types.MatchPhraseQuery{
					"keyword.keyword": types.MatchPhraseQuery{
						Query: kw,
					},
				},
			})
		}
		qry.Filter = append(qry.Filter, types.Query{Bool: bq})
	}

	if len(params.CollectionIDs) > 0 {
		bq := &types.BoolQuery{MinimumShouldMatch: 1}
		for _, cid := range params.CollectionIDs {
			bq.Should = append(bq.Should, types.Query{
				Term: map[string]types.TermQuery{
					"collection_id": types.TermQuery{
						Value: cid,
					},
				},
			})
		}
		qry.Filter = append(qry.Filter, types.Query{Bool: bq})
	}

	if params.IntervalTop != nil || params.IntervalBottom != nil {
		nrq := types.NumberRangeQuery{}
		if params.IntervalTop != nil {
			n := types.Float64(*params.IntervalTop)
			nrq.Gte = &n
		}
		if params.IntervalBottom != nil {
			n := types.Float64(*params.IntervalBottom)
			nrq.Lte = &n
		}

		qry.Filter = append(qry.Filter, types.Query{
			Range: map[string]types.RangeQuery{
				"interval": nrq,
			},
		})
	}

	if len(params.ProspectIDs) > 0 {
		bq := &types.BoolQuery{MinimumShouldMatch: 1}
		for _, cid := range params.ProspectIDs {
			bq.Should = append(bq.Should, types.Query{
				Term: map[string]types.TermQuery{
					"prospect.id": types.TermQuery{
						Value: cid,
					},
				},
			})
		}
		qry.Filter = append(qry.Filter, types.Query{Bool: bq})
	}

	// Exclude private inventory if needed
	if !params.IncludePrivate {
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
		Private: params.IncludePrivate,
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
