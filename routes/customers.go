package routes

import (
	"net/http"
	"strings"

	"github.com/Hello256World/shop-api/models"
	"github.com/Hello256World/shop-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerHandler struct {
	customerService    *models.CustomerService
	cartProductService *models.CartProductService
}

func NewUserHandler(db *gorm.DB) *CustomerHandler {
	return &CustomerHandler{
		customerService:    models.NewCustomerService(db),
		cartProductService: models.NewCartProductService(db),
	}
}

func (u *CustomerHandler) getMe(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		return
	}

	if index := strings.Index(token, " "); index != -1 {
		token = token[index+1:]
	}

	res, err := utils.ValidateToken(token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "دوباره ثبت نام کنید"})
		return
	}

	id, ok := res["customerId"].(float64)

	if !ok {
		return
	}

	customer, err := u.customerService.GetById(uint64(id))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "دوباره ثبت نام کنید"})
		return
	}

	count, err := u.cartProductService.Count(customer.ID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت تعداد محصولات در سبد خرید"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": customer, "cart_product_count": count})
}
