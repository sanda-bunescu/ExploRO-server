package models

type Group struct {
	BaseEntity
	Name      string
	ImageURL  string
	TripPlans []TripPlan `gorm:"foreignKey:GroupId"`
	Expenses  []Expense  `gorm:"foreignKey:GroupId"`
}
