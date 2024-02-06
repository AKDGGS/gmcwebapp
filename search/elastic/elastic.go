package elastic

import (
	"fmt"

	"gmc/search/util"

	elastic "github.com/elastic/go-elasticsearch/v8"
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
