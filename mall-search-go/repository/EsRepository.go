package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	Recommend(id int64, product model.EsProduct, pageNum int, pageSize int) (model.Page, error)
	SearchRelated(keyword string) (model.EsProductRelatedInfo, error)
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

	//如果没有提供关键字（keyword为空），则使用match_all查询
	if keyword != "" {
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	} else {
		//使用multi_match查询搜索多个字段。这些字段的权重如下：
		//name: 权重为10
		//subTitle: 权重为5
		//keywords: 权重为2
		//所以，如果关键字在name字段中出现，它的重要性是在subTitle字段中出现的2倍，是在keywords字段中出现的5倍。
		query["query"] = map[string]interface{}{
			"function_score": map[string]interface{}{
				"query": map[string]interface{}{
					"multi_match": map[string]interface{}{
						"query":  keyword,
						"fields": []string{"name^10", "subTitle^5", "keywords"},
					},
				},
				"score_mode": "sum",
				"min_score":  2,
			},
		}
	}

	//如果提供了brandId或productCategoryId，则将它们添加为term查询来过滤结果
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
	//1: 根据id降序
	//2: 根据sale降序
	//3: 根据price升序
	//4: 根据price降序
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

func (repo *esProductRepositoryImpl) Recommend(id int64, product model.EsProduct, pageNum int, pageSize int) (model.Page, error) {
	var result model.Page

	query := map[string]interface{}{
		//must_not 确保原始商品不会被包含在结果中。
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must_not": map[string]interface{}{
					"term": map[string]interface{}{
						"id": id,
					},
				},
				//构建了一个加权的查询
				"should": []map[string]interface{}{
					{"match": map[string]interface{}{"name": map[string]interface{}{"query": product.Name, "boost": 8}}},
					{"match": map[string]interface{}{"subTitle": map[string]interface{}{"query": product.SubTitle, "boost": 2}}},
					{"match": map[string]interface{}{"keywords": map[string]interface{}{"query": product.Keywords, "boost": 2}}},
					{"match": map[string]interface{}{"brandId": map[string]interface{}{"query": product.BrandId, "boost": 5}}},
					{"match": map[string]interface{}{"productCategoryId": map[string]interface{}{"query": product.ProductCategoryId, "boost": 3}}},
				},
			},
		},
		"from": (pageNum - 1) * pageSize,
		"size": pageSize,
	}

	// Convert query to JSON and make the request
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return result, err
	}

	res, err := repo.client.Search(
		repo.client.Search.WithContext(context.Background()),
		repo.client.Search.WithIndex(repo.index),
		repo.client.Search.WithBody(&buf),
	)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	// Parse the response
	var searchResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return result, err
	}

	// Map the results
	var products []model.EsProduct
	for _, hit := range searchResult["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var product model.EsProduct
		mapstructure.Decode(hit.(map[string]interface{})["_source"], &product)
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

func (repo *esProductRepositoryImpl) SearchRelated(keyword string) (model.EsProductRelatedInfo, error) {
	var info model.EsProductRelatedInfo

	query := make(map[string]interface{})
	if keyword == "" {
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	} else {
		query["query"] = map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  keyword,
				"fields": []string{"name", "subTitle", "keywords"},
			},
		}
	}

	//aggregations
	query["aggs"] = map[string]interface{}{
		"brandNames": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "brandName",
			},
		},
		"productCategoryNames": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "productCategoryName",
			},
		},
		"allAttrValues": map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "attrValueList",
			},
			"aggs": map[string]interface{}{
				"productAttrs": map[string]interface{}{
					"terms": map[string]interface{}{
						"attrValueList.type": 1,
					},
				},
				"aggs": map[string]interface{}{
					"attrIds": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "attrValueList.productAttributeId",
						},
						"aggs": map[string]interface{}{
							"attrValues": map[string]interface{}{
								"terms": map[string]interface{}{
									"field": "attrValueList.value",
								},
							},
							"attrNames": map[string]interface{}{
								"terms": map[string]interface{}{
									"field": "attrValueList.name",
								},
							},
						},
					},
				},
			},
		},
	}
	// Convert query to JSON and make the request
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return info, err
	}

	req := esapi.SearchRequest{
		Index: []string{repo.index},
		Body:  &buf,
	}

	res, err := req.Do(context.Background(), repo.client)
	if err != nil {
		return info, err
	}
	defer res.Body.Close()

	// Parse the response
	var searchResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return info, err
	}

	return convertProductRelatedInfo(searchResult)

}

func convertProductRelatedInfo(response map[string]interface{}) (model.EsProductRelatedInfo, error) {
	var info model.EsProductRelatedInfo

	// Extract aggregations
	aggs, ok := response["aggregations"].(map[string]interface{})
	if !ok {
		return info, errors.New("Error parsing aggregations")
	}

	// Extract brand names
	if brandNamesAgg, exists := aggs["brandNames"].(map[string]interface{}); exists {
		for _, bucket := range brandNamesAgg["buckets"].([]interface{}) {
			brand := bucket.(map[string]interface{})
			brandName := brand["key"].(string)
			info.BrandNames = append(info.BrandNames, brandName)
		}
	}

	// Extract product category names
	if productCategoryNamesAgg, exists := aggs["productCategoryNames"].(map[string]interface{}); exists {
		for _, bucket := range productCategoryNamesAgg["buckets"].([]interface{}) {
			category := bucket.(map[string]interface{})
			categoryName := category["key"].(string)
			info.ProductCategoryNames = append(info.ProductCategoryNames, categoryName)
		}
	}

	// Extract product attributes
	if allAttrValuesAgg, exists := aggs["allAttrValues"].(map[string]interface{}); exists {
		productAttrsAgg := allAttrValuesAgg["productAttrs"].(map[string]interface{})
		attrIdsAgg := productAttrsAgg["attrIds"].(map[string]interface{})
		for _, bucket := range attrIdsAgg["buckets"].([]interface{}) {
			attr := model.ProductAttr{}
			attrBucket := bucket.(map[string]interface{})
			attr.AttrId = int64(attrBucket["key"].(float64))

			// Extract attribute values and names
			attrValuesAgg := attrBucket["attrValues"].(map[string]interface{})
			for _, valueBucket := range attrValuesAgg["buckets"].([]interface{}) {
				value := valueBucket.(map[string]interface{})["key"].(string)
				attr.AttrValues = append(attr.AttrValues, value)
			}

			attrNamesAgg := attrBucket["attrNames"].(map[string]interface{})
			if len(attrNamesAgg["buckets"].([]interface{})) > 0 {
				attr.AttrName = attrNamesAgg["buckets"].([]interface{})[0].(map[string]interface{})["key"].(string)
			}

			info.ProductAttrs = append(info.ProductAttrs, attr)
		}
	}

	return info, nil
}

//func (repo *esProductRepositoryImpl) Recommend(id int64, pageNum, pageSize int) (model.Page, error) {
//	var result model.Page
//
//
//
//}
