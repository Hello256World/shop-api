package routes

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Hello256World/shop-api/database"
	"github.com/Hello256World/shop-api/models"
	"github.com/Hello256World/shop-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	customerService *models.CustomerService
	cartService     *models.CartService
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		customerService: models.NewCustomerService(db),
		cartService:     models.NewCartService(db),
	}
}

func (a *AuthHandler) signup(c *gin.Context) {
	var customer models.Customer

	err := c.ShouldBind(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = a.customerService.Create(&customer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = a.cartService.Create(&models.Cart{
		CustomerID: customer.ID,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "شما با موفقیت ثبت نام کردید"})
}

func (a *AuthHandler) otp(c *gin.Context) {
	var request struct {
		Phone string `json:"phone" binding:"required"`
	}

	err := c.ShouldBind(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if pass := database.RDB.Get(context.Background(), request.Phone).Val(); pass != "" {
		// todo : send this pass to user again
		c.JSON(http.StatusOK, gin.H{"message": "رمز عبور موقت دوباره برای شما ارسال شد", "password": pass})
		return
	}

	customer, err := a.customerService.GetCustomerByPhone(request.Phone)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	rand.NewSource(time.Now().UnixNano())
	pass := rand.Intn(999999-111111) + 111111

	phoneKey := fmt.Sprintf("%v", customer.Phone)

	err = database.RDB.Set(context.Background(), phoneKey, pass, 6*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// todo : send password to the phone number

	c.JSON(http.StatusOK, gin.H{"message": "رمز عبور به تلفن همراه شما ارسال شد", "password": database.RDB.Get(context.Background(), phoneKey).Val()})
}

func (a *AuthHandler) signin(c *gin.Context) {
	var input struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" bining:"required"`
	}

	err := c.ShouldBind(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tempPass := database.RDB.Get(context.Background(), input.Phone).Val()

	if tempPass != input.Password {
		c.JSON(http.StatusBadRequest, gin.H{"message": "رمز عبور شما نامعتبر است"})
		return
	}

	customer, err := a.customerService.GetCustomerByPhone(input.Phone)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	token, err := utils.CreateToken("Customer", customer.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "شما با موفقیت وارد شدید", "token": token})
}
