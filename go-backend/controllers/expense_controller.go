package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/models/requests"
	"github.com/sanda-bunescu/ExploRO/services"
	"net/http"
	"strconv"
)

type ExpenseController struct {
	ExpenseService services.ExpenseServiceInterface
}

func NewExpenseController(expenseService services.ExpenseServiceInterface) *ExpenseController {
	return &ExpenseController{
		ExpenseService: expenseService,
	}
}

func (ec *ExpenseController) GetAllExpensesByGroupId(ctx *gin.Context) {
	groupIdStr := ctx.Query("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil || groupId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid groupId provided"})
		return
	}

	expenses, err := ec.ExpenseService.GetAllExpensesByGroupId(ctx, uint(groupId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expenses", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

func (ec *ExpenseController) CreateExpenseWithDebts(ctx *gin.Context) {
	var req requests.NewExpenseRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if err := ec.ExpenseService.CreateExpenseWithDebts(ctx, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense with debts", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Expense with debts created successfully"})
}

func (ec *ExpenseController) EditExpenseWithDebts(ctx *gin.Context) {
	var req requests.EditExpenseRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	err := ec.ExpenseService.EditExpenseWithDebts(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit expense", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Expense updated successfully"})
}

func (ec *ExpenseController) DeleteExpenseByID(ctx *gin.Context) {
	expenseIdString := ctx.Query("id")

	expenseId, err := strconv.Atoi(expenseIdString)
	if err != nil || expenseId < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id provided"})
		return
	}

	err = ec.ExpenseService.SoftDeleteExpenseByID(ctx, uint(expenseId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}
