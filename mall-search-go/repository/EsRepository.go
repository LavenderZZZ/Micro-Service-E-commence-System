package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/mitchellh/mapstructure"
	"log"
	"mall-search-go/model"
	"strconv"
)

type EsProductRepository interface {
	SaveAll([]model.EsProduct) (int, error)
	Save(*model.EsProduct) (*model.EsProduct, error)
	Delete(int64) error
	DeletaBatch([]int64) error
	Search(keyword string, pageNum, pageSize int) (model.Page, error)
	SearchById(keyword string, brandId *int64, productCategoryId *int64, pageNum int, pageSize int, sort int) (model.Page, error)
}

type esProductRepositoryImpl struct {
	client *elasticsearch.Client
	index  string
}

func init() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	Repo = NewEsProductRepository(es, "pms")
}

var Repo EsProductRepository

func NewEsProductRepository(client *elasticsearch.Client, index string) EsProductRepository {
	return &esProductRepositoryImpl{client: client, index: index}
}

func (repo *esProductRepositoryImpl) SaveAll(products []model.EsProduct) (int, error) {

	var buf bytes.Buffer
	for _, product := range products {
		meta := []byte(`{"index" : {"_id" : "` + strconv.FormatInt(product.ID, 10) + `" }} ` + "\n")
		data, err := json.Marshal(product)
		if err != nil {
			return 0, err
		}
		data = append(data, "\n"...)
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
	}
	req := esapi.BulkRequest{
		Index:   repo.index,
		Body:    &buf,
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), repo.client)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}
	return len(products), nil
}

func (repo *esProductRepositoryImpl) Save(product *model.EsProduct) (*model.EsProduct, error) {
	req := esapi.IndexRequest{
		Index:      repo.index,
		DocumentID: strconv.FormatInt(product.ID, 10),
		Body:       esutil.NewJSONReader(product),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), repo.client)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("Error indexing document ID=%d: %s", product.ID, res.Status())
	}

	return product, nil
}

func (repo *esProductRepositoryImpl) Delete(id int64) error {
	req := esapi.DeleteRequest{
		Index:      repo.index,
		DocumentID: strconv.FormatInt(id, 10),
	}
	res, err := req.Do(context.Background(), repo.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Error: %s", res.String())
	}
	return nil
}

func (repo *esProductRepositoryImpl) DeletaBatch(ids []int64) error {
	var buf bytes.Buffer
	for _, id := range ids {
		meta := []byte(`{"delete" : {"_id" : "` + strconv.FormatInt(id, 10) + `" }} ` + "\n")
		buf.Grow(len(meta))
		buf.Write(meta)
	}

	req := esapi.BulkRequest{
		Index:   repo.index,
		Body:    &buf,
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), repo.client)
	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("Error deleting documents: %s", res.Status())
	}
	return nil
}

func (repo *esProductRepositoryImpl) Search(keyword string, pageNum, pageSize int) (model.Page, error) {
	var result model.Page

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{"match": map[string]interface{}{"name": keyword}},
					{"match": map[string]interface{}{"subTitle": keyword}},
					{"match": map[string]interface{}{"keywords": keyword}},
				},
			},
		},
		"from": (pageNum - 1) * pageSize,
		"size": pageSize,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return result, err
	}

	res, err := repo.client.Search(
		repo.client.Search.WithContext(context.Background()),
		repo.client.Search.WithIndex(repo.index),
		repo.client.Search.WithBody(&buf),
		repo.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return result, err
	}

	var searchResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return result, err
	}

	var products []model.EsProduct
	for _, hit := range searchResult["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var product model.EsProduct
		err := mapstructure.Decode(hit.(map[string]interface{})["_source"], &product)
		if err != nil {
			return result, err
		}
		products = append(products, product)
	}

	totalHits := int(searchResult["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	totalPages := (totalHits + pageSize - 1) / pageSize

	result.Content = products
	result.PageInfo = model.PageInfo{
		TotalElements: totalHits,
		TotalPages:    totalPages,
		Number:        pageNum,
		Size:          pageSize,
	}
	return result, nil
}

func (repo *esProductRepositoryImpl) SearchById(keyword string, brandId *int64, productCategoryId *int64, pageNum int, pageSize int, sort int) (model.Page, error) {
	var result model.Page
	query := make(map[string]interface{})

	if keyword != "" {
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	} else {
		query["query"] = map[string]interface{}{
			"function_score": map[string]interface{}{
				"query": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":  keyword,
						"fields": []string{"name^10", "subTitle", "keywords"},
					},
				},
				"score_mode": "sum",
				"min_score":  2,
			},
		}
	}

	boolFilter := map[string]interface{}{
		"filter": []map[string]interface{}{},
	}

	if brandId != nil {
		boolFilter["filter"] = append(boolFilter["filter"].([]map[string]interface{}), map[string]interface{}{
			"term": map[string]interface{}{
				"brandId": brandId,
			},
		})
	}

	if productCategoryId != nil {
		boolFilter["filter"] = append(boolFilter["filter"].([]map[string]interface{}), map[string]interface{}{
			"term": map[string]interface{}{
				"productCategoryId": productCategoryId,
			},
		})
	}

	query["query"] = map[string]interface{}{"bool": boolFilter}

	//Sorting
	var sortField string
	var sortOrder string
	switch sort {
	case 1:
		sortField = "id"
		sortOrder = "desc"
	case 2:
		sortField = "sale"
		sortOrder = "desc"
	case 3:
		sortField = "price"
		sortOrder = "asc"
	case 4:
		sortField = "price"
		sortOrder = "desc"
	default:
		sortField = "_score"
		sortOrder = "desc"
	}

	query["sort"] = []map[string]interface{}{
		{
			sortField: map[string]interface{}{
				"order": sortOrder,
			},
		},
	}

	//Pagination
	query["from"] = (pageNum - 1) * pageSize
	query["size"] = pageSize

	//Convert query to JSON and make the request
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return result, fmt.Errorf("Error encoding query: %s", err)
	}

	res, err := repo.client.Search(
		repo.client.Search.WithContext(context.Background()),
		repo.client.Search.WithIndex(repo.index),
		repo.client.Search.WithBody(&buf),
	)
	if err != nil {
		return result, fmt.Errorf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	var searchResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return result, fmt.Errorf("Error parsing the response body: %s", err)
	}

	var products []model.EsProduct
	for _, hits := range searchResult["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var product model.EsProduct
		err := mapstructure.Decode(hits.(map[string]interface{})["_source"], &product)
		if err != nil {
			return result, err
		}
		products = append(products, product)
	}

	totalHits := int(searchResult["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	totalPages := (totalHits + pageSize - 1) / pageSize

	result.Content = products
	result.PageInfo = model.PageInfo{
		TotalElements: totalHits,
		TotalPages:    totalPages,
		Number:        pageNum,
		Size:          pageSize,
	}
	return result, nil

}
