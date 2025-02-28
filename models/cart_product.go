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

func (c *CartProductService) GetById(id uint64) (*CartProduct, error) {
	return c.repo.GetByID(id)
}

func (c *CartProductService) Update(entity *CartProduct) error {
	return c.repo.Update(entity)
}

func (c *CartProductService) Delete(id uint64) error {
	return c.repo.Delete(id)
}

func (c *CartProductService) DeleteAll(cartId uint64) error {
	return c.repo.GetQuery().Where("cart_id = ?", cartId).Delete(&CartProduct{}).Error
}

func (c *CartProductService) Create(entity *CartProduct) error {
	return c.repo.Create(entity)
}
