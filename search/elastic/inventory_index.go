package elastic

import (
	"fmt"
	"time"

	"gmc/db/model"
	"gmc/search/util"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

var (
	yes = true
	no  = false
)

func (es *Elastic) NewInventoryIndex() (util.InventoryIndex, error) {
	iname := fmt.Sprintf("inventory-%x", time.Now().UnixMicro())
	err := es.createIndex(iname,
		&types.TypeMapping{
			Properties: map[string]types.Property{
				"collection": &types.TextProperty{Store: &yes, Index: &yes},
				"barcode":    &types.TextProperty{Store: &yes, Index: &yes},
				"remark":     &types.TextProperty{Store: &no, Index: &yes},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return &ElasticInventoryIndex{
		es: es, bulk: es.bulk(iname), index: iname,
	}, nil
}

type ElasticInventoryIndex struct {
	es    *Elastic
	bulk  *bulk.Bulk
	count int
	index string
}

func (ii *ElasticInventoryIndex) Count() int {
	return ii.count
}

func (ii *ElasticInventoryIndex) Add(f *model.FlatInventory) error {
	sid := f.StringID()
	err := ii.bulk.IndexOp(types.IndexOperation{Id_: &sid}, f)
	if err == nil {
		ii.count++
	}
	return err
}

func (ii *ElasticInventoryIndex) Flush() error {
	if ii.count > 0 {
		if _, err := ii.bulk.Do(nil); err != nil {
			fmt.Printf("\n%+v\n", err)
			return err
		}
		ii.count = 0
		ii.bulk = ii.es.bulk(ii.index)
	}
	return nil
}

func (ii *ElasticInventoryIndex) Rollback() error {
	return ii.es.deleteIndex(ii.index)
}

func (ii *ElasticInventoryIndex) Commit() error {
	// Flush any remaining bulk updates
	if err := ii.Flush(); err != nil {
		return err
	}
	// Refresh the newly created index
	if err := ii.es.refreshIndex(ii.index); err != nil {
		return err
	}
	// Replace the inventory alias with the newly created index
	old, err := ii.es.replaceAlias("inventory", ii.index)
	if err != nil {
		return err
	}
	// Delete old indexes previously pointing to inventory alias
	for _, idx := range old {
		if err := ii.es.deleteIndex(idx); err != nil {
			return err
		}
	}
	return nil
}
