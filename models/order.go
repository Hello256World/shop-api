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
	StatusConfirmed     OrderStatus = "confirmed"
	StatusWaitingForIPG OrderStatus = "waiting_for_ipg"
	StatusRejected      OrderStatus = "rejected"
	StatusNew           OrderStatus = "new"
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

type Order struct {
	ID              uint64 `gorm:"primaryKey"`
	CustomerID      uint64 `gorm:"not null"`
	AddressID       uint64 `gorm:"not null"`
	CustomerName    string `gorm:"not null"`
	Phone           string `gorm:"not null"`
	Description     *string
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

func (o *OrderService) GetByCustomerId(customerId uint64, customerName, sortBy, order string, take, skip int) (*[]Order, error) {
	var orders []Order
	query := o.repo.GetQuery().Where("customer_id = ?", customerId)

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

	return &orders, query.Error
}

func (o *OrderService) Create(entity *Order) error {
	return o.repo.Create(entity)
}

func (o *OrderService) BeginTransaction() *gorm.DB {
	return o.repo.GetQuery().Begin()
}
