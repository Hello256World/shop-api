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

type ImageProductHandler struct {
	imageProductService *models.ImageProductService
	productService      *models.ProductService
}

func NewImageProductHandler(db *gorm.DB) *ImageProductHandler {
	return &ImageProductHandler{
		imageProductService: models.NewImageProductService(db),
		productService:      models.NewProductService(db),
	}
}

func (i *ImageProductHandler) getAll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	images, err := i.imageProductService.GetAll(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "عکس های محصول", "images": images})
}

func (i *ImageProductHandler) create(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	if !i.productService.IsProductById(id) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "محصولی یافت نشد"})
		return
	}

	var inputImageProducts struct {
		Image    *multipart.FileHeader `form:"file" binding:"required"`
		Priority *int                  `form:"priority" binding:"required,gte=0"`
	}

	if err := c.ShouldBind(&inputImageProducts); err != nil {
		getError := utils.FormValidation(err.Error(), map[string]string{"Image": "تصویر", "Priority": "اولویت"})
		c.JSON(http.StatusBadRequest, gin.H{"message": getError, "error": err.Error()})
		return
	}

	imageLocation, err := utils.AddImageToServer(c, "productsimage", "images", inputImageProducts.Image)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	imageProduct := models.ImageProduct{
		Image:     *imageLocation,
		Priority:  *inputImageProducts.Priority,
		ProductID: id,
	}

	if err := i.imageProductService.Create(imageProduct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "خطا در ذخیره عکس محصول در دیتا بیس", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "عکس محصول با موفقیت ثبت شد"})
}

func (i *ImageProductHandler) getById(c *gin.Context) {
	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	imageId, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه عکس محصول"})
		return
	}

	if !i.productService.IsProductById(productId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "محصولی با این شناسه یافت نشد"})
		return
	}

	imageProduct, err := i.imageProductService.GetById(imageId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if imageProduct.ProductID != productId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "شناسه عکس محصول غیر مجاز است"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"image_product": imageProduct})
}

func (i *ImageProductHandler) update(c *gin.Context) {
	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	imageId, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه عکس محصول"})
		return
	}

	if !i.productService.IsProductById(productId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "محصولی با این شناسه یافت نشد"})
		return
	}

	imageProduct, err := i.imageProductService.GetById(imageId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if imageProduct.ProductID != productId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "شناسه عکس محصول غیر مجاز است"})
		return
	}

	var inputImageProducts struct {
		Image    *multipart.FileHeader `form:"file"`
		Priority int                   `form:"priority" binding:"required,gte=0"`
		IsActive *bool                 `form:"is_active" binding:"required"`
		IsDelete *bool                 `form:"is_delete" binding:"required"`
	}

	if err := c.ShouldBind(&inputImageProducts); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت اطلاعات"})
		return
	}

	if inputImageProducts.Image != nil {
		imageLocation, err := utils.AddImageToServer(c, "productsimage", "images", inputImageProducts.Image)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		imageIndex := strings.Index(imageProduct.Image, "productsimage")
		utils.DeleteImageOfServer("productsimage", imageProduct.Image[imageIndex:])
		imageProduct.Image = *imageLocation
	}
	now := time.Now()
	imageProduct.IsActive = inputImageProducts.IsActive
	imageProduct.IsDelete = inputImageProducts.IsDelete
	imageProduct.Priority = inputImageProducts.Priority
	imageProduct.ProductID = productId
	imageProduct.ModifiedAt = &now

	if err := i.imageProductService.Update(imageProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی عکس محصول", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "عکس محصول با موفقیت آپدیت شد"})
}

func (i *ImageProductHandler) delete(c *gin.Context) {
	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	imageId, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه عکس محصول"})
		return
	}

	if !i.productService.IsProductById(productId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "محصولی با این شناسه یافت نشد"})
		return
	}

	imageProduct, err := i.imageProductService.GetById(imageId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if imageProduct.ProductID != productId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "شناسه عکس محصول غیر مجاز است"})
		return
	}

	if err := i.imageProductService.Delete(imageId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطادر حذف عکس محصول", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "محصول با موفقیت حذف شد"})
}
