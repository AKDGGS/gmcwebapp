package elastic

import (
	"fmt"

	"gmc/search/util"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Elastic struct {
	client *elastic.TypedClient
	index  string
}

func New(cfg map[string]interface{}) (*Elastic, error) {
	index, ok := cfg["index"].(string)
	if !ok {
		index = "gmc"
	}

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

	cli, err := elastic.NewTypedClient(escfg)
	if err != nil {
		return nil, err
	}

	return &Elastic{client: cli, index: index}, nil
}

func (es *Elastic) Shutdown() {}

func (es *Elastic) NewInventoryIndex() (util.InventoryIndex, error) {
	ii := ElasticInventoryIndex{}
	return &ii, nil
}

func (es *Elastic) createIndex(name string, tmap *types.TypeMapping) error {
	_, err := es.client.Indices.Create(name).Request(
		&create.Request{Mappings: tmap},
	).Do(nil)
	return err
}

func (es *Elastic) deleteIndex(name string) error {
	_, err := es.client.Indices.Delete(name).Do(nil)
	return err
}

func (es *Elastic) getIndicesFromAlias(name string) ([]string, error) {
	r, err := es.client.Indices.GetAlias().Name("stupid").Do(nil)
	if err != nil {
		return nil, err
	}

	indexes := []string{}
	for k, _ := range r {
		indexes = append(indexes, k)
	}
	return indexes, nil
}
