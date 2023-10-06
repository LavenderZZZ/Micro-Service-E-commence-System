package model

import "math/big"

type EsProduct struct {
	Id                  int                     `json:"id"`
	ProductSn           string                  `json:"productSn"`
	BrandId             int                     `json:"brandId"`
	BrandName           string                  `json:"brandName"`
	ProductCategoryId   int                     `json:"productCategoryId"`
	ProductCategoryName string                  `json:"productCategoryName"`
	Pic                 string                  `json:"pic"`
	Name                string                  `json:"name"`
	SubTitle            string                  `json:"subTitle"`
	Keywords            string                  `json:"keywords"`
	Price               *big.Float              `json:"price"`
	Sale                int                     `json:"sale"`
	NewStatus           int                     `json:"newStatus"`
	RecommandStatus     int                     `json:"recommandStatus"`
	Stock               int                     `json:"stock"`
	PromotionType       int                     `json:"promotionType"`
	Sort                int                     `json:"sort"`
	AttrValueList       []ProductAttributeValue `json:"attrValueList"`
}

type ProductAttributeValue struct {
	Id                 int    `json:"id"`
	ProductAttributeID int    `json:"productAttributeId"`
	Value              string `json:"value"`
	Type               int    `json:"type"` // 0->specification; 1->parameter
	Name               string `json:"name"` // Attribute name
}
