package routes

import (
	"fmt"
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

type CategoryHandler struct {
	categoryService *models.CategoryService
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{
		categoryService: models.NewCategoryService(db),
	}
}

func (ch *CategoryHandler) getAll(c *gin.Context) {
	name := c.Query("name")
	sortBy := c.Query("sortBy")
	order := c.Query("order")
	take := c.Query("take")
	skip := c.Query("skip")
	takeInt, err := strconv.Atoi(take)
	if err != nil {
		takeInt = 10
	}
	skipInt, err := strconv.Atoi(skip)
	if err != nil {
		skipInt = 0
	}

	categories, err := ch.categoryService.GetAll(name, sortBy, order, takeInt, skipInt)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

func (ch *CategoryHandler) getAllActive(c *gin.Context) {
	name := c.Query("name")
	sortBy := c.Query("sortBy")
	order := c.Query("order")
	take := c.Query("take")
	skip := c.Query("skip")
	takeInt, err := strconv.Atoi(take)
	if err != nil {
		takeInt = 10
	}
	skipInt, err := strconv.Atoi(skip)
	if err != nil {
		skipInt = 0
	}

	categories, err := ch.categoryService.GetAllActive(name, sortBy, order, takeInt, skipInt)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت دسته بندی"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"categories": categories})
}

func (ch *CategoryHandler) create(c *gin.Context) {
	var inputCategory struct {
		Name     string                `form:"name" binding:"required"`
		ParentID *uint64               `form:"parentId"`
		File     *multipart.FileHeader `form:"file" binding:"required"`
	}

	if err := c.ShouldBind(&inputCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	imageLocation, err := utils.AddImageToServer(c, "category", "thumbnail", inputCategory.File)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	category := models.Category{
		Name:     inputCategory.Name,
		Image:    fmt.Sprintf("https://%v", *imageLocation),
		ParentID: inputCategory.ParentID,
	}

	isOk := ch.categoryService.Create(category)

	if !isOk {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "مشکلی در ذخیره دسته بندی در دیتابیس پیش آمده"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "دسته بندی با موفقیت دخیره شد"})
}

func (ch *CategoryHandler) getById(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("categoryId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "خطا در شناسه دسته بندی"})
		return
	}

	category, err := ch.categoryService.GetById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

func (ch *CategoryHandler) update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("categoryId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه دسته بندی"})
		return
	}

	category, err := ch.categoryService.GetById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var inputCategory struct {
		Name     string                `form:"name" binding:"required"`
		File     *multipart.FileHeader `form:"file"`
		ParentID *uint64               `form:"parent_Id"`
		IsActive *bool                 `form:"is_active"`
	}

	if err := c.ShouldBind(&inputCategory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "مشکلی در دریافت اطلاعات پیش آمده", "error": err.Error()})
		return
	}

	if inputCategory.File != nil {
		imageLocation, err := utils.AddImageToServer(c, "category", "thumbnail", inputCategory.File)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		imageIndex := strings.Index(category.Image, "thumbnail")
		isOk := utils.DeleteImageOfServer("category", category.Image[imageIndex:])

		if !isOk {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "مشکلی در حذف عکس از سرور پیش آمده"})
			return
		}

		category.Image = fmt.Sprintf("https://%v", *imageLocation)
	}

	now := time.Now()
	category.ModifiedAt = &now
	category.Name = inputCategory.Name
	category.ParentID = inputCategory.ParentID
	if inputCategory.IsActive != nil {
		category.IsActive = inputCategory.IsActive
	}
	isOk := ch.categoryService.Update(*category)
	if !isOk {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "مشکلی در آپدیت دسته بندی پیش آمده"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "دسته بندی با موفقیت آپدیت شد"})
}

func (ch *CategoryHandler) delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("categoryId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "خطا در شناسه دسته بندی"})
		return
	}

	err = ch.categoryService.Delete(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "دسته بندی با موفقیت حذف شد"})
}
