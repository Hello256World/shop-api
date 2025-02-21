package models

import (
	"errors"
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type SuperAdmin struct {
	ID         uint64     `gorm:"primaryKey"`
	Username   string     `gorm:"unique;not null"`
	Password   string     `gorm:"not null"`
	ModifiedAt *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt  time.Time  `gorm:"type:timestamp with time zone;default:now()"`
}

func (SuperAdmin) TableName() string {
	return "super_admin"
}

type SuperAdminService struct {
	repo repository.Repository[SuperAdmin]
}

func NewSuperAdminSerivce(db *gorm.DB) *SuperAdminService {
	return &SuperAdminService{
		repo: repository.NewGenericRepository[SuperAdmin](db),
	}
}

func (sa *SuperAdminService) GetByUsername(username string) (*SuperAdmin, error) {
	var superAdmin SuperAdmin
	res := sa.repo.GetQuery().Where("username = ?", username).First(&superAdmin)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("کاربری با این نام کاربری یافت نشد")
		}
		return nil, res.Error
	}

	return &superAdmin, nil
}
