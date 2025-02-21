package models

import (
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type ImageProduct struct {
	ID         uint64     `gorm:"primaryKey"`
	Image      string     `gorm:"not null"`
	Priority   int        `gorm:"not null"`
	ProductID  uint64     `gorm:"not null"`
	IsActive   *bool      `gorm:"default:true"`
	IsDelete   *bool      `gorm:"default:false"`
	ModifiedAt *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt  time.Time  `gorm:"type:timestamp with time zone;default:now()"`
}

func (ImageProduct) TableName() string {
	return "image_product"
}

type ImageProductService struct {
	repo repository.Repository[ImageProduct]
}

func NewImageProductService(db *gorm.DB) *ImageProductService {
	return &ImageProductService{
		repo: repository.NewGenericRepository[ImageProduct](db),
	}
}

func (i *ImageProductService) GetAll(id uint64) (*[]ImageProduct, error) {
	var images []ImageProduct
	res := i.repo.GetQuery().Where("product_id = ?", id).Find(&images)
	return &images, res.Error
}

func (i *ImageProductService) Create(image ImageProduct) error {
	return i.repo.Create(&image)
}

func (i *ImageProductService) GetById(id uint64) (*ImageProduct, error) {
	return i.repo.GetByID(id)
}

func (i *ImageProductService) Update(imageProduct *ImageProduct) error {
	return i.repo.Update(imageProduct)
}

func (i *ImageProductService) Delete(id uint64) error {
	return i.repo.Delete(id)
}