package migrate

import (
	"log"

	"github.com/Hello256World/shop-api/database"
	"github.com/Hello256World/shop-api/models"
)

func Init() {
	err := database.DB.AutoMigrate(&models.Customer{}, &models.Order{}, &models.Category{}, &models.Product{}, &models.OrderProduct{}, &models.Specification{}, &models.ImageProduct{}, &models.CompareProduct{}, &models.Cart{}, &models.CartProduct{}, &models.Admin{}, &models.SuperAdmin{})
	if err != nil {
		log.Fatal(err)
	}
}
