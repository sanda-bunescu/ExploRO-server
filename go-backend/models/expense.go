package models

import "time"

type Expense struct {
	BaseEntity
	Name        string  `gorm:"type:varchar(255);not null"`
	Amount      float64 `gorm:"not null"`
	Type        string  `gorm:"type:varchar(255)"`
	Description string
	Date        time.Time
	GroupId     uint   `gorm:"not null"`
	PayerId     string `gorm:"not null"`

	Debtors []Debt `gorm:"foreignKey:ExpenseId"`
}
