package models

import "time"

type Users struct {
	Id        string `gorm:"column:id;primaryKey;type:varchar(255)"`
	CreatedAt time.Time
	DeletedAt *time.Time
	Name      string `gorm:"type:varchar(255)"`
	Email     string `gorm:"type:varchar(255);not null"`

	Debts    []Debt    `gorm:"foreignKey:UserId"`
	Expenses []Expense `gorm:"foreignKey:PayerId"`
	Cities   []City    `gorm:"many2many:user_cities;"`
}
