package requests

import "time"

type EditExpenseRequest struct {
	Id          uint          `json:"id"`
	Name        string        `json:"name"`
	Amount      float64       `json:"amount"`
	Type        string        `json:"type"`
	Date        time.Time     `json:"date"`
	Description string        `json:"description"`
	Debtors     []DebtRequest `json:"debtors"`
}
