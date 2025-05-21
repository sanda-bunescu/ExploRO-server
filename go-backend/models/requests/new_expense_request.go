package requests

import "time"

type NewExpenseRequest struct {
	Name        string        `json:"name"`
	GroupId     uint          `json:"group_id"`
	PayerId     string        `json:"payer_id"`
	Date        time.Time     `json:"date"`
	Amount      float64       `json:"amount"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	Debtors     []DebtRequest `json:"debtors"`
}
