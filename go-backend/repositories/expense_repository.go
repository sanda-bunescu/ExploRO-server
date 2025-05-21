package repositories

import (
	"context"
	"github.com/sanda-bunescu/ExploRO/models"
	"gorm.io/gorm"
)

type ExpenseRepositoryInterface interface {
	BaseRepositoryInterface[models.Expense]
	GetAllWithDebtsByGroupId(groupId uint) ([]*models.Expense, error)
	WithTx(tx *gorm.DB) *ExpenseRepository
	Transaction(fn func(tx *gorm.DB) error) error
	GetAllByPayerId(ctx context.Context, payerId string) ([]*models.Expense, error)
}

type ExpenseRepository struct {
	BaseRepository[models.Expense]
}

func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{BaseRepository[models.Expense]{DB: db}}
}

var _ ExpenseRepositoryInterface = (*ExpenseRepository)(nil)

func (er *ExpenseRepository) GetAllWithDebtsByGroupId(groupId uint) ([]*models.Expense, error) {
	var expenses []*models.Expense
	err := er.DB.
		Preload("Debtors", "deleted_at IS NULL").
		Where("group_id = ? AND deleted_at IS NULL", groupId).
		Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func (er *ExpenseRepository) WithTx(tx *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{BaseRepository[models.Expense]{DB: tx}}
}

func (er *ExpenseRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return er.DB.Transaction(fn)
}

func (er *ExpenseRepository) GetAllByPayerId(ctx context.Context, payerId string) ([]*models.Expense, error) {
	var expenses []*models.Expense
	err := er.DB.WithContext(ctx).
		Where("payer_id = ? AND deleted_at IS NULL", payerId).
		Find(&expenses).Error
	if err != nil {
		return nil, err
	}
	return expenses, nil
}
