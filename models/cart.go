package models

import (
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type Cart struct {
	ID         uint64     `gorm:"primaryKey"`
	CustomerID uint64     `gorm:"unique;not null"`
	IsActive   *bool      `gorm:"default:true"`
	IsDelete   *bool      `gorm:"default:false"`
	ModifiedAt *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt  time.Time  `gorm:"type:timestamp with time zone;default:now()"`

	// Relations
	CartProducts []CartProduct `gorm:"foreignKey:CartID"`
}

func (Cart) TableName() string {
	return "cart"
}

type CartService struct {
	repo repository.Repository[Cart]
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{
		repo: repository.NewGenericRepository[Cart](db),
	}
}

func (c *CartService) Create(cart *Cart) error {
	return c.repo.Create(cart)
}

func (c *CartService) GetByCustomerId(customerId uint64) (*Cart, error) {
	var cart Cart
	res := c.repo.GetQuery().Where("customer_id = ? ", customerId).Preload("CartProducts").Find(&cart)

	if res.Error != nil {
		return nil, res.Error
	}

	return &cart, nil
}

func (c *CartService) GetAll(customerId uint64) (*Cart, error) {
	var cart Cart
	err := c.repo.GetQuery().Where("customer_id = ?", customerId).Preload("CartProducts").First(&cart).Error
	return &cart, err
}
