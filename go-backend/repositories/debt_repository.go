package repositories

import (
	"context"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"gorm.io/gorm"
)

type DebtRepositoryInterface interface {
	BaseRepositoryInterface[models.Debt]
	GetByExpenseID(ctx context.Context, expenseId uint) ([]*models.Debt, error)
	GetWithPayerByUserIdAndGroupId(ctx context.Context, userId string, groupId uint) ([]*responses.DebtDetailResponse, error)
	WithTx(tx *gorm.DB) *DebtRepository
	Transaction(fn func(tx *gorm.DB) error) error
}

type DebtRepository struct {
	BaseRepository[models.Debt]
}

func NewDebtRepository(db *gorm.DB) *DebtRepository {
	return &DebtRepository{BaseRepository[models.Debt]{DB: db}}
}

var _ DebtRepositoryInterface = (*DebtRepository)(nil)

func (dr *DebtRepository) GetByExpenseID(ctx context.Context, expenseId uint) ([]*models.Debt, error) {
	var debts []*models.Debt
	err := dr.DB.WithContext(ctx).
		Where("expense_id = ? AND deleted_at IS NULL", expenseId).
		Find(&debts).Error

	if err != nil {
		return nil, err
	}
	return debts, nil
}

func (dr *DebtRepository) GetWithPayerByUserIdAndGroupId(ctx context.Context, userId string, groupId uint) ([]*responses.DebtDetailResponse, error) {
	var debts []*responses.DebtDetailResponse
	err := dr.DB.WithContext(ctx).
		Table("debts").
		Joins("JOIN expenses ON expenses.id = debts.expense_id").
		Where("debts.user_id = ? AND expenses.group_id = ? AND debts.deleted_at IS NULL AND expenses.deleted_at IS NULL", userId, groupId).
		Select("debts.id, debts.user_id, debts.amount_to_pay, expenses.payer_id").
		Find(&debts).Error

	if err != nil {
		return nil, err
	}
	for _, debt := range debts {
		debt.UserName = ""
		debt.PayerName = ""
	}

	return debts, nil
}

func (dr *DebtRepository) WithTx(tx *gorm.DB) *DebtRepository {
	return &DebtRepository{BaseRepository[models.Debt]{DB: tx}}
}

func (dr *DebtRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return dr.DB.Transaction(fn)
}
