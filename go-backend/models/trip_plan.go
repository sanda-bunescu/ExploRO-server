package models

import "time"

type TripPlan struct {
	BaseEntity
	Name      string    `gorm:"type:varchar(255);not null"`
	StartDate time.Time `gorm:"not null"`
	EndDate   time.Time `gorm:"not null"`
	GroupId   uint      `gorm:"not null"`
	CityId    uint      `gorm:"not null"`

	Itineraries []Itinerary `gorm:"foreignKey:TripPlanId"`
}
