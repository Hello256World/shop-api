package models

import (
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type Specification struct {
	ID         uint64     `gorm:"primaryKey"`
	Key        string     `gorm:"not null"`
	Value      string     `gorm:"not null"`
	ProductID  uint64     `gorm:"not null"`
	IsActive   *bool      `gorm:"default:true"`
	ModifiedAt *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt  time.Time  `gorm:"type:timestamp with time zone;default:now()"`
}

func (Specification) TableName() string {
	return "specification"
}

type SpecificationService struct {
	repo repository.Repository[Specification]
}

func NewSpecificationService(db *gorm.DB) *SpecificationService {
	return &SpecificationService{
		repo: repository.NewGenericRepository[Specification](db),
	}
}

func (s *SpecificationService) GetAll(id uint64) (*[]Specification, error) {
	var entities []Specification
	res := s.repo.GetQuery().Where("product_id = ?", id).Find(&entities)
	return &entities, res.Error
}

func (s *SpecificationService) Create(entity Specification) error {
	return s.repo.Create(&entity)
}

func (s *SpecificationService) GetById(id uint64) (*Specification, error) {
	return s.repo.GetByID(id)
}

func (s *SpecificationService) Update(entity *Specification) error {
	return s.repo.Update(entity)
}

func (s *SpecificationService) Delete(id uint64)error {
	return s.repo.Delete(id)
}