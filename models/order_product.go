package models

import "time"

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
