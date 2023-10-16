package service

import (
	//"mall-search-go/model"
	"mall-search-go/model"
)

type EsProductService interface {
	// Import all products from the database to ES
	ImportAll() (int, error)

	// Delete a product from ES
	Delete(id int64) error

	// Create a product in ES
	Create(id int64) (*model.EsProduct, error)

	// Search for products in ES
	DeleteBatch(ids []int64) error

	SearchByNameOrSubTitleOrKeywords(keyword string, pageNum, pageSize int) (model.Page, error)

	SearchByProductCategoryId(keyword string, brandId *int64, productCategoryId *int64, pageNum int, pageSize int, sort int) (model.Page, error)

	// recommend products based on product id
	Recommend(id int64, pageNum int, pageSize int) (model.Page, error)

	// SearchRelated products based on keyword
	SearchRelated(keyword string) (model.EsProductRelatedInfo, error)
}
