package elastic

import (
	"fmt"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Elastic struct {
	client *elastic.TypedClient
	index  string
}

func New(cfg map[string]interface{}) (*Elastic, error) {
	url, ok := cfg["url"].(string)
	if !ok {
		return nil, fmt.Errorf("elastic url must exist and be a string")
	}

	escfg := elastic.Config{Addresses: []string{url}}

	user, _ := cfg["user"].(string)
	pass, _ := cfg["pass"].(string)
	if user != "" && pass != "" {
		escfg.Username = user
		escfg.Password = pass
	}

	index, _ := cfg["index"].(string)
	if index == "" {
		index = "gmc"
	}

	cli, err := elastic.NewTypedClient(escfg)
	if err != nil {
		return nil, err
	}

	return &Elastic{client: cli, index: index}, nil
}

func (es *Elastic) Shutdown() {}

func (es *Elastic) Name() string {
	return "elastic"
}

func (es *Elastic) createIndex(name string, tmap *types.TypeMapping) error {
	window := 1000000
	_, err := es.client.Indices.Create(name).Request(
		&create.Request{
			Mappings: tmap,
			Settings: &types.IndexSettings{
				Analysis: &types.IndexSettingsAnalysis{
					Normalizer: map[string]types.Normalizer{
						"clean": map[string]interface{}{
							"type":   "custom",
							"filter": []string{"lowercase", "trim"},
						},
					},
				},
				MaxResultWindow: &window,
			},
		},
	).Do(nil)
	return err
}

func (es *Elastic) deleteIndex(name string) error {
	_, err := es.client.Indices.Delete(name).Do(nil)
	if err != nil {
		// Ignore errors from attempts to delete non-existant indexes
		if e, ok := err.(*types.ElasticsearchError); ok && e.Status == 404 {
			return nil
		}
		return err
	}
	return nil
}

func (es *Elastic) getIndicesFromAlias(name string) ([]string, error) {
	r, err := es.client.Indices.GetAlias().Name(name).Do(nil)
	if err != nil {
		// Ignore errors from attempts to get non-existant indexes
		if e, ok := err.(*types.ElasticsearchError); ok && e.Status == 404 {
			return nil, nil
		}
		return nil, err
	}

	indexes := []string{}
	for k, _ := range r {
		indexes = append(indexes, k)
	}
	return indexes, nil
}

func (es *Elastic) replaceAlias(name, index string) ([]string, error) {
	actions := []types.IndicesAction{
		types.IndicesAction{Add: &types.AddAction{
			Alias: &name, Index: &index,
		}},
	}

	old, err := es.getIndicesFromAlias(name)
	if err != nil {
		return nil, err
	}

	// Don't add a removal action if the existing alias doesn't exist
	if old != nil {
		actions = append(actions, types.IndicesAction{
			Remove: &types.RemoveAction{Alias: &name, Indices: old},
		})
	}

	_, err = es.client.Indices.UpdateAliases().Actions(actions...).Do(nil)
	if err != nil {
		return nil, err
	}
	return old, nil
}

func (es *Elastic) bulk(index string) *bulk.Bulk {
	return es.client.Bulk().Index(index)
}

func (es *Elastic) refreshIndex(index string) error {
	_, err := es.client.Indices.Refresh().Index(index).Do(nil)
	return err
}
