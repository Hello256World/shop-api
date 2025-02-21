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

type CompareProductHandler struct {
	compareProductService models.CompareProductService
	productService        models.ProductService
}

func NewCompareProductHandler(db *gorm.DB) *CompareProductHandler {
	return &CompareProductHandler{
		compareProductService: *models.NewCompareProductService(db),
		productService:        *models.NewProductService(db),
	}
}

func (cp *CompareProductHandler) getAll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	compareProducts, err := cp.compareProductService.GetAll(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"compare_products": compareProducts, "message": "محصولات مشابه"})
}

func (cp *CompareProductHandler) getById(c *gin.Context) {
	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول مشابه"})
		return
	}

	compareProduct, err := cp.compareProductService.GetById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if compareProduct.ProductID != productId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "محصولی با این شناسه یافت نشد"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"compare_product": compareProduct})
}

func (cp *CompareProductHandler) create(c *gin.Context) {
	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	if !cp.productService.IsProductById(productId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "محصولی با این شناسه یافت نشد"})
		return
	}

	var inputCompareProduct struct {
		Name  string                `form:"name" binding:"required"`
		Link  string                `form:"link" binding:"required"`
		Price float64               `form:"price" binding:"required"`
		Image *multipart.FileHeader `form:"file" binding:"required"`
	}

	if err := c.ShouldBind(&inputCompareProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت دیتا", "error": err.Error()})
		return
	}

	imageAddress, err := utils.AddImageToServer(c, "productsimage", "compare images", inputCompareProduct.Image)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	compareProduct := models.CompareProduct{
		ProductID: productId,
		Name:      inputCompareProduct.Name,
		Link:      inputCompareProduct.Link,
		Price:     inputCompareProduct.Price,
		Image:     *imageAddress,
	}

	if err := cp.compareProductService.Create(compareProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در ذخیره محصول مشابه", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "محصول با موفقیت ذخیره شد"})
}

func (cp *CompareProductHandler) update(c *gin.Context) {
	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول مشابه"})
		return
	}

	compareProduct, err := cp.compareProductService.GetById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if compareProduct.ProductID != productId {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "این عملیات نامعتبر است"})
		return
	}

	var inputCompareProduct struct {
		Name     string                `form:"name" binding:"required"`
		Link     string                `form:"link" binding:"required"`
		Price    float64               `form:"price" binding:"required"`
		IsActive *bool                 `form:"is_active" binding:"required"`
		IsDelete *bool                 `form:"is_delete" binding:"required"`
		Image    *multipart.FileHeader `form:"file"`
	}

	if err := c.ShouldBind(&inputCompareProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت اطلاعات", "error": err.Error()})
		return
	}

	if inputCompareProduct.Image != nil {
		imageLocation, err := utils.AddImageToServer(c, "productsimage", "compare images", inputCompareProduct.Image)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		imageIndex := strings.Index(compareProduct.Image, "productsimage")
		utils.DeleteImageOfServer("productsimage", compareProduct.Image[imageIndex:])
		compareProduct.Image = *imageLocation
	}
	now := time.Now()
	compareProduct.IsActive = inputCompareProduct.IsActive
	compareProduct.IsDelete = inputCompareProduct.IsDelete
	compareProduct.Link = inputCompareProduct.Link
	compareProduct.Name = inputCompareProduct.Name
	compareProduct.Price = inputCompareProduct.Price
	compareProduct.ProductID = productId
	compareProduct.ModifiedAt = &now

	if err := cp.compareProductService.Update(compareProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در ذخیره محصول مشابه", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "محصول با موفقیت آپدیت شد"})
}

func (cp *CompareProductHandler) delete(c *gin.Context) {
	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول مشابه"})
		return
	}

	compareProduct, err := cp.compareProductService.GetById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if compareProduct.ProductID != productId {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "این عملیات نامعتبر است"})
		return
	}

	if err := cp.compareProductService.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در حذف محصول مشابه", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "محصول با موفقیت حذف شد"})
}
