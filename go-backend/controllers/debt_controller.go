package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sanda-bunescu/ExploRO/services"
	"net/http"
	"strconv"
)

type DebtController struct {
	DebtService services.DebtServiceInterface
}

func NewDebtController(DebtService services.DebtServiceInterface) *DebtController {
	return &DebtController{
		DebtService: DebtService,
	}
}

func (dc *DebtController) GetByGroupIdAndUser(ctx *gin.Context) {
	groupIdStr := ctx.Query("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil || groupId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid groupId provided"})
		return
	}

	debts, err := dc.DebtService.GetByUserIdAndGroupId(ctx, uint(groupId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"debts": debts})
}

func (dc *DebtController) DeleteDebt(ctx *gin.Context) {
	debtIdStr := ctx.Query("Id")
	debtId, err := strconv.Atoi(debtIdStr)
	if err != nil || debtId <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid debtId provided"})
		return
	}

	err = dc.DebtService.Delete(ctx, uint(debtId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Debt deleted successfully"})
}
