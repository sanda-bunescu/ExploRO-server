package models

type Group struct {
	BaseEntity
	Name      string `gorm:"type:varchar(255);not null"`
	ImageURL  string
	TripPlans []TripPlan `gorm:"foreignKey:GroupId"`
	Expenses  []Expense  `gorm:"foreignKey:GroupId"`
}
