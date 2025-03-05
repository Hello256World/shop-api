package models

import (
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type Address struct {
	ID           uint64     `gorm:"primaryKey"`
	CustomerID   uint64     `gorm:"not null"`
	ReceiverName string     `gorm:"not null"`
	Address      string     `gorm:"not null"`
	Phone        string     `gorm:"not null"`
	NO           string     `gorm:"not null"`
	Unit         string     `gorm:"not null"`
	IsDelete     *bool      `gorm:"default:false"`
	ModifiedAt   *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt    time.Time  `gorm:"type:timestamp with time zone;default:now()"`

	//Relation
	Orders []Order `gorm:"foreignKey:AddressID"`
}

func (Address) TableName() string {
	return "address"
}

type AddressService struct {
	repo repository.Repository[Address]
}

func NewAddressService(db *gorm.DB) *AddressService {
	return &AddressService{
		repository.NewGenericRepository[Address](db),
	}
}

func (a *AddressService) GetById(id uint64) (*Address, error) {
	return a.repo.GetByID(id)
}

func (a *AddressService) GetAllActive(customerId uint64) (*[]Address, error) {
	var addresses []Address
	err := a.repo.GetQuery().Where("customer_id = ? AND is_delete = ?", customerId, false).Find(&addresses).Error
	return &addresses, err
}

func (a *AddressService) Create(entity *Address) error {
	return a.repo.Create(entity)
}

func (a *AddressService) Update(entity *Address) error {
	return a.repo.Update(entity)
}

func (a *AddressService) Delete(id uint64) error {
	return a.repo.Delete(id)
}
