package models

type Itinerary struct {
	BaseEntity
	DayNr      uint `gorm:"not null"`
	TripPlanId uint `gorm:"not null"`
}
