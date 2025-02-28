package migrate

import (
	"fmt"
	"log"

	"github.com/Hello256World/shop-api/database"
	"github.com/Hello256World/shop-api/models"
)

func Init() {
	if err := database.DB.Exec("DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN CREATE TYPE order_status AS ENUM ('confirmed', 'waiting', 'rejected'); END IF; END $$;").Error; err != nil {
		fmt.Println("Failed to create enum type:", err)
		return
	}
	err := database.DB.AutoMigrate(&models.Customer{}, &models.Order{}, &models.Category{}, &models.Product{}, &models.OrderProduct{}, &models.Specification{}, &models.ImageProduct{}, &models.CompareProduct{}, &models.Cart{}, &models.CartProduct{}, &models.Admin{}, &models.SuperAdmin{})
	if err != nil {
		log.Fatal(err)
	}
}
