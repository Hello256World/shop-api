package models

import (
	"errors"
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type Category struct {
	ID         uint64 `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	Image      string `gorm:"not null"`
	ParentID   *uint64
	IsActive   *bool `gorm:"default:true"`
	IsDelete   *bool `gorm:"default:false"`
	ModifiedAt *time.Time
	CreateAt   time.Time `gorm:"not null;default:now()"`
	// Relations
	SubCategories []Category `gorm:"foreignKey:ParentID" json:"-"`
	Products      []Product  `gorm:"foreignKey:CategoryID" json:"-"`
}

func (Category) TableName() string {
	return "category"
}

type CategoryService struct {
	repo repository.Repository[Category]
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{
		repo: repository.NewGenericRepository[Category](db),
	}
}

func (c *CategoryService) GetAll(name, sortBy, order string, take, skip int) (*[]Category, error) {
	var categories []Category
	query := c.repo.GetQuery()

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if sortBy != "" {
		if order == "desc" {
			query = query.Order(sortBy + " desc")
		} else {
			query = query.Order(sortBy + " asc")
		}
	}

	query = query.Offset(skip).Limit(take).Find(&categories)

	return &categories, query.Error
}

func (c *CategoryService) GetAllActive(name, sortBy, order string, take, skip int) (*[]Category, error) {
	var categories []Category

	query := c.repo.GetQuery().Where("is_active = ? AND is_delete = ?", true, false)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if sortBy != "" {
		if order == "desc" {
			query = query.Order(sortBy + " desc")
		} else {
			query = query.Order(sortBy + " asc")
		}
	}

	query = query.Offset(skip).Limit(take).Find(&categories)

	return &categories, query.Error
}

func (c *CategoryService) Create(category Category) bool {
	err := c.repo.Create(&category)
	return err == nil
}

func (c *CategoryService) GetById(id uint64) (*Category, error) {
	var entity Category
	res := c.repo.GetQuery().Preload("SubCategories").Preload("Products").First(&entity, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("دسته بندی با این شناسه یافت نشد")
	}
	return &entity, res.Error
}

func (c *CategoryService) Update(category Category) bool {
	err := c.repo.Update(&category)
	return err == nil
}

func (c *CategoryService) Delete(id uint64) error {
	// category, err := c.GetById(id)
	// if err != nil {
	// 	return err
	// }
	// if len(category.Products) > 0 || len(category.SubCategories) > 0 {
	// 	return errors.New("نمی توان این دسته بندی را حذف کرد زیرا با محصولات یا دسته بندی دیگری ارتباط دارد")
	// }

	return c.repo.Delete(id)
}
