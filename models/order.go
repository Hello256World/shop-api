package models

import (
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusConfirmed OrderStatus = "confirmed"
	StatusWaiting   OrderStatus = "waiting"
	StatusRejected  OrderStatus = "rejected"
)

type Order struct {
	ID              uint64 `gorm:"primaryKey"`
	CustomerID      uint64 `gorm:"not null"`
	CustomerName    string `gorm:"not null"`
	Phone           string `gorm:"not null"`
	Description     *string
	Weight          float64 `gorm:"not null"`
	DeliverMethod   string  `gorm:"not null"`
	RejectionReason *string
	TotalAmount     float64     `gorm:"not null"`
	DeliveryAddress string      `gorm:"not null"`
	Status          OrderStatus `gorm:"not null"`
	ModifiedAt      *time.Time  `gorm:"type:timestamp with time zone"`
	CreatedAt       time.Time   `gorm:"type:timestamp with time zone;default:now()"`

	// Relations
	OrderProducts []OrderProduct `gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return "order"
}

type OrderService struct {
	repo repository.Repository[Order]
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		repo: repository.NewGenericRepository[Order](db),
	}
}

func (o *OrderService) GetAll() (*[]Order, error) {
	return o.repo.GetAll()
}

func (o *OrderService) GetById(id uint64) (*Order, error) {
	return o.repo.GetByID(id)
}

func (o *OrderService) Update(order *Order) error {
	return o.repo.Update(order)
}