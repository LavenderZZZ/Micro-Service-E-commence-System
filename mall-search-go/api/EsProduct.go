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

func (ctrl *EsProductController) ImportAllList(c *gin.Context) {
	count, _ := ctrl.Service.ImportAll()
	c.JSON(http.StatusOK, gin.H{"data": count})
}

func (ctrl *EsProductController) Delete(c *gin.Context) {
	idStr, _ := c.Params.Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil})
		return
	}
	ctrl.Service.Delete(id)
	c.JSON(http.StatusOK, gin.H{"data": nil})
}

func (ctrl *EsProductController) DeleteBatch(c *gin.Context) {
	var ids []int64
	err := c.BindJSON(&ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil})
		return
	}
	err = ctrl.Service.DeleteBatch(ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": nil})
}

func (ctrl *EsProductController) Create(c *gin.Context) {
	idStr, _ := c.Params.Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil})
		return
	}
	product, err := ctrl.Service.Create(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": product})
}

func (ctrl *EsProductController) SearchSimple(c *gin.Context) {
	keyword := c.Query("keyword")
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	result, err := ctrl.Service.SearchByNameOrSubTitleOrKeywords(keyword, pageNum, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

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
		c.JSON(http.StatusBadRequest, gin.H{"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})

}

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
		c.JSON(http.StatusBadRequest, gin.H{"data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (ctrl *EsProductController) SearchRelatedInfo(c *gin.Context) {
	keyword := c.Query("keyword")
	result, err := ctrl.Service.SearchRelated(keyword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
