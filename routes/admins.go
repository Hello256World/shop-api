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

type AdminHandler struct {
	adminService *models.AdminService
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{
		adminService: models.NewAdminService(db),
	}
}

func (a *AdminHandler) create(c *gin.Context) {
	var inputAdmin struct {
		Username        string `form:"username" binding:"required"`
		Password        string `form:"password" binding:"required"`
		ConfirmPassword string `form:"confirm_password" binding:"required,eqfield=Password"`
		Phone           string `form:"phone" binding:"required,phone"`
	}

	if err := c.ShouldBind(&inputAdmin); err != nil {
		myError := utils.FormValidation(err.Error(), map[string]string{"Username": "نام کاربری", "ConfirmPassword": "تکرار رمز عبور", "Password": "رمز عبور", "Phone": "تلفن همراه"})
		c.JSON(http.StatusNotAcceptable, gin.H{"message": myError, "error": err.Error()})
		return
	}

	if _, err := a.adminService.GetByUsername(inputAdmin.Username); err == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "نام کاربری تکراری می باشد"})
		return
	}

	if _, err := a.adminService.GetByPhone(inputAdmin.Phone); err == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "شماره تلفن همراه تکراری می باشد"})
		return
	}

	hashPass, err := utils.HashPassword(inputAdmin.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در ساخت ادمین"})
		return
	}

	admin := models.Admin{
		Username: inputAdmin.Username,
		Password: hashPass,
		Phone:    inputAdmin.Phone,
	}

	if err := a.adminService.Create(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "ادمین با موفقیت ساخته شد"})
}

func (a *AdminHandler) getAll(c *gin.Context) {
	admins, err := a.adminService.GetAll()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در دریافت ادمین ها", "error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"admins": admins})
}

func (a *AdminHandler) update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه ادمین"})
		return
	}

	admin, err := a.adminService.GetById(id)

	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}

	var inputAdmin struct {
		Username        string `form:"username" binding:"required"`
		Password        string `form:"password" binding:"required"`
		ConfirmPassword string `form:"confirm_password" binding:"required,eqfield=Password"`
		Phone           string `form:"phone" binding:"required,phone"`
		IsActive        *bool  `form:"is_active" binding:"required"`
		IsDelete        *bool  `form:"is_delete" binding:"required"`
	}

	if err := c.ShouldBind(&inputAdmin); err != nil {
		getErrors := utils.FormValidation(err.Error(), map[string]string{"Username": "نام کاربری", "Password": "رمز عبور", "ConfirmPassword": "تکرار رمز عبور", "Phone": "تلفن همراه", "IsActive": "فعال/غیرفعال", "IsDelete": "حذف"})
		c.JSON(http.StatusNotAcceptable, gin.H{"message": getErrors})
		return
	}

	if _, err := a.adminService.GetByUsername(inputAdmin.Username); err == nil && admin.Username != inputAdmin.Username {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "ادمین دیگری با این نام کاربری وجود دارد"})
		return
	}

	if _, err := a.adminService.GetByPhone(inputAdmin.Phone); err == nil && admin.Phone != inputAdmin.Phone {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "ادمین دیگری با این تلفن همراه وجود دارد"})
		return
	}
	hashPass, err := utils.HashPassword(inputAdmin.ConfirmPassword)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در ساخت ادمین"})
		return
	}

	now := time.Now()
	admin.Password = hashPass
	admin.IsActive = inputAdmin.IsActive
	admin.IsDelete = inputAdmin.IsDelete
	admin.Phone = inputAdmin.Phone
	admin.Username = inputAdmin.Username
	admin.ModifiedAt = &now

	if err = a.adminService.Update(admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "ادمین با موفقیت آپدیت شد"})
}

func (a *AdminHandler) delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در شناسه ادمین"})
		return
	}

	admin, err := a.adminService.GetById(id)

	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": err.Error()})
		return
	}
	adTrue := true
	admin.IsDelete = &adTrue

	if err := a.adminService.Update(admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "خطا در حذف ادمین"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "ادمین با موفقیت حذف شد"})
}

func (a *AdminHandler) signin(c *gin.Context) {
	var inputAdmin struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	err := c.ShouldBind(&inputAdmin)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	admin, err := a.adminService.GetByUsername(inputAdmin.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !utils.CheckHashPass(inputAdmin.Password, admin.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "نام کاربری یا رمز عبور اشتباه است"})
		return
	}

	token, err := utils.CreateToken("Admin", admin.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "شما با موفقیت وارد شدید", "token": token})
}
