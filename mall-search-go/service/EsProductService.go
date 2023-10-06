package service

import (
	//"mall-search-go/model"
	"mall-search-go/model"
)

type EsProductService interface {
	// Import all products from the database to ES
	ImportAll() (int, error)

	// Delete a product from ES
	Delete(id int) error

	// Delete multiple products from ES
	DeleteBatch(ids []int) error

	// Create a product in ES
	Create(id int) (model.EsProduct, error)

	// Create multiple products in ES
	CreateBatch(esProductList []model.EsProduct) (int, error)

	// Search for products in ES
	Search(keyword string, pageNum int, pageSize int) ([]model.EsProduct, int64, error)

	//// Search for products in ES by category
	//SearchByProductCategoryId(productCategoryId int, keyword string, pageNum int, pageSize int, sort int) ([]EsProduct, int64, error)
	//
	//// Search for products in ES by brand
	//SearchByBrandId(brandId int, keyword string, pageNum int, pageSize int, sort int) ([]EsProduct, int64, error)

	// recommend products based on product id
	Recommend(id int, pageNum int, pageSize int) ([]model.EsProduct, error)

	// SearchRelated products based on keyword
	SearchRelated(keyword string, pageNum int, pageSize int) ([]model.EsProduct, error)
}
