package models

type Debt struct {
	BaseEntity
	ExpenseId   uint    `gorm:"not null"`
	UserId      string  `gorm:"not null"`
	AmountToPay float64 `gorm:"not null"`
}
