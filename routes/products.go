package routes

import (
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Hello256World/shop-api/models"
	"github.com/Hello256World/shop-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	productService  *models.ProductService
	categoryService *models.CategoryService
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{
		productService:  models.NewProductService(db),
		categoryService: models.NewCategoryService(db),
	}
}

func (p *ProductHandler) getAll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("categoryId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	name := c.Query("name")
	minPriceStr := c.Query("minPrice")
	maxPriceStr := c.Query("maxPrice")
	sortBy := c.Query("sortBy")
	order := c.Query("order")
	take := c.Query("take")
	skip := c.Query("skip")
	productId := c.Query("productId")
	var productIdUint uint64
	if productId != "" {
		if parseId, err := strconv.ParseUint(productId, 10, 64); err == nil {
			productIdUint = parseId
		}
	}

	takeInt, err := strconv.Atoi(take)
	if err != nil {
		takeInt = 10
	}
	skipInt, err := strconv.Atoi(skip)
	if err != nil {
		skipInt = 0
	}

	var minPrice, maxPrice float64
	if minPriceStr != "" {
		if parsedMinPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			minPrice = parsedMinPrice
		}
	}

	if maxPriceStr != "" {
		if parsedMaxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			maxPrice = parsedMaxPrice
		}
	}

	products, err := p.productService.GetAll(id, productIdUint, minPrice, maxPrice, name, sortBy, order, takeInt, skipInt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (p *ProductHandler) getAllActive(c *gin.Context) {
	name := c.Query("name")
	minPriceStr := c.Query("minPrice")
	maxPriceStr := c.Query("maxPrice")
	sortBy := c.Query("sortBy")
	order := c.Query("order")
	take := c.Query("take")
	skip := c.Query("skip")
	productId := c.Query("productId")
	var productIdUint uint64
	if productId != "" {
		if parseId, err := strconv.ParseUint(productId, 10, 64); err == nil {
			productIdUint = parseId
		}
	}
	categoryId := c.Query("categoryId")
	var categoryIdUint uint64
	if categoryId != "" {
		if parseId, err := strconv.ParseUint(categoryId, 10, 64); err == nil {
			categoryIdUint = parseId
		}
	}
	takeInt, err := strconv.Atoi(take)
	if err != nil {
		takeInt = 10
	}
	skipInt, err := strconv.Atoi(skip)
	if err != nil {
		skipInt = 0
	}

	var minPrice, maxPrice float64
	if minPriceStr != "" {
		if parsedMinPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			minPrice = parsedMinPrice
		}
	}

	if maxPriceStr != "" {
		if parsedMaxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			maxPrice = parsedMaxPrice
		}
	}

	products, err := p.productService.GetAllActive(productIdUint, categoryIdUint, minPrice, maxPrice, name, sortBy, order, takeInt, skipInt)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت محصولات", "error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"products": products})
}

func (p *ProductHandler) create(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("categoryId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}
	_, err = p.categoryService.GetById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var inputProduct struct {
		Name        string                `json:"name" form:"name" binding:"required"`
		Price       float64               `json:"price" form:"price" binding:"required"`
		Stock       int                   `json:"stock" form:"stock" binding:"required"`
		Description *string               `json:"description" form:"description"`
		File        *multipart.FileHeader `json:"file" form:"file" binding:"required"`
	}

	if err := c.ShouldBind(&inputProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	imageLocation, err := utils.AddImageToServer(c, "productsimage", "thumbnail", inputProduct.File)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	product := models.Product{
		Description: inputProduct.Description,
		Name:        inputProduct.Name,
		Price:       inputProduct.Price,
		Stock:       inputProduct.Stock,
		Thumbnail:   *imageLocation,
		CategoryID:  id,
	}

	if err = p.productService.Create(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "محصول با موفقیت اضافه شد"})
}

func (p *ProductHandler) update(c *gin.Context) {
	catId, err := strconv.ParseUint(c.Param("categoryId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه دسته بندی"})
		return
	}
	_, err = p.categoryService.GetById(catId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	product, err := p.productService.GetById(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if product.CategoryID != catId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "این عملیات امکان پذیز نمی باشد"})
		return
	}

	var inputProduct struct {
		Name        string                `json:"name" form:"name" binding:"required"`
		Description *string               `json:"description" form:"description"`
		Price       float64               `json:"price" form:"price" binding:"required"`
		Stock       *int                  `json:"stock" form:"stock" binding:"required"`
		Thumbnail   *multipart.FileHeader `json:"file" form:"file"`
		IsActive    *bool                 `json:"is_active" form:"is_active"`
	}

	if err := c.ShouldBind(&inputProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "مشکلی در دریافت اطلاعات پیش آمده", "error": err.Error()})
		return
	}

	if inputProduct.Thumbnail != nil {
		imageLocation, err := utils.AddImageToServer(c, "productsimage", "thumbnail", inputProduct.Thumbnail)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		imageIndex := strings.Index(product.Thumbnail, "productsimage")
		go func() {
			if imageIndex != -1 {
				utils.DeleteImageOfServer("productsimage", product.Thumbnail[imageIndex:])
			}
		}()

		product.Thumbnail = *imageLocation
	}
	now := time.Now()
	product.Name = inputProduct.Name
	product.Description = inputProduct.Description
	product.Price = inputProduct.Price
	product.Stock = *inputProduct.Stock
	product.ModifiedAt = &now
	if inputProduct.IsActive != nil {
		product.IsActive = inputProduct.IsActive
	}

	if err = p.productService.Update(product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی محصول"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "محصول با موفقیت بروزرسانی شد"})
}

func (p *ProductHandler) delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	if err = p.productService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "خطا در خذف محصول", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "محصول با موفقیت حذف شد"})
}
