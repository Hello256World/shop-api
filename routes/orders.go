package routes

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	transactionService  *models.TransactionService
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{
		orderService:        models.NewOrderService(db),
		addressService:      models.NewAddressService(db),
		cartService:         models.NewCartService(db),
		productService:      models.NewProductService(db),
		orderProductService: models.NewOrderProductService(db),
		cartProductService:  models.NewCartProductService(db),
		transactionService:  models.NewTransactionService(db),
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

	if inputOrder.Status == models.OrderStatusRejected {
		if order.Status == models.OrderStatusConfirmed {
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
	if err != nil || len(cart.CartProducts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "سبد خرید شما خالی می باشد"})
		return
	}

	var productsId []uint64
	var productsMap = make(map[int]int)
	for _, val := range cart.CartProducts {
		productsMap[int(val.ProductID)] = val.Quantity
		productsId = append(productsId, val.ProductID)
	}

	var weight, totalAmount float64
	products, err := o.productService.GetProductsById(productsId...)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت محصولات سبد خرید"})
		return
	}

	for _, val := range *products {
		weight += val.ShipmentWeight * float64(productsMap[int(val.ID)])
		totalAmount += val.Price * float64(productsMap[int(val.ID)])
	}

	userAgent := c.Request.UserAgent()

	var deviceType string

	if strings.Contains(userAgent, "Mobile") || strings.Contains(userAgent, "Android") || strings.Contains(userAgent, "iPhone") {
		deviceType = "mobile"
	} else {
		deviceType = "browser"
	}

	tx := o.orderService.BeginTransaction()
	defer tx.Rollback()

	transaction := models.Transaction{
		CustomerID: customerId,
		Type:       "default",
		Device:     deviceType,
		Status:     models.TransactionStatusNew,
		Amount:     totalAmount,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در ساخت تراکنش"})
		return
	}

	customerOrder := models.Order{
		AddressID:       address.ID,
		CustomerID:      customerId,
		TransactionID:   transaction.ID,
		CustomerName:    address.ReceiverName,
		Phone:           address.Phone,
		Description:     inputOrder.Description,
		DeliverMethod:   "post",
		DeliveryAddress: address.Address,
		Status:          models.OrderStatusWaitingForIPG,
		Weight:          weight,
		TotalAmount:     totalAmount,
	}

	if err := tx.Create(&customerOrder).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در ساخت سفارش", "error": err.Error()})
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
		orderProducts = append(orderProducts, models.OrderProduct{
			OrderID:   customerOrder.ID,
			Quantity:  val.Quantity,
			ProductID: val.ProductID,
			Price:     price,
		})
	}

	zarinPay, err := utils.NewZarinpal(os.Getenv("MERCHANT_ID"), true)

	if err != nil {
		transaction.Status = models.TransactionStatusFailed
		customerOrder.Status = models.OrderStatusFailed
		if err := tx.Save(&transaction).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی وضعیت تراکنش", "error": err.Error()})
			return
		}
		if err := tx.Save(&customerOrder).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی وضعیت سفارش", "error": err.Error()})
			return
		}
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در کامیت تراکنش", "error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "بخش اول : خطا در ایجاد درگاه پرداخت", "error": err.Error()})
		return
	}

	paymentURL, authority, statusCode, err := zarinPay.NewPaymentRequest(2000000000, fmt.Sprintf("http://localhost:8080/v1/public/orders/%v", customerOrder.ID), "test", "test@test.com", "09900994735")
	if err != nil {
		transaction.Status = models.TransactionStatusFailed
		customerOrder.Status = models.OrderStatusFailed
		if err := tx.Save(&transaction).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی وضعیت تراکنش", "error": err.Error()})
			return
		}
		if err := tx.Save(&customerOrder).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی وضعیت سفارش", "error": err.Error()})
			return
		}
		if err := tx.Commit().Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در کامیت تراکنش", "error": err.Error()})
			return
		}

		if statusCode == -3 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "مبلغ کل برای سیستم بانکی قابل قبول نیست"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "بخش دوم : خطا در ایجاد درگاه پرداخت", "error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در کامیت تراکنش", "error": err.Error()})
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

	c.JSON(http.StatusCreated, gin.H{"message": "سفارش با موفقیت ثبت شد", "URL": paymentURL, "authority": authority})
}

func (o *OrderHandler) paymentUpdate(c *gin.Context) {
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

	transaction, err := o.transactionService.GetById(order.TransactionID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var inputOrder struct {
		Authority string `form:"authority" binding:"required"`
		Status    string `form:"status" binding:"required"`
	}

	if err := c.ShouldBind(&inputOrder); err != nil {
		getErrors := utils.FormValidation(err.Error(), map[string]string{"Authority": "شناسه پرداخت"})
		c.JSON(http.StatusBadRequest, gin.H{"message": getErrors})
		return
	}

	zarinPay, err := utils.NewZarinpal(os.Getenv("MERCHANT_ID"), true)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در اطلاعات درگاه پرداخت", "error": err.Error()})
		return
	}

	verified, refID, statusCode, err := zarinPay.PaymentVerification(2000000000, inputOrder.Authority)
	if err != nil {
		if statusCode == 101 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "این پرداخت از قبل تایید شده است"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("status code is : %v", statusCode), "error": err.Error()})
		return
	}

	if verified {
		order.Status = models.OrderStatusNew
		transaction.Status = models.TransactionStatusSucceed

		if err := o.orderService.Update(order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی سفارش", "error": err.Error()})
			return
		}

		if err := o.transactionService.Update(transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی تراکنش", "error": err.Error()})
			return
		}
	} else {
		order.Status = models.OrderStatusFailed
		transaction.Status = models.TransactionStatusFailed

		if err := o.orderService.Update(order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی سفارش", "error": err.Error()})
			return
		}

		if err := o.transactionService.Update(transaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی تراکنش", "error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusAccepted, gin.H{"message": fmt.Sprintf("پرداخت شما تایید شد/نشد : %v", verified), "refID": refID})
}

func (o *OrderHandler) callBackUrl(c *gin.Context) {
	id := c.Param("id")

	authority := c.Query("Authority")
	status := c.Query("Status")

	html := `
		<!DOCTYPE html>
		<html lang="fa">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>صفحه دکمه</title>
			<script>
				window.onload = function() {
					const id = '` + id + `';
					const authority = '` + authority + `';
					const status = '` + status + `';

					const body = JSON.stringify({
						authority: authority,
						status: status
					});

					fetch('/v1/public/orders/' + id, {
						method: 'PUT',
						headers: {
							'Content-Type': 'application/json'
						},
						body: body
					})
					.then(response => response.text())
					.then(data => console.log(data))
					.catch(error => console.error('Error:', error));
				};
			</script>
			<style>
				body {
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
					background-color: #f0f0f0;
				}
				button {
					padding: 10px 20px;
					font-size: 16px;
				}
			</style>
		</head>
		<body>
			<button onclick="alert('دکمه کلیک شد!')">کلیک کن!</button>
		</body>
		</html>
		`

	c.Header("Content-Type", "text/html")
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
