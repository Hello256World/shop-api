package models

import (
	"github.com/Hello256World/shop-api/repository"
	"gorm.io/gorm"
)

type TransactionStatus string

func (s TransactionStatus) isValid() bool {
	switch s {
	case TransactionStatusFailed, TransactionStatusNew, TransactionStatusRollBack, TransactionStatusSucceed:
		return true
	default:
		return false
	}
}

const (
	TransactionStatusNew        TransactionStatus = "new"
	TransactionStatusSucceed    TransactionStatus = "succeed"
	TransactionStatusFailed     TransactionStatus = "failed"
	TransactionStatusRollBack   TransactionStatus = "roll-back"
	TransactionStatusInProgress TransactionStatus = "in-progress"
	TransactionStatusExpired    TransactionStatus = "expired"
)

type Transaction struct {
	ID                       uint64            `gorm:"primaryKey"`
	CustomerID               uint64            `gorm:"not null"`
	Device                   string            `gorm:"not null"`
	Type                     string            `gorm:"not null"`
	Status                   TransactionStatus `gorm:"type:transaction_status;not null"`
	RetrievalReferenceNumber *string           `gorm:"null;unique"`
	FailureCause             *string           `gorm:"null"`
	Amount                   float64           `gorm:"not null"`
	Description              *string           `gorm:"null;type:text"`

	// Relations
	Order Order `gorm:"foreignKey:TransactionID"`
}

func (Transaction) TableName() string {
	return "transaction"
}

type TransactionService struct {
	repo repository.Repository[Transaction]
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{
		repo: repository.NewGenericRepository[Transaction](db),
	}
}

func (t *TransactionService) GetById(id uint64) (*Transaction, error) {
	return t.repo.GetByID(id)
}

func (t *TransactionService) Update(entity *Transaction) error {
	return t.repo.Update(entity)
}
