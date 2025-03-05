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

type CartProductHandler struct {
	cartProductService *models.CartProductService
	cartService        *models.CartService
	productService     *models.ProductService
}

func NewCartProductHandler(db *gorm.DB) *CartProductHandler {
	return &CartProductHandler{
		models.NewCartProductService(db),
		models.NewCartService(db),
		models.NewProductService(db),
	}
}

func (ch *CartProductHandler) delete(c *gin.Context) {
	cartProductId, err := strconv.ParseUint(c.Param("cartProductId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه محصول"})
		return
	}

	cartProduct, err := ch.cartProductService.GetById(cartProductId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if cartProduct.Quantity > 1 {
		cartProduct.Quantity = cartProduct.Quantity - 1
		if err := ch.cartProductService.Update(cartProduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در کم کردن تعداد محصول"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "تعداد محصول با موفقیت کم شد"})
	} else {
		if err := ch.cartProductService.Delete(cartProductId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در کم کردن تعداد محصول"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "محصول با موفقیت حذف شد"})
	}

}

func (ch *CartProductHandler) deleteAll(c *gin.Context) {
	cartId, err := strconv.ParseUint(c.Param("cartId"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه سبد خرید"})
		return
	}

	if err := ch.cartProductService.DeleteAll(cartId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "حطا در حدف محصولات سبد خرید", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "محصولات سبد خرید با موفقیت حذف شد"})
}

func (ch *CartProductHandler) create(c *gin.Context) {
	var inputCartProduct struct {
		ProductID uint64 `form:"product_id" binding:"required"`
		Quantity  *int   `form:"quantity" binding:"required,gt=0"`
	}

	if err := c.ShouldBind(&inputCartProduct); err != nil {
		getErrors := utils.FormValidation(err.Error(), map[string]string{"ProductID": "شناسه محصول", "Quantity": "تعداد محصول"})
		c.JSON(http.StatusBadRequest, gin.H{"message": getErrors, "error": err.Error()})
		return
	}

	product, err := ch.productService.GetById(inputCartProduct.ProductID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !*product.IsActive || *product.IsDelete {
		c.JSON(http.StatusBadRequest, gin.H{"message": "محصول نامعتبر است"})
		return
	}

	if product.Stock < *inputCartProduct.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"message": "درخواست تعداد محصول بیشتر از موجودی می باشد"})
		return
	}

	cart, err := ch.cartService.GetByCustomerId(c.GetUint64("customerId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا دریافت سبد خرید", "error": err.Error()})
		return
	}

	for _, value := range cart.CartProducts {
		if value.ProductID == inputCartProduct.ProductID {
			now := time.Now()
			value.Quantity += *inputCartProduct.Quantity
			value.ModifiedAt = &now
			if err := ch.cartProductService.Update(&value); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در بروزرسانی سبد خرید", "error": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"message": "محصول سبد خرید با موفقیت بروزرسانی شد"})
			return
		}
	}

	cartProduct := models.CartProduct{
		ProductID: inputCartProduct.ProductID,
		Quantity:  *inputCartProduct.Quantity,
		CartID:    cart.ID,
	}

	if err := ch.cartProductService.Create(&cartProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در ذخیره در سبد خرید", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "محصول با موفقیت در سبد خرید اضافه شد"})
}
