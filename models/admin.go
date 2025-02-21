package models

import (
	"errors"
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type Admin struct {
	ID         uint64     `gorm:"primaryKey"`
	Username   string     `gorm:"unique;not null"`
	Password   string     `gorm:"not null"`
	Phone      string     `gorm:"unique;not null"`
	IsActive   *bool      `gorm:"default:true"`
	IsDelete   *bool      `gorm:"default:false"`
	ModifiedAt *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt  time.Time  `gorm:"type:timestamp with time zone;default:now()"`
}

func (Admin) TableName() string {
	return "admin"
}

type AdminService struct {
	repo repository.Repository[Admin]
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{
		repo: repository.NewGenericRepository[Admin](db),
	}
}

func (a *AdminService) Create(admin *Admin) error {
	err := a.repo.Create(admin)
	return err
}

func (a *AdminService) GetByUsername(username string) (*Admin, error) {
	var admin Admin
	res := a.repo.GetQuery().Where("username = ? ", username).First(&admin)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("کاربری با این نام کاربری یافت نشد")
	}

	return &admin, res.Error
}

func (a *AdminService) Update(admin *Admin) error {
	return a.repo.Update(admin)
}

func (a *AdminService) GetByPhone(phone string) (*Admin, error) {
	var admin Admin
	res := a.repo.GetQuery().Where("phone = ?", phone).First(&admin)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("کاربری با این شماره تلفن یافت نشد")
	}
	return &admin, res.Error
}

func (a *AdminService) GetAll() (*[]Admin, error) {
	return a.repo.GetAll()
}

func (a *AdminService) GetById(id uint64) (*Admin, error) {
	return a.repo.GetByID(id)
}
