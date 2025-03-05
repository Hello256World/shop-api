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

type AddressHandler struct {
	addressService *models.AddressService
}

func NewAddressHandler(db *gorm.DB) *AddressHandler {
	return &AddressHandler{
		addressService: models.NewAddressService(db),
	}
}

func (a *AddressHandler) getAllActive(c *gin.Context) {
	addresses, err := a.addressService.GetAllActive(c.GetUint64("customerId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت آدرس ها", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"addresses": addresses})
}

func (a *AddressHandler) create(c *gin.Context) {
	var inputAddress struct {
		ReceiverName string `form:"receiver_name" binding:"required"`
		Address      string `form:"address" binding:"required"`
		Phone        string `form:"phone" binding:"required,phone"`
		NO           string `form:"no" binding:"required"`
		Unit         string `form:"unit" binding:"required"`
	}

	if err := c.ShouldBind(&inputAddress); err != nil {
		getErrors := utils.FormValidation(err.Error(), map[string]string{"ReceiverName": "نام دریافت کننده", "Address": "آدرس", "Phone": "تلفن همراه", "NO": "پلاک", "Unit": "واحد"})
		c.JSON(http.StatusBadRequest, gin.H{"message": getErrors, "error": err.Error()})
		return
	}

	address := models.Address{
		ReceiverName: inputAddress.ReceiverName,
		CustomerID:   c.GetUint64("customerId"),
		Address:      inputAddress.Address,
		Phone:        inputAddress.Phone,
		NO:           inputAddress.NO,
		Unit:         inputAddress.Unit,
	}

	if err := a.addressService.Create(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در ذخیره آدرس", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "آدرس با موفقیت ذخیره شد"})
}

func (a *AddressHandler) update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	customerId := c.GetUint64("customerId")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه آدرس"})
		return
	}

	address, err := a.addressService.GetById(id)

	if err != nil || address.CustomerID != customerId {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "آدرس موردنظر یافت نشد"})
		return
	}

	var inputAddress struct {
		ReceiverName string `form:"receiver_name" binding:"required"`
		Address      string `form:"address" binding:"required"`
		Phone        string `form:"phone" binding:"required,phone"`
		NO           string `form:"no" binding:"required"`
		Unit         string `form:"unit" binding:"required"`
		IsDelete     *bool  `form:"is_delete" binding:"required"`
	}

	if err := c.ShouldBind(&inputAddress); err != nil {
		getErrors := utils.FormValidation(err.Error(), map[string]string{"ReceiverName": "نام دریافت کننده", "Address": "آدرس", "Phone": "تلفن همراه", "NO": "پلاک", "Unit": "واحد", "IsDelete": "غیرفعال"})
		c.JSON(http.StatusBadRequest, gin.H{"message": getErrors, "error": err.Error()})
		return
	}

	now := time.Now()
	address.Address = inputAddress.Address
	address.ReceiverName = inputAddress.ReceiverName
	address.Phone = inputAddress.Phone
	address.NO = inputAddress.NO
	address.Unit = inputAddress.Unit
	address.IsDelete = inputAddress.IsDelete
	address.CustomerID = customerId
	address.ModifiedAt = &now

	if err := a.addressService.Update(address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی آدرس", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "آدرس با موفقیت بروزرسانی شد"})
}

func (a *AddressHandler) delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه آدرس"})
		return
	}

	address, err := a.addressService.GetById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "آدرس موردنظر یافت نشد"})
		return
	}

	if address.CustomerID != c.GetUint64("customerId") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "آدرس موردنظر یافت نشد"})
		return
	}

	if err := a.addressService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "خطا در حذف آدرس", "error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "آدرس با موفقیت حذف شد"})
}
