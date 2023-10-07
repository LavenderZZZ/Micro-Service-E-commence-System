package store

import (
	"gorm.io/gorm"
	"mall-search-go/model"
)

type EsproductDao interface {
	GetAllProductList(id *int64) ([]model.EsProduct, error)
}

type EsProductDaoImpl struct {
	db *gorm.DB
}

func (e *EsProductDaoImpl) GetAllProductList(id *int64) ([]model.EsProduct, error) {
	var esProducts []model.EsProduct
	query := e.db.Preload("AttrValueList").Where("delete_status = ? AND publish_status = ?", 0, 1)

	if id != nil {
		query = query.Where("id = ?", id)
	}

	err := e.db.Find(&esProducts).Error
	if err != nil {
		return nil, err
	}

	return esProducts, err
}

func NewEsProductDao(db *gorm.DB) EsproductDao {
	return &EsProductDaoImpl{db: db}
}
