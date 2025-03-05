package models

import (
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type OrderProduct struct {
	ID         uint64     `gorm:"primaryKey"`
	OrderID    uint64     `gorm:"not null"`
	ProductID  uint64     `gorm:"not null"`
	Quantity   int        `gorm:"not null"`
	Price      float64    `gorm:"not null"`
	ModifiedAt *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt  time.Time  `gorm:"type:timestamp with time zone;default:now()"`
}

func (OrderProduct) TableName() string {
	return "order_product"
}

type OrderProductService struct{
	repo repository.Repository[OrderProduct]
}

func NewOrderProductService(db *gorm.DB) *OrderProductService{
	return &OrderProductService{
		repo: repository.NewGenericRepository[OrderProduct](db),
	}
}

func (o *OrderProductService) CreateRange(orderProducts ...OrderProduct) error {
	return o.repo.GetQuery().Create(orderProducts).Error
}