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

func (es *Elastic) InventorySortByFields() [][2]string {
	return [][2]string{
		[2]string{"_score", "Best Match"},
		[2]string{"borehole.display_name", "Borehole"},
		[2]string{"box", "Box"},
		[2]string{"collection", "Collection"},
		[2]string{"core", "Core Number"},
		[2]string{"keyword", "Keywords"},
		[2]string{"prospect.display_name", "Prospect"},
		[2]string{"sample", "Sample"},
		[2]string{"set", "Set Number"},
		[2]string{"top", "Top"},
		[2]string{"bottom", "Bottom"},
		[2]string{"well.display_name", "Well"},
		[2]string{"well.number", "Well Number"},
	}
}

func (es *Elastic) NewInventoryIndex() (util.InventoryIndex, error) {
	iname := fmt.Sprintf("inventory-%x", time.Now().UnixMicro())
	err := es.createIndex(iname,
		&types.TypeMapping{
			Properties: map[string]types.Property{
				"collection":    &types.KeywordProperty{Index: &yes, Normalizer: &clean},
				"collection_id": &types.IntegerNumberProperty{Index: &yes},
				"sample":        &types.KeywordProperty{Index: &yes, Normalizer: &clean},
				"slide":         &types.TextProperty{Index: &yes},
				"box":           &types.KeywordProperty{Index: &yes, Normalizer: &clean},
				"set":           &types.KeywordProperty{Index: &yes, Normalizer: &clean},
				"core":          &types.KeywordProperty{Index: &yes, Normalizer: &clean},
				"diameter":      &types.FloatNumberProperty{Index: &yes},
				"core_name":     &types.TextProperty{Index: &yes},
				"core_unit":     &types.TextProperty{Index: &yes},
				"top":           &types.FloatNumberProperty{Index: &yes},
				"bottom":        &types.FloatNumberProperty{Index: &yes},
				"unit":          &types.TextProperty{Index: &yes},
				"keyword":       &types.KeywordProperty{Index: &yes, Normalizer: &clean},
				"barcode":       &types.TextProperty{Index: &yes},
				"container_id":  &types.IntegerNumberProperty{Index: &yes},
				"path_cache":    &types.TextProperty{Index: &yes},
				"remark":        &types.TextProperty{Index: &yes},
				"project_id":    &types.IntegerNumberProperty{Index: &yes},
				"project":       &types.TextProperty{Index: &yes},
				"can_publish":   &types.BooleanProperty{Index: &yes},
				"description":   &types.TextProperty{Index: &yes},
				"note":          &types.TextProperty{Index: &yes},
				"issue":         &types.TextProperty{Index: &yes},
				"well": &types.ObjectProperty{
					Dynamic: &dynamicmapping.False,
					Enabled: &yes,
					Properties: map[string]types.Property{
						"id":           &types.IntegerNumberProperty{Index: &yes},
						"name":         &types.KeywordProperty{Index: &yes},
						"display_name": &types.KeywordProperty{Index: &no, Normalizer: &clean},
						"number":       &types.KeywordProperty{Index: &yes, Normalizer: &clean},
						"api":          &types.TextProperty{Index: &yes},
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
						"id":           &types.IntegerNumberProperty{Index: &yes},
						"name":         &types.KeywordProperty{Index: &yes},
						"display_name": &types.KeywordProperty{Index: &no, Normalizer: &clean},
						"prospect": &types.ObjectProperty{
							Dynamic: &dynamicmapping.False,
							Enabled: &yes,
							Properties: map[string]types.Property{
								"id":           &types.IntegerNumberProperty{Index: &yes},
								"name":         &types.KeywordProperty{Index: &yes},
								"display_name": &types.KeywordProperty{Index: &no, Normalizer: &clean},
								"ardf":         &types.TextProperty{Index: &yes},
							},
						},
					},
				},
				"prospect": &types.ObjectProperty{
					Dynamic: &dynamicmapping.False,
					Enabled: &yes,
					Properties: map[string]types.Property{
						"id":           fieldAlias("borehole.prospect.id"),
						"name":         fieldAlias("borehole.prospect.name"),
						"display_name": fieldAlias("borehole.prospect.display_name"),
						"ardf":         fieldAlias("borehole.prospect.ardf"),
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
