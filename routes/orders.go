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

type OrderHandler struct {
	orderService *models.OrderService
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{
		orderService: models.NewOrderService(db),
	}
}

func (o *OrderHandler) getAll(c *gin.Context) {
	orders, err := o.orderService.GetAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت سفارش ها", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (o *OrderHandler) update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه سفارش"})
		return
	}

	order, err := o.orderService.GetById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var inputOrder struct {
		CustomerName    string             `form:"customer_name" binding:"required"`
		Phone           string             `form:"phone" binding:"required,phone"`
		Description     *string            `form:"description"`
		Weight          *float64           `form:"weight" binding:"required,gt=0"`
		DeliverMethod   string             `form:"deliver_method" binding:"required"`
		RejectionReason *string            `form:"rejection_reason"`
		TotalAmount     *float64           `form:"total_amount" binding:"required,gt=0"`
		DeliveryAddress string             `form:"delivery_address" binding:"required"`
		Status          models.OrderStatus `form:"status" binding:"required"`
	}

	if err := c.ShouldBind(&inputOrder); err != nil {
		getErrors := utils.FormValidation(err.Error(), map[string]string{"CustomerName": "نام مشتری", "Phone": "تلفن همراه", "Description": "توضیحات", "Weight": "وزن", "DeliverMethod": "روش ارسال", "RejectionReason": "دلیل رد شدن", "TotalAmount": "کل حساب", "DeliveryAddress": "آدرس ارسال", "Status": "وضعیت"})
		c.JSON(http.StatusNotAcceptable, gin.H{"message": getErrors})
		return
	}

	if inputOrder.Status == models.StatusRejected {
		if order.Status == models.StatusConfirmed {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "سفارش تایید شده است و نمی توان آن را رد کرد"})
			return
		} else if inputOrder.RejectionReason == nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "دلیل رد شدن سفارش را وارد کنید"})
			return
		}
	}
	now := time.Now()
	order.CustomerID = id
	order.CustomerName = inputOrder.CustomerName
	order.DeliverMethod = inputOrder.DeliverMethod
	order.DeliveryAddress = inputOrder.DeliveryAddress
	order.Description = inputOrder.Description
	order.ModifiedAt = &now
	order.Phone = inputOrder.Phone
	order.RejectionReason = inputOrder.RejectionReason
	order.Status = inputOrder.Status
	order.TotalAmount = *inputOrder.TotalAmount
	order.Weight = *inputOrder.Weight

	if err := o.orderService.Update(order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی سفارش"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "سفارش با موفقیت بروزرسانی شد"})
}
