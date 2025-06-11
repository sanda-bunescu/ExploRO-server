package models

type StopPoint struct {
	BaseEntity
	ItineraryStopNr     int                 `gorm:"not null"`
	ItineraryId         uint                `gorm:"not null"`
	Itinerary           Itinerary           `gorm:"foreignKey:ItineraryId;references:Id"`
	AttractionId        uint                `gorm:"not null"`
	TouristicAttraction TouristicAttraction `gorm:"foreignKey:AttractionId;references:Id"`
}
