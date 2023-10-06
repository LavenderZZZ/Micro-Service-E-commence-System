package model

import "time"

type PmsProduct struct {
	ID                         int64     `gorm:"column:id;primaryKey"`
	BrandID                    int64     `gorm:"column:brand_id"`
	ProductCategoryID          int64     `gorm:"column:product_category_id"`
	FeightTemplateID           int64     `gorm:"column:feight_template_id"`
	ProductAttributeCategoryID int64     `gorm:"column:product_attribute_category_id"`
	Name                       string    `gorm:"column:name"`
	Pic                        string    `gorm:"column:pic"`
	ProductSn                  string    `gorm:"column:product_sn"`
	DeleteStatus               int       `gorm:"column:delete_status"`
	PublishStatus              int       `gorm:"column:publish_status"`
	NewStatus                  int       `gorm:"column:new_status"`
	RecommandStatus            int       `gorm:"column:recommand_status"`
	VerifyStatus               int       `gorm:"column:verify_status"`
	Sort                       int       `gorm:"column:sort"`
	Sale                       int       `gorm:"column:sale"`
	Price                      float64   `gorm:"column:price"`
	PromotionPrice             float64   `gorm:"column:promotion_price"`
	GiftGrowth                 int       `gorm:"column:gift_growth"`
	GiftPoint                  int       `gorm:"column:gift_point"`
	UsePointLimit              int       `gorm:"column:use_point_limit"`
	SubTitle                   string    `gorm:"column:sub_title"`
	OriginalPrice              float64   `gorm:"column:original_price"`
	Stock                      int       `gorm:"column:stock"`
	LowStock                   int       `gorm:"column:low_stock"`
	Unit                       string    `gorm:"column:unit"`
	Weight                     float64   `gorm:"column:weight"`
	PreviewStatus              int       `gorm:"column:preview_status"`
	ServiceIds                 string    `gorm:"column:service_ids"`
	Keywords                   string    `gorm:"column:keywords"`
	Note                       string    `gorm:"column:note"`
	AlbumPics                  string    `gorm:"column:album_pics"`
	DetailTitle                string    `gorm:"column:detail_title"`
	PromotionStartTime         time.Time `gorm:"column:promotion_start_time"`
	PromotionEndTime           time.Time `gorm:"column:promotion_end_time"`
	PromotionPerLimit          int       `gorm:"column:promotion_per_limit"`
	PromotionType              int       `gorm:"column:promotion_type"`
	BrandName                  string    `gorm:"column:brand_name"`
	ProductCategoryName        string    `gorm:"column:product_category_name"`
	Description                string    `gorm:"column:description"`
	DetailDesc                 string    `gorm:"column:detail_desc"`
	DetailHtml                 string    `gorm:"column:detail_html"`
	DetailMobileHtml           string    `gorm:"column:detail_mobile_html"`
}
