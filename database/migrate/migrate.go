package migrate

import (
	"log"

	"github.com/Hello256World/shop-api/database"
	"github.com/Hello256World/shop-api/models"
)

func Init() {
	if err := database.DB.Exec(`
    DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
            CREATE TYPE order_status AS ENUM ('confirmed', 'waiting_for_ipg', 'rejected', 'new', 'failed');
        END IF;
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'transaction_status') THEN
            CREATE TYPE transaction_status AS ENUM ('new', 'succeed', 'failed', 'roll-backed', 'in-progress', 'expired');
        END IF;
    END $$;
`).Error; err != nil {
		log.Fatalf("Failed to create enum types: %v", err)
		return
	}

	err := database.DB.AutoMigrate(&models.Customer{}, &models.Transaction{}, &models.Order{}, &models.Category{}, &models.Product{}, &models.OrderProduct{}, &models.Specification{}, &models.ImageProduct{}, &models.CompareProduct{}, &models.Cart{}, &models.CartProduct{}, &models.Admin{}, &models.SuperAdmin{})
	if err != nil {
		log.Fatal(err)
	}
}
