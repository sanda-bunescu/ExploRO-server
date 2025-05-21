package models

import "time"

type TripPlan struct {
	BaseEntity
	Name      string
	StartDate time.Time
	EndDate   time.Time
	GroupId   uint
	CityId    uint

	Itineraries []Itinerary `gorm:"foreignKey:TripPlanId"`
}
