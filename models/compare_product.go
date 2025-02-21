package models

import (
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type CompareProduct struct {
	ID         uint64     `gorm:"primaryKey"`
	ProductID  uint64     `gorm:"not null"`
	Name       string     `gorm:"not null"`
	Link       string     `gorm:"not null"`
	Price      float64    `gorm:"not null"`
	Image      string     `gorm:"not null"`
	IsActive   *bool      `gorm:"default:true"`
	IsDelete   *bool      `gorm:"default:false"`
	ModifiedAt *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt  time.Time  `gorm:"type:timestamp with time zone;default:now()"`
}

func (CompareProduct) TableName() string {
	return "compare_product"
}

type CompareProductService struct {
	repo repository.Repository[CompareProduct]
}

func NewCompareProductService(db *gorm.DB) *CompareProductService {
	return &CompareProductService{
		repo: repository.NewGenericRepository[CompareProduct](db),
	}
}

func (c *CompareProductService) GetAll(id uint64) (*[]CompareProduct, error) {
	var compareProducts []CompareProduct
	res := c.repo.GetQuery().Where("product_id = ?", id).Find(&compareProducts)
	return &compareProducts, res.Error
}

func (c *CompareProductService) GetById(id uint64) (*CompareProduct, error) {
	return c.repo.GetByID(id)
}

func (c *CompareProductService) Create(compareProduct CompareProduct) error {
	return c.repo.Create(&compareProduct)
}

func (c *CompareProductService) Update(compareProduct *CompareProduct) error {
	return c.repo.Update(compareProduct)
}

func (c *CompareProductService) Delete(id uint64) error {
	return c.repo.Delete(id)
}