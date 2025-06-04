package elastic

import (
	"encoding/json"
	"fmt"
	"time"

	"gmc/db/model"
	"gmc/search/util"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/dynamicmapping"
)

var (
	yes   = true
	no    = false
	clean = "clean"
)

func fieldAlias(s string) *types.FieldAliasProperty {
	return &types.FieldAliasProperty{Path: &s}
}

func (es *Elastic) InventorySortByFields(full bool) [][2]string {
	if full {
		return [][2]string{
			[2]string{"_score", "Best Match"},
			[2]string{"borehole.name_sort", "Borehole"},
			[2]string{"box.sort", "Box"},
			[2]string{"collection.sort", "Collection"},
			[2]string{"core.sort", "Core Number"},
			[2]string{"keyword.sort", "Keywords"},
			[2]string{"prospect.name_sort", "Prospect"},
			[2]string{"sample.sort", "Sample"},
			[2]string{"set.sort", "Set Number"},
			[2]string{"top", "Top"},
			[2]string{"bottom", "Bottom"},
			[2]string{"well.name_sort", "Well"},
			[2]string{"well.number.sort", "Well Number"},
			[2]string{"path_cache.sort", "Location"},
			[2]string{"barcode", "Barcode"},
		}
	}
	return [][2]string{
		[2]string{"_score", "Best Match"},
		[2]string{"borehole.name_sort", "Borehole"},
		[2]string{"box.sort", "Box"},
		[2]string{"collection.sort", "Collection"},
		[2]string{"core.sort", "Core Number"},
		[2]string{"keyword.sort", "Keywords"},
		[2]string{"prospect.name_sort", "Prospect"},
		[2]string{"sample.sort", "Sample"},
		[2]string{"set.sort", "Set Number"},
		[2]string{"top", "Top"},
		[2]string{"bottom", "Bottom"},
		[2]string{"well.name_sort", "Well"},
		[2]string{"well.number.sort", "Well Number"},
	}
}

func (es *Elastic) NewInventoryIndex() (util.InventoryIndex, error) {
	// Any field used for sorting needs to be normalized, even if it's
	// not indexed
	iname := fmt.Sprintf("%s-inventory-%x", es.index, time.Now().UnixMicro())
	err := es.createIndex(iname,
		&types.TypeMapping{
			Properties: map[string]types.Property{
				"collection": &types.TextProperty{
					Index: &yes,
					Fields: map[string]types.Property{
						"sort": &types.KeywordProperty{Index: &yes, Normalizer: &clean},
					},
				},
				"collection_id": &types.IntegerNumberProperty{Index: &yes},
				"slide":         &types.TextProperty{Index: &yes},
				"sample": &types.TextProperty{
					Index: &yes,
					Fields: map[string]types.Property{
						"sort": &types.KeywordProperty{Index: &yes, Normalizer: &clean},
					},
				},
				"box": &types.TextProperty{
					Index: &yes,
					Fields: map[string]types.Property{
						"sort": &types.KeywordProperty{Index: &yes, Normalizer: &clean},
					},
				},
				"set": &types.TextProperty{
					Index: &yes,
					Fields: map[string]types.Property{
						"sort": &types.KeywordProperty{Index: &yes, Normalizer: &clean},
					},
				},
				"core": &types.TextProperty{
					Index: &yes,
					Fields: map[string]types.Property{
						"sort": &types.KeywordProperty{Index: &yes, Normalizer: &clean},
					},
				},
				"diameter":  &types.FloatNumberProperty{Index: &yes},
				"core_name": &types.TextProperty{Index: &yes},
				"core_unit": &types.TextProperty{Index: &yes},
				"interval":  &types.FloatRangeProperty{Index: &yes},
				"top":       &types.FloatNumberProperty{Index: &yes},
				"bottom":    &types.FloatNumberProperty{Index: &yes},
				"unit":      &types.TextProperty{Index: &yes},
				"keyword": &types.TextProperty{
					Index: &yes,
					Fields: map[string]types.Property{
						"sort": &types.KeywordProperty{Index: &yes, Normalizer: &clean},
					},
				},
				"barcode":      &types.KeywordProperty{Index: &yes},
				"alt_barcode":  &types.KeywordProperty{Index: &yes},
				"container_id": &types.IntegerNumberProperty{Index: &yes},
				"path_cache": &types.TextProperty{
					Index: &yes,
					Fields: map[string]types.Property{
						"sort": &types.KeywordProperty{Index: &yes, Normalizer: &clean},
					},
				},
				"remark":      &types.TextProperty{Index: &yes},
				"project_id":  &types.IntegerNumberProperty{Index: &yes},
				"project":     &types.TextProperty{Index: &yes},
				"can_publish": &types.BooleanProperty{Index: &yes},
				"description": &types.TextProperty{Index: &yes},
				"geometries":  &types.GeoShapeProperty{},
				"longitude":   &types.FloatNumberProperty{Index: &no},
				"latitude":    &types.FloatNumberProperty{Index: &no},
				"note":        &types.TextProperty{Index: &yes},
				"issue":       &types.TextProperty{Index: &yes},
				"well": &types.ObjectProperty{
					Dynamic: &dynamicmapping.False,
					Enabled: &yes,
					Properties: map[string]types.Property{
						"id":         &types.IntegerNumberProperty{Index: &yes},
						"name":       &types.TextProperty{Index: &yes, CopyTo: []string{"well.name_sort"}},
						"name_sort":  &types.KeywordProperty{Index: &yes, Normalizer: &clean},
						"alt_names":  &types.TextProperty{Index: &yes, CopyTo: []string{"well.name"}},
						"api":        &types.KeywordProperty{Index: &yes, Normalizer: &clean},
						"is_onshore": &types.BooleanProperty{Index: &yes},
						"number": &types.TextProperty{
							Index: &yes,
							Fields: map[string]types.Property{
								"sort": &types.KeywordProperty{Index: &yes, Normalizer: &clean},
							},
						},
					},
				},
				"outcrop": &types.ObjectProperty{
					Dynamic: &dynamicmapping.False,
					Enabled: &yes,
					Properties: map[string]types.Property{
						"id":     &types.IntegerNumberProperty{Index: &yes},
						"name":   &types.TextProperty{Index: &yes},
						"number": &types.TextProperty{Index: &yes},
						"year":   &types.IntegerNumberProperty{Index: &yes},
					},
				},
				"borehole": &types.ObjectProperty{
					Dynamic: &dynamicmapping.False,
					Enabled: &yes,
					Properties: map[string]types.Property{
						"id":        &types.IntegerNumberProperty{Index: &yes},
						"name":      &types.TextProperty{Index: &yes, CopyTo: []string{"borehole.name_sort"}},
						"name_sort": &types.KeywordProperty{Index: &yes, Normalizer: &clean},
						"alt_names": &types.TextProperty{Index: &yes, CopyTo: []string{"borehole.name"}},
						"prospect": &types.ObjectProperty{
							Dynamic: &dynamicmapping.False,
							Enabled: &yes,
							Properties: map[string]types.Property{
								"id":        &types.IntegerNumberProperty{Index: &yes},
								"name":      &types.TextProperty{Index: &yes, CopyTo: []string{"borehole.prospect.name_sort"}},
								"name_sort": &types.KeywordProperty{Index: &yes, Normalizer: &clean},
								"alt_names": &types.TextProperty{Index: &yes, CopyTo: []string{"borehole.prospect.name"}},
								"ardf":      &types.TextProperty{Index: &yes},
							},
						},
					},
				},
				"prospect": &types.ObjectProperty{
					Dynamic: &dynamicmapping.False,
					Enabled: &yes,
					Properties: map[string]types.Property{
						"id":        fieldAlias("borehole.prospect.id"),
						"name":      fieldAlias("borehole.prospect.name"),
						"name_sort": fieldAlias("borehole.prospect.name_sort"),
						"ardf":      fieldAlias("borehole.prospect.ardf"),
					},
				},
				"shotline": &types.ObjectProperty{
					Dynamic: &dynamicmapping.False,
					Enabled: &yes,
					Properties: map[string]types.Property{
						"id":   &types.IntegerNumberProperty{Index: &yes},
						"name": &types.TextProperty{Index: &yes},
						"year": &types.IntegerNumberProperty{Index: &yes},
						"min":  &types.FloatNumberProperty{Index: &yes},
						"max":  &types.FloatNumberProperty{Index: &yes},
					},
				},
				"publication": &types.ObjectProperty{
					Dynamic: &dynamicmapping.False,
					Enabled: &yes,
					Properties: map[string]types.Property{
						"id":          &types.IntegerNumberProperty{Index: &yes},
						"title":       &types.TextProperty{Index: &yes},
						"description": &types.TextProperty{Index: &yes},
						"year":        &types.IntegerNumberProperty{Index: &yes},
						"number":      &types.TextProperty{Index: &yes},
						"series":      &types.TextProperty{Index: &yes},
					},
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return &ElasticInventoryIndex{
		es: es, bulk: es.bulk(iname), index: iname,
		alias: fmt.Sprintf("%s-inventory", es.index),
	}, nil
}

type ElasticInventoryIndex struct {
	es    *Elastic
	bulk  *bulk.Bulk
	count int
	index string
	alias string
}

func (ii *ElasticInventoryIndex) Count() int {
	return ii.count
}

func (ii *ElasticInventoryIndex) Add(f *model.FlatInventory) error {
	sid := f.StringID()
	js, err := json.Marshal(f)
	if err != nil {
		return err
	}

	err = ii.bulk.IndexOp(types.IndexOperation{Id_: &sid}, js)
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
	old, err := ii.es.replaceAlias(ii.alias, ii.index)
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
