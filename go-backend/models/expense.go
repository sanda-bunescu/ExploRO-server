package models

import "time"

type Expense struct {
	BaseEntity
	Name        string
	Amount      float64
	Type        string
	Description string
	Date        time.Time
	GroupId     uint
	PayerId     string

	Debtors []Debt `gorm:"foreignKey:ExpenseId"`
}
