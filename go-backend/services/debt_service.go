package services

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"github.com/sanda-bunescu/ExploRO/repositories"
	"gorm.io/gorm"
)

type DebtServiceInterface interface {
	CreateDebt(ctx context.Context, debt *models.Debt) error
	SoftDeleteDebtsByExpenseID(ctx context.Context, expenseId uint) error
	GetByUserIdAndGroupId(ginCtx *gin.Context, groupId uint) ([]*responses.DebtDetailResponse, error)
	Delete(ctx context.Context, debtId uint) error
	WithTx(tx *gorm.DB) *DebtService
}

type DebtService struct {
	debtRepo    repositories.DebtRepositoryInterface
	userRepo    repositories.UserRepositoryInterface
	expenseRepo repositories.ExpenseRepositoryInterface
}

func NewDebtService(debtRepo repositories.DebtRepositoryInterface, userRepo repositories.UserRepositoryInterface, expenseRepo repositories.ExpenseRepositoryInterface) *DebtService {
	return &DebtService{
		debtRepo:    debtRepo,
		userRepo:    userRepo,
		expenseRepo: expenseRepo,
	}
}

var _ DebtServiceInterface = (*DebtService)(nil)

func (ds *DebtService) CreateDebt(ctx context.Context, debt *models.Debt) error {
	err := ds.debtRepo.Create(ctx, debt)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DebtService) WithTx(tx *gorm.DB) *DebtService {
	return &DebtService{
		debtRepo:    ds.debtRepo.WithTx(tx),
		userRepo:    ds.userRepo,
		expenseRepo: ds.expenseRepo.WithTx(tx),
	}
}

func (ds *DebtService) SoftDeleteDebtsByExpenseID(ctx context.Context, expenseId uint) error {
	debts, err := ds.debtRepo.GetByExpenseID(ctx, expenseId)
	if err != nil {
		return err
	}
	return ds.debtRepo.SoftDeleteRange(ctx, debts)
}

func (ds *DebtService) GetByUserIdAndGroupId(ginCtx *gin.Context, groupId uint) ([]*responses.DebtDetailResponse, error) {
	firebaseUID, exists := ginCtx.Get("firebaseUID")
	if !exists {
		return nil, fmt.Errorf("unauthorized user: no firebaseUID in context")
	}
	//verify if account is deleted or not
	userRecord, err := ds.userRepo.GetByID(ginCtx, firebaseUID.(string))
	if err != nil {
		return nil, err
	}

	debtsResponse, err := ds.debtRepo.GetWithPayerByUserIdAndGroupId(ginCtx, firebaseUID.(string), groupId)
	if err != nil {
		return nil, err
	}
	for _, debt := range debtsResponse {
		debt.UserName = userRecord.Name
		payerUser, err := ds.userRepo.GetByID(ginCtx, debt.PayerId)
		if err != nil {
			return nil, err
		}
		debt.PayerName = payerUser.Name
	}
	return debtsResponse, nil
}

func (ds *DebtService) Delete(ctx context.Context, debtId uint) error {
	return ds.debtRepo.Transaction(func(tx *gorm.DB) error {
		debtserviceTx := ds.WithTx(tx)

		debt, err := debtserviceTx.debtRepo.GetByID(ctx, debtId)
		if err != nil {
			return err
		}

		// Soft delete the debt
		if err := debtserviceTx.debtRepo.SoftDelete(ctx, debt); err != nil {
			return err
		}

		// Now handle the expense after the debt is deleted
		debts, err := debtserviceTx.debtRepo.GetByExpenseID(ctx, debt.ExpenseId)
		if err != nil {
			return err
		}

		expense, err := debtserviceTx.expenseRepo.GetByID(ctx, debt.ExpenseId)
		if err != nil {
			return err
		}

		if len(debts) == 0 {
			// No more debts => delete the expense
			if err := debtserviceTx.expenseRepo.SoftDelete(ctx, expense); err != nil {
				return err
			}
		} else {
			// Still debts left => recalculate amount
			var newAmount float64
			for _, d := range debts {
				newAmount += d.AmountToPay
			}

			expense.Amount = newAmount
			if err := debtserviceTx.expenseRepo.Update(ctx, expense); err != nil {
				return err
			}
		}

		return nil
	})
}
