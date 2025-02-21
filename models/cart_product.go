package models

import (
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type CartProduct struct {
	ID         uint64     `gorm:"primaryKey"`
	CartID     uint64     `gorm:"not null"`
	ProductID  uint64     `gorm:"not null"`
	Quantity   int        `gorm:"type:int"`
	IsActive   *bool      `gorm:"default:true"`
	IsDelete   *bool      `gorm:"default:true"`
	ModifiedAt *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt  time.Time  `gorm:"type:timestamp with time zone;default:now()"`
}

func (CartProduct) TableName() string {
	return "cart_product"
}

type CartProductService struct {
	repo repository.Repository[CartProduct]
}

func NewCartProductService(db *gorm.DB) *CartProductService {
	return &CartProductService{
		repo: repository.NewGenericRepository[CartProduct](db),
	}
}

func (c *CartProductService) Count(customerID uint64) (*int64, error) {
	var count int64
	err := c.repo.GetQuery().Joins("JOIN cart ON cart.id = cart_product.cart_id").
	Where("cart.customer_id = ?", customerID).
	Count(&count).Error
	return &count, err
}
