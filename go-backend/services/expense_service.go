package services

import (
	"context"
	"github.com/sanda-bunescu/ExploRO/models"
	"github.com/sanda-bunescu/ExploRO/models/requests"
	"github.com/sanda-bunescu/ExploRO/models/responses"
	"github.com/sanda-bunescu/ExploRO/repositories"
	"gorm.io/gorm"
)

type ExpenseServiceInterface interface {
	GetAllExpensesByGroupId(ctx context.Context, groupId uint) ([]responses.ExpenseResponse, error)
	CreateExpenseWithDebts(ctx context.Context, req requests.NewExpenseRequest) error
	EditExpenseWithDebts(ctx context.Context, req requests.EditExpenseRequest) error
	SoftDeleteExpenseByID(ctx context.Context, expenseId uint) error
	SoftDeleteExpenseByPayerId(ctx context.Context, payerId string) error
}

type ExpenseService struct {
	expenseRepo repositories.ExpenseRepositoryInterface
	userRepo    repositories.UserRepositoryInterface
	debtService DebtServiceInterface
}

func NewExpenseService(repo repositories.ExpenseRepositoryInterface, userRepo repositories.UserRepositoryInterface, debtService DebtServiceInterface) *ExpenseService {
	return &ExpenseService{
		expenseRepo: repo,
		userRepo:    userRepo,
		debtService: debtService,
	}
}

var _ ExpenseServiceInterface = (*ExpenseService)(nil)

func (es *ExpenseService) GetAllExpensesByGroupId(ctx context.Context, groupId uint) ([]responses.ExpenseResponse, error) {

	expenses, err := es.expenseRepo.GetAllWithDebtsByGroupId(groupId)
	if err != nil {
		return nil, err
	}

	var expensesResponse []responses.ExpenseResponse
	for _, expense := range expenses {
		payerUser, err := es.userRepo.GetByID(ctx, expense.PayerId)
		if err != nil {
			return nil, err
		}

		debtorsResponse := []responses.DebtResponse{}

		for _, debtor := range expense.Debtors {
			debtorIdentity, err := es.userRepo.GetByID(ctx, debtor.UserId)
			if err != nil {
				return nil, err
			}
			debtorsResponse = append(debtorsResponse, responses.DebtResponse{
				Id:          debtor.Id,
				UserId:      debtor.UserId,
				UserName:    debtorIdentity.Email,
				AmountToPay: debtor.AmountToPay,
			})
		}

		expensesResponse = append(expensesResponse, responses.ExpenseResponse{
			Id:            expense.Id,
			Name:          expense.Name,
			GroupId:       expense.GroupId,
			PayerUserName: payerUser.Name,
			Date:          expense.Date,
			Amount:        expense.Amount,
			Description:   expense.Description,
			Type:          expense.Type,
			Debtors:       debtorsResponse,
		})
	}
	return expensesResponse, nil
}

func (es *ExpenseService) CreateExpenseWithDebts(ctx context.Context, req requests.NewExpenseRequest) error {
	return es.expenseRepo.Transaction(func(tx *gorm.DB) error {
		expenseRepoTx := es.expenseRepo.WithTx(tx)
		debtServiceTx := es.debtService.WithTx(tx)

		expense := &models.Expense{
			Name:        req.Name,
			Amount:      req.Amount,
			Type:        req.Type,
			Description: req.Description,
			Date:        req.Date,
			GroupId:     req.GroupId,
			PayerId:     req.PayerId,
		}

		err := expenseRepoTx.Create(ctx, expense)
		if err != nil {
			return err
		}
		for _, debt := range req.Debtors {
			debtModel := &models.Debt{
				ExpenseId:   expense.Id,
				UserId:      debt.UserId,
				AmountToPay: debt.AmountToPay,
			}
			err = debtServiceTx.CreateDebt(ctx, debtModel)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (es *ExpenseService) EditExpenseWithDebts(ctx context.Context, req requests.EditExpenseRequest) error {
	return es.expenseRepo.Transaction(func(tx *gorm.DB) error {
		expenseRepoTx := es.expenseRepo.WithTx(tx)
		debtServiceTx := es.debtService.WithTx(tx)

		existingExpense, err := expenseRepoTx.GetByID(ctx, req.Id)
		if err != nil {
			return err
		}

		existingExpense.Name = req.Name
		existingExpense.Amount = req.Amount
		existingExpense.Type = req.Type
		existingExpense.Description = req.Description
		existingExpense.Date = req.Date

		if err := expenseRepoTx.Update(ctx, existingExpense); err != nil {
			return err
		}

		if err := debtServiceTx.SoftDeleteDebtsByExpenseID(ctx, req.Id); err != nil {
			return err
		}

		for _, debtor := range req.Debtors {
			debt := &models.Debt{
				ExpenseId:   req.Id,
				UserId:      debtor.UserId,
				AmountToPay: debtor.AmountToPay,
			}
			if err := debtServiceTx.CreateDebt(ctx, debt); err != nil {
				return err
			}
		}

		return nil
	})
}

func (es *ExpenseService) SoftDeleteExpenseByID(ctx context.Context, expenseId uint) error {
	return es.expenseRepo.Transaction(func(tx *gorm.DB) error {
		expenseRepoTx := es.expenseRepo.WithTx(tx)
		debtServiceTx := es.debtService.WithTx(tx)

		expense, err := expenseRepoTx.GetByID(ctx, expenseId)
		if err != nil {
			return err
		}

		err = expenseRepoTx.SoftDelete(ctx, expense)
		if err != nil {
			return err
		}

		err = debtServiceTx.SoftDeleteDebtsByExpenseID(ctx, expense.Id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (es *ExpenseService) SoftDeleteExpenseByPayerId(ctx context.Context, payerId string) error {
	return es.expenseRepo.Transaction(func(tx *gorm.DB) error {
		expenseRepoTx := es.expenseRepo.WithTx(tx)
		debtServiceTx := es.debtService.WithTx(tx)

		expenses, err := expenseRepoTx.GetAllByPayerId(ctx, payerId)
		if err != nil {
			return err
		}

		for _, expense := range expenses {
			err := expenseRepoTx.SoftDelete(ctx, expense)
			if err != nil {
				return err
			}

			err = debtServiceTx.SoftDeleteDebtsByExpenseID(ctx, expense.Id)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
