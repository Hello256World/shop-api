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
	orderService        *models.OrderService
	addressService      *models.AddressService
	cartService         *models.CartService
	productService      *models.ProductService
	orderProductService *models.OrderProductService
	cartProductService  *models.CartProductService
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{
		orderService:        models.NewOrderService(db),
		addressService:      models.NewAddressService(db),
		cartService:         models.NewCartService(db),
		productService:      models.NewProductService(db),
		orderProductService: models.NewOrderProductService(db),
		cartProductService:  models.NewCartProductService(db),
	}
}

func (o *OrderHandler) getAll(c *gin.Context) {
	customerName := c.Query("customerName")
	sortBy := c.Query("sortBy")
	order := c.Query("order")
	take := c.Query("take")
	skip := c.Query("skip")
	customerId := c.Query("customerId")
	id := c.Query("id")
	var customerIdInt uint64
	if customerId != "" {
		if parsedId, err := strconv.ParseUint(customerId, 10, 64); err == nil {
			customerIdInt = parsedId
		}
	}
	var idInt uint64
	if id != "" {
		if parsedId, err := strconv.ParseUint(id, 10, 64); err == nil {
			idInt = parsedId
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

	orders, err := o.orderService.GetAll(idInt, customerIdInt, customerName, sortBy, order, takeInt, skipInt)
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

func (o *OrderHandler) getByCustomer(c *gin.Context) {
	id := c.Query("id")
	var idUint uint64
	if id != "" {
		newId, err := strconv.ParseUint(id, 10, 64)
		if err == nil {
			idUint = newId
		}
	}
	status := c.Query("status")
	var finalstate models.OrderStatus 
	if status != "" {
		finalstate = models.OrderStatus(status)
	}
	start := c.Query("start")
	end := c.Query("end")
	customerName := c.Query("customerName")
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

	orders, err := o.orderService.GetByCustomerId(idUint, c.GetUint64("customerId"), start, end, customerName, sortBy, order, finalstate, takeInt, skipInt)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت سفارش ها", "error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"orders": orders})
}

func (o *OrderHandler) create(c *gin.Context) {
	customerId := c.GetUint64("customerId")
	var inputOrder struct {
		AddressID   uint64  `form:"address_id" binding:"required"`
		Description *string `form:"description"`
	}

	if err := c.ShouldBind(&inputOrder); err != nil {
		getErrors := utils.FormValidation(err.Error(), map[string]string{"AddressID": "شناسه آدرس", "Description": "توضیحات"})
		c.JSON(http.StatusBadRequest, gin.H{"message": getErrors, "error": err.Error()})
		return
	}

	address, err := o.addressService.GetById(inputOrder.AddressID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت آدرس کاربر", "error": err.Error()})
		return
	}

	if address.CustomerID != customerId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت آدرس کاربر"})
		return
	}

	cart, err := o.cartService.GetByCustomerId(customerId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت محصولات سبد خرید"})
		return
	} else if len(cart.CartProducts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "سبد خرید شما خالی می باشد"})
		return
	}

	var productsId []uint64
	var productsMap = make(map[int]int)

	for _, val := range cart.CartProducts {
		productsMap[int(val.ProductID)] = val.Quantity
		productsId = append(productsId, val.ProductID)
	}

	var weight float64
	var totalAmount float64

	products, err := o.productService.GetProductsById(productsId...)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت محصولات سبد خرید"})
		return
	}

	for _, val := range *products {
		weight += val.ShipmentWeight * float64(productsMap[int(val.ID)])
		totalAmount += val.Price
	}

	customerOrder := models.Order{
		AddressID:       address.ID,
		CustomerID:      customerId,
		CustomerName:    address.ReceiverName,
		Phone:           address.Phone,
		Description:     inputOrder.Description,
		DeliverMethod:   "post",
		DeliveryAddress: address.Address,
		Status:          models.StatusWaitingForIPG,
		Weight:          weight,
		TotalAmount:     totalAmount,
	}

	tx := o.orderService.BeginTransaction()
	defer tx.Rollback()

	if err := tx.Create(&customerOrder).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در ساخت سفارش"})
		return
	}

	var orderProducts []models.OrderProduct

	for _, val := range cart.CartProducts {
		var price float64
		for _, value := range *products {
			if value.ID == val.ProductID {
				price = value.Price
			}
		}
		orderProduct := models.OrderProduct{
			OrderID:   customerOrder.ID,
			Quantity:  val.Quantity,
			ProductID: val.ProductID,
			Price:     price,
		}

		orderProducts = append(orderProducts, orderProduct)
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در ساخت سفارش"})
		return
	}

	if err := o.orderProductService.CreateRange(orderProducts...); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در ساخت سفارش محصولات"})
		return
	}

	if err := o.cartProductService.DeleteAll(cart.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در حذف کردن محصولات سبد خرید"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "سفارش با موفقیت ثبت شد", "URL": "https://hello.com"})
}

func (o *OrderHandler) customerUpdate(c *gin.Context) {
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

	if order.CustomerID != c.GetUint64("customerId") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "عملیات نا معتبر است"})
		return
	}

}
