package model

type Page struct {
	Content  []EsProduct
	PageInfo PageInfo
}

type PageInfo struct {
	TotalPages    int
	TotalElements int
	Number        int
	Size          int
}

type EsProduct struct {
	ID                  int64 `gorm:"primaryKey"`
	ProductSn           string
	BrandId             int64
	BrandName           string
	ProductCategoryId   int64
	ProductCategoryName string
	Pic                 string
	Name                string
	SubTitle            string
	Price               float64
	Sale                int64
	NewStatus           int64
	RecommendStatus     int64
	Stock               int64
	PromotionType       int64
	Keywords            string
	Sort                int64
	AttrValueList       []EsProductAttributeValue `gorm:"foreignKey:ProductID"`
}

type EsProductAttributeValue struct {
	ID                 int64 `gorm:"primaryKey"`
	Value              string
	ProductAttributeID int64
	Type               string
	Name               string
	ProductID          int64
}

// EsProductRelatedInfo represents the product-related information for search results.
type EsProductRelatedInfo struct {
	BrandNames           []string      `json:"brandNames"`
	ProductCategoryNames []string      `json:"productCategoryNames"`
	ProductAttrs         []ProductAttr `json:"productAttrs"`
}

// ProductAttr represents a product attribute.
type ProductAttr struct {
	AttrId     int64    `json:"attrId"`
	AttrName   string   `json:"attrName"`
	AttrValues []string `json:"attrValues"`
}
