package models

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusConfirmed     OrderStatus = "confirmed"
	OrderStatusWaitingForIPG OrderStatus = "waiting_for_ipg"
	OrderStatusRejected      OrderStatus = "rejected"
	OrderStatusFailed        OrderStatus = "failed"
	OrderStatusNew           OrderStatus = "new"
)

func (s *OrderStatus) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}

	v, ok := value.(string)

	if !ok {
		return errors.New("failed to scan OrderStatus")
	}

	*s = OrderStatus(v)

	return nil
}

func (s OrderStatus) Value() (driver.Value, error) {
	return string(s), nil
}

func (s OrderStatus) isValid() bool {
	switch s {
	case OrderStatusConfirmed, OrderStatusFailed, OrderStatusRejected, OrderStatusNew, OrderStatusWaitingForIPG:
		return true
	default:
		return false
	}
}

type Order struct {
	ID              uint64 `gorm:"primaryKey"`
	CustomerID      uint64 `gorm:"not null"`
	AddressID       uint64 `gorm:"not null"`
	TransactionID   uint64 `gorm:"null"`
	CustomerName    string `gorm:"not null"`
	Phone           string `gorm:"not null"`
	Description     *string `gorm:"null;type:text"`
	Weight          float64 `gorm:"not null"`
	DeliverMethod   string  `gorm:"not null"`
	RejectionReason *string
	TotalAmount     float64     `gorm:"not null"`
	DeliveryAddress string      `gorm:"not null"`
	Status          OrderStatus `gorm:"type:order_status;not null"`
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

func (o *OrderService) GetAll(id, customerId uint64, customerName, sortBy, order string, take, skip int) (*[]Order, error) {
	var orders []Order

	query := o.repo.GetQuery()

	if id > 0 {
		query = query.Where("id = ?", id).Limit(1).Preload("OrderProducts").Find(&orders)
	} else {
		if customerId > 0 {
			query = query.Where("customer_id = ?", customerId)
		}
		if customerName != "" {
			query = query.Where("customer_name LIKE ?", "%"+customerName+"%")
		}
		if sortBy != "" {
			if order == "desc" {
				query = query.Order(sortBy + " desc")
			} else {
				query = query.Order(sortBy + " asc")
			}
		}

		query = query.Offset(skip).Limit(take).Preload("OrderProducts").Find(&orders)
	}

	return &orders, query.Error
}

func (o *OrderService) GetById(id uint64) (*Order, error) {
	return o.repo.GetByID(id)
}

func (o *OrderService) Update(order *Order) error {
	return o.repo.Update(order)
}

func (o *OrderService) GetByCustomerId(id, customerId uint64, start, end, customerName, sortBy, order string, status OrderStatus, take, skip int) (*[]Order, error) {
	var orders []Order
	query := o.repo.GetQuery().Where("customer_id = ?", customerId)

	if id > 0 {
		query = query.Where("id = ?", id).Limit(1).Preload("OrderProducts").Find(&orders)
	} else {
		if customerName != "" {
			query = query.Where("customer_name LIKE ?", "%"+customerName+"%")
		}
		if start != "" {
			startTime, err := time.Parse(time.RFC3339, start)
			if err == nil {
				query = query.Where("created_at >= ?", startTime)
			}
		}
		if end != "" {
			endTime, err := time.Parse(time.RFC3339, end)
			if err == nil {
				query = query.Where("created_at <= ?", endTime)
			}
		}
		if status.isValid() {
			query = query.Where("status = ?", status)
		}
		if sortBy != "" {
			if order == "desc" {
				query = query.Order(sortBy + " desc")
			} else {
				query = query.Order(sortBy + " asc")
			}
		}

		query = query.Offset(skip).Limit(take).Preload("OrderProducts").Find(&orders)
	}

	return &orders, query.Error
}

func (o *OrderService) Create(entity *Order) error {
	return o.repo.Create(entity)
}

func (o *OrderService) BeginTransaction() *gorm.DB {
	return o.repo.GetQuery().Begin()
}
