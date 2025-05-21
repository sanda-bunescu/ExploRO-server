package responses

import "time"

type ExpenseResponse struct {
	Id            uint           `json:"id"`
	Name          string         `json:"name"`
	GroupId       uint           `json:"group_id"`
	PayerUserName string         `json:"payer_user_name"`
	Date          time.Time      `json:"date"`
	Amount        float64        `json:"amount"`
	Description   string         `json:"description"`
	Type          string         `json:"type"`
	Debtors       []DebtResponse `json:"debtors"`
}
