package models

import (
	"errors"
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type Customer struct {
	ID         uint64     `gorm:"primaryKey"`
	Fullname   string     `gorm:"not null" binding:"required"`
	Email      *string    `gorm:"unique"`
	Phone      string     `gorm:"unique;not null" binding:"required"`
	Birthday   *time.Time `gorm:"type:timestamp"`
	Gender     *Gender    `gorm:"column:gender"`
	ModifiedAt *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt  time.Time  `gorm:"type:timestamp with time zone;default:now()"`

	// Relations
	Cart      Cart      `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;" json:"-"`
	Orders    []Order   `gorm:"foreignKey:CustomerID" json:"-"`
	Addresses []Address `gorm:"foreignKey:CustomerID" json:"-"`
}

func (Customer) TableName() string {
	return "customer"
}

type CustomerService struct {
	repo repository.Repository[Customer]
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{
		repo: repository.NewGenericRepository[Customer](db),
	}
}

func (cs *CustomerService) Create(c *Customer) error {
	return cs.repo.Create(c)
}

func (cs *CustomerService) GetByPhone(phone string) (*Customer, error) {
	var customer Customer
	result := cs.repo.GetQuery().Where("phone = ?", phone).First(&customer)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("کاربر مورد نظر پیدا نشد")
	}
	return &customer, result.Error
}

func (cs *CustomerService) GetById(id uint64) (*Customer, error) {
	return cs.repo.GetByID(id)
}
