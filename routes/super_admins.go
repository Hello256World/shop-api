package routes

import (
	"net/http"

	"github.com/Hello256World/shop-api/models"
	"github.com/Hello256World/shop-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SuperAdminHandler struct {
	superAdminService *models.SuperAdminService
}

func NewSuperAdminHandler(db *gorm.DB) *SuperAdminHandler {
	return &SuperAdminHandler{
		superAdminService: models.NewSuperAdminSerivce(db),
	}
}

func (sa *SuperAdminHandler) signin(c *gin.Context) {
	var inputSuperAdmin struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBind(&inputSuperAdmin)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	superAdmin, err := sa.superAdminService.GetByUsername(inputSuperAdmin.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if superAdmin.Password != inputSuperAdmin.Password {
		c.JSON(http.StatusBadRequest, gin.H{"message": "رمز عبور یا نام کاربری اشتباه است"})
		return
	}

	token, err := utils.CreateToken("SuperAdmin", superAdmin.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "با موفقیت وارد شدید", "token": token})
}
