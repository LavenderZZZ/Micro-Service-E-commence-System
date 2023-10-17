package api

import (
	"github.com/gin-gonic/gin"
	"mall-search-go/service"
	"net/http"
	"strconv"
)

type EsProductController struct {
	Service service.EsProductService
}

func NewEsProductController(service service.EsProductService) *EsProductController {
	return &EsProductController{Service: service}
}

func (ctrl *EsProductController) RegisterRoutes(router *gin.Engine) {

	esProductGroup := router.Group("/esProduct")

	esProductGroup.POST("/importAll", ctrl.ImportAllList)
	esProductGroup.GET("/delete/:id", ctrl.Delete)
	esProductGroup.POST("/delete/batch", ctrl.DeleteBatch)
	esProductGroup.POST("/create/:id", ctrl.Create)
	esProductGroup.GET("/search/simple", ctrl.SearchSimple)
	esProductGroup.GET("/search", ctrl.Search)
	esProductGroup.GET("/recommend/:id", ctrl.Recommend)
	esProductGroup.GET("/search/relate", ctrl.SearchRelatedInfo)
}

// @Summary Import all products to Elasticsearch
// @Description Import all products from the database to Elasticsearch
// @Tags esProduct
// @Accept  json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /esProduct/importAll [post]
func (ctrl *EsProductController) ImportAllList(c *gin.Context) {
	count, err := ctrl.Service.ImportAll()
	if err != nil {
		res := Failed("Failed to import" + err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusOK, Success(count))
}

// @Summary Delete a product in Elasticsearch by ID
// @Description Delete a specific product in Elasticsearch using its ID
// @Tags esProduct
// @Accept  json
// @Produce json
// @Param  id   path   int64  true  "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /esProduct/delete/{id} [get]
func (ctrl *EsProductController) Delete(c *gin.Context) {
	idStr, _ := c.Params.Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil})
		return
	}
	err = ctrl.Service.Delete(id)
	if err != nil {
		res := Failed("Failed to delete" + err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusOK, Success(nil))
}

// @Summary Batch delete products in Elasticsearch
// @Description Delete multiple products in Elasticsearch using their IDs
// @Tags esProduct
// @Accept  json
// @Produce json
// @Param  ids   body  []int64  true  "Array of Product IDs"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /esProduct/delete/batch [post]
func (ctrl *EsProductController) DeleteBatch(c *gin.Context) {
	var ids []int64
	err := c.BindJSON(&ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil})
		return
	}
	err = ctrl.Service.DeleteBatch(ids)
	if err != nil {
		res := Failed("Failed to delete" + err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": nil})
}

// @Summary Create a product in Elasticsearch by ID
// @Description Add a specific product to Elasticsearch using its database ID
// @Tags esProduct
// @Accept  json
// @Produce json
// @Param  id   path   int64  true  "Database Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /esProduct/create/{id} [post]
func (ctrl *EsProductController) Create(c *gin.Context) {
	idStr, _ := c.Params.Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		res := Failed("Invalid ID")
		c.JSON(http.StatusBadRequest, res)
		return
	}
	product, err := ctrl.Service.Create(id)
	if err != nil {
		res := Failed("Failed to create product" + err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusOK, Success(product))
}

// @Summary Simple search in Elasticsearch
// @Description Search products by name, subtitle, or keywords
// @Tags esProduct
// @Accept  json
// @Produce json
// @Param  keyword   query   string  true  "Keyword for search"
// @Param  pageNum   query   int     false "Page number"
// @Param  pageSize  query   int     false "Number of items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /esProduct/search/simple [get]
func (ctrl *EsProductController) SearchSimple(c *gin.Context) {
	keyword := c.Query("keyword")
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	result, err := ctrl.Service.SearchByNameOrSubTitleOrKeywords(keyword, pageNum, pageSize)
	if err != nil {
		res := Failed("Failed to search" + err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusOK, Success(result))

}

// @Summary Detailed search in Elasticsearch
// @Description Search products by multiple criteria
// @Tags esProduct
// @Accept  json
// @Produce json
// @Param  keyword              query   string  false "Keyword for search"
// @Param  brandId              query   int64   false "Brand ID"
// @Param  productCategoryId    query   int64   false "Product Category ID"
// @Param  pageNum              query   int     false "Page number"
// @Param  pageSize             query   int     false "Number of items per page"
// @Param  sort                 query   int     false "Sort order"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /esProduct/search [get]
func (ctrl *EsProductController) Search(c *gin.Context) {
	keyword := c.DefaultQuery("keyword", "")
	brandIdStr := c.DefaultQuery("brandId", "0")
	brandId, _ := strconv.ParseInt(brandIdStr, 10, 64)
	productCategoryIdStr := c.DefaultQuery("productCategoryId", "0")
	productCategoryId, _ := strconv.ParseInt(productCategoryIdStr, 10, 64)
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	sort, _ := strconv.Atoi(c.DefaultQuery("sort", "0"))

	result, err := ctrl.Service.SearchByProductCategoryId(keyword, &brandId, &productCategoryId, pageNum, pageSize, sort)
	if err != nil {
		res := Failed("Failed to search" + err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusOK, Success(result))

}

// @Summary Recommend products
// @Description Recommend products based on a specific product ID
// @Tags esProduct
// @Accept  json
// @Produce json
// @Param  id       path   int64  true  "Product ID"
// @Param  pageNum  query   int     false "Page number"
// @Param  pageSize query   int     false "Number of items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /esProduct/recommend/{id} [get]
func (ctrl *EsProductController) Recommend(c *gin.Context) {
	idStr, _ := c.Params.Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil})
		return
	}
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "0"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
	result, err := ctrl.Service.Recommend(id, pageNum, pageSize)
	if err != nil {
		res := Failed("Failed to Recommend" + err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusOK, Success(result))
}

// @Summary Search related product info
// @Description Search for product-related information based on a keyword
// @Tags esProduct
// @Accept  json
// @Produce json
// @Param  keyword   query   string  true  "Keyword for search"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /esProduct/search/relate [get]
func (ctrl *EsProductController) SearchRelatedInfo(c *gin.Context) {
	keyword := c.Query("keyword")
	result, err := ctrl.Service.SearchRelated(keyword)
	if err != nil {
		res := Failed("Failed to search" + err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusOK, Success(result))
}
