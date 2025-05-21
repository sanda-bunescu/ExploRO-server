package models

type Debt struct {
	BaseEntity
	ExpenseId   uint
	UserId      string
	AmountToPay float64
}
