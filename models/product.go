package models

import (
	"time"

	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type Product struct {
	ID             uint64     `gorm:"primaryKey"`
	Name           string     `gorm:"not null"`
	Description    *string    `gorm:"column:description"`
	Price          float64    `gorm:"not null;type:float"`
	Stock          int        `gorm:"not null;type:int"`
	Thumbnail      string     `gorm:"not null;type:varchar"`
	CategoryID     uint64     `gorm:"not null;column:category_id"`
	ShipmentWeight float64    `gorm:"not null"`
	IsActive       *bool      `gorm:"default:true"`
	IsDelete       *bool      `gorm:"default:false"`
	ModifiedAt     *time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt      time.Time  `gorm:"type:timestamp with time zone;default:now()"`

	// Relations
	CartProducts    []CartProduct    `gorm:"foreignKey:ProductID" json:"-"`
	ImageProducts   []ImageProduct   `gorm:"foreignKey:ProductID" json:"-"`
	CompareProducts []CompareProduct `gorm:"foreignKey:ProductID" json:"-"`
	OrderProducts   []OrderProduct   `gorm:"foreignKey:ProductID" json:"-"`
	Specifications  []Specification  `gorm:"foreignKey:ProductID" json:"-"`
}

func (Product) TableName() string {
	return "product"
}

type ProductService struct {
	repo               repository.Repository[Product]
	cartProductService CartProductService
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{
		repo:               repository.NewGenericRepository[Product](db),
		cartProductService: *NewCartProductService(db),
	}
}

func (p *ProductService) GetAll(catId, productId uint64, minPrice, maxPrice float64, name, sortBy, order string, take, skip int) (*[]Product, error) {
	var products []Product
	query := p.repo.GetQuery().Where("category_id = ?", catId)

	if productId > 0 {
		query = query.Where("id = ?", productId).Limit(1).Find(&products)
	} else {
		if name != "" {
			query = query.Where("name LIKE ?", "%"+name+"%")
		}
		if minPrice > 0 {
			query = query.Where("price >= ?", minPrice)
		}
		if maxPrice > 0 {
			query = query.Where("price <= ?", maxPrice)
		}

		if sortBy != "" {
			if order == "desc" {
				query = query.Order(sortBy + " desc")
			} else {
				query = query.Order(sortBy + " asc")
			}
		}

		query = query.Offset(skip).Limit(take).Find(&products)
	}

	return &products, query.Error
}

func (p *ProductService) GetAllActive(productId, categoryId uint64, minPrice, maxPrice float64, name, sortBy, order string, take, skip int) (*[]Product, error) {
	var products []Product
	query := p.repo.GetQuery().
		Joins("JOIN category ON category.id = product.category_id").
		Where("product.is_active = ? AND product.is_delete = ? AND category.is_active = ? AND category.is_delete = ?", true, false, true, false)

	if productId > 0 {
		query = query.Where("product.id = ?", productId).Limit(1).Find(&products)
	} else {
		if categoryId > 0 {
			query = query.Where("product.category_id = ?", categoryId)
		}
		if name != "" {
			query = query.Where("product.name LIKE ?", "%"+name+"%")
		}
		if minPrice > 0 {
			query = query.Where("product.price >= ?", minPrice)
		}
		if maxPrice > 0 {
			query = query.Where("product.price <= ?", maxPrice)
		}

		if sortBy != "" {
			if order == "desc" {
				query = query.Order("product." + sortBy + " desc")
			} else {
				query = query.Order("product." + sortBy + " asc")
			}
		}
		query = query.Offset(skip).Limit(take).Find(&products)
	}

	return &products, query.Error
}

func (p *ProductService) Create(product Product) error {
	return p.repo.Create(&product)
}

func (p *ProductService) Update(product *Product) error {
	if product.Stock == 0 {
		p.cartProductService.repo.GetQuery().Where("product_id = ?", product.ID).Delete(&CartProduct{})
	}
	return p.repo.Update(product)
}

func (p *ProductService) GetById(id uint64) (*Product, error) {
	return p.repo.GetByID(id)
}

func (p *ProductService) Delete(id uint64) error {
	return p.repo.Delete(id)
}

func (p *ProductService) IsProductById(id uint64) bool {
	_, err := p.repo.GetByID(id)
	return err == nil
}

func (p *ProductService) GetProductsById(ids ...uint64) (*[]Product,error){
	var products []Product
	err := p.repo.GetQuery().Where("id IN ?",ids).Find(&products).Error
	return &products,err
}