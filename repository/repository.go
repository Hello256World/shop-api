package repository

import (
	"errors"

	"gorm.io/gorm"
)

type Repository[T any] interface {
	GetQuery() *gorm.DB
	Create(entity *T) error
	GetAll() (*[]T, error)
	GetByID(id uint64) (*T, error)
	Update(entity *T) error
	Delete(id uint64) error
}

type GenericRepository[T any] struct {
	db *gorm.DB
}

func NewGenericRepository[T any](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{db: db}
}

func (r *GenericRepository[T]) GetQuery() *gorm.DB {
	return r.db.Model(new(T))
}

func (r *GenericRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *GenericRepository[T]) GetAll() (*[]T, error) {
	var entities []T
	result := r.db.Find(&entities)
	return &entities, result.Error
}

func (r *GenericRepository[T]) GetByID(id uint64) (*T, error) {
	var entity T
	result := r.db.First(&entity, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("دیتا یافت نشد")
	}
	return &entity, result.Error
}

func (r *GenericRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *GenericRepository[T]) Delete(id uint64) error {
	return r.db.Delete(new(T), id).Error
}
