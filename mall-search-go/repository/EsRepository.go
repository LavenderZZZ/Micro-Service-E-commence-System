package repository

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"mall-search-go/model"
)

var es *elasticsearch.Client

func init() {
	cfg = elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
}

func SaveAll([]model.EsProduct) (int, error) {

}

func indexProduct(product model.EsProduct) bool {
	req := esapi.IndexRequest{
		Index: "products",
	}

}
