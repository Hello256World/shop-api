package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Hello256World/shop-api/models"
	"github.com/Hello256World/shop-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SpecificationHandler struct {
	specificationService *models.SpecificationService
	productService       *models.ProductService
}

func NewSpecificationHandler(db *gorm.DB) *SpecificationHandler {
	return &SpecificationHandler{
		specificationService: models.NewSpecificationService(db),
		productService:       models.NewProductService(db),
	}
}

func (s *SpecificationHandler) getAll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	if !s.productService.IsProductById(id) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "محصول مورد نظر یافت نشد"})
		return
	}

	specifications, err := s.specificationService.GetAll(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت مشخصات محصول"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"specifications": specifications})
}

func (s *SpecificationHandler) create(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	if !s.productService.IsProductById(id) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "محصول مورد نظر یافت نشد"})
		return
	}

	var inputSpecification struct {
		Key   string `form:"key" binding:"required"`
		Value string `form:"value" binding:"required"`
	}

	if err := c.ShouldBind(&inputSpecification); err != nil {
		getErrors := utils.FormValidation(err.Error(), map[string]string{"Key": "کلید", "Value": "مقدار"})
		c.JSON(http.StatusNotAcceptable, gin.H{"message": getErrors, "error": err.Error()})
		return
	}

	specification := models.Specification{
		Key:       inputSpecification.Key,
		Value:     inputSpecification.Value,
		ProductID: id,
	}

	if err := s.specificationService.Create(specification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در ذخیره مشخصات محصول"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "مشخصات محصول با موفقیت ثبت شد"})
}

func (s *SpecificationHandler) update(c *gin.Context) {
	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	if !s.productService.IsProductById(productId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "محصول مورد نظر یافت نشد"})
		return
	}

	specification, err := s.specificationService.GetById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if specification.ProductID != productId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "این عملیات غیر مجاز است"})
		return
	}

	var inputSpecification struct {
		Key      string `form:"key" binding:"required"`
		Value    string `form:"value" binding:"required"`
		IsActive *bool  `form:"is_active" binding:"required"`
	}

	if err := c.ShouldBind(&inputSpecification); err != nil {
		getErrors := utils.FormValidation(err.Error(), map[string]string{"Key": "کلید", "Value": "مقدار", "IsActive": "فعال / غیر فعال"})
		c.JSON(http.StatusBadRequest, gin.H{"message": getErrors})
		return
	}

	now := time.Now()
	specification.Value = inputSpecification.Value
	specification.ModifiedAt = &now
	specification.IsActive = inputSpecification.IsActive
	specification.Key = inputSpecification.Key
	specification.ProductID = productId

	if err := s.specificationService.Update(specification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی مشخصات محصول"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "مشخصات محصول با موفقیت بروزرسانی شد"})
}

func (s *SpecificationHandler) delete(c *gin.Context) {
	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	if !s.productService.IsProductById(productId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "محصول مورد نظر یافت نشد"})
		return
	}

	specification, err := s.specificationService.GetById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if specification.ProductID != productId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "این عملیات غیر مجاز است"})
		return
	}

	if err := s.specificationService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در خذف مشخصات محصول"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "مشخصات محصول با موفقیت حذف شد"})
}
