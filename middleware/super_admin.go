package middleware

import (
	"net/http"
	"strings"

	"github.com/Hello256World/shop-api/utils"
	"github.com/gin-gonic/gin"
)

func SuperAdminAccess(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "باید ثبت نام کنید"})
		return
	}

	if index := strings.Index(token, " "); index != -1 {
		token = token[index+1:]
	}

	data, err := utils.ValidateToken(token)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		// c.Redirect(http.StatusUnauthorized,"/")
		return
	}

	role, ok := data["role"].(string)

	if !ok || role != "SuperAdmin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "شما دسترسی ندارید"})
		return
	}

	c.Next()
}
