package routes

import (
	"net/http"

	"github.com/Hello256World/shop-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandler struct {
	cartService *models.CartService
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{
		cartService: models.NewCartService(db),
	}
}

func (ch *CartHandler) getAll(c *gin.Context) {
	customerId := c.GetUint64("customerId")
	cart, err := ch.cartService.GetAll(customerId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت محصولات سبد خرید", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}