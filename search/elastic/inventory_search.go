package elastic

import (
	"encoding/json"
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
		id, err := strconv.Atoi(hit.Id_)
		if err != nil {
			return nil, err
		}

		ih := util.InventoryHit{ID: id}
		if err := hfType[string](hit, "collection", &ih.Collection); err != nil {
			return nil, err
		}
		if err := hfType[string](hit, "sample", &ih.SampleNumber); err != nil {
			return nil, err
		}
		if err := hfType[string](hit, "slide", &ih.SlideNumber); err != nil {
			return nil, err
		}
		if err := hfType[string](hit, "box", &ih.BoxNumber); err != nil {
			return nil, err
		}
		if err := hfType[string](hit, "set", &ih.SetNumber); err != nil {
			return nil, err
		}
		if err := hfType[string](hit, "core", &ih.CoreNumber); err != nil {
			return nil, err
		}
		if err := hfType[float64](hit, "diameter", &ih.CoreDiameter); err != nil {
			return nil, err
		}
		if err := hfType[float64](hit, "top", &ih.IntervalTop); err != nil {
			return nil, err
		}
		if err := hfType[float64](hit, "bottom", &ih.IntervalBottom); err != nil {
			return nil, err
		}
		if err := hfType[string](hit, "keywords", &ih.Keywords); err != nil {
			return nil, err
		}
		if err := hfType[string](hit, "barcode", &ih.Barcode); err != nil {
			return nil, err
		}
		res.Hits = append(res.Hits, ih)
	}

	return res, nil
}

// Returns a hit's field n as type T in ptr
func hfType[T any](hit types.Hit, n string, ptr *T) error {
	f, ok := hit.Fields[n]
	if !ok {
		return nil
	}

	var arr []T
	if err := json.Unmarshal(f, &arr); err != nil {
		return err
	}

	if len(arr) < 1 {
		return nil
	}

	*ptr = arr[0]
	return nil
}
