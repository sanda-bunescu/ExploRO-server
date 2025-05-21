package models

type StopPoint struct {
	BaseEntity
	ItineraryStopNr     int
	ItineraryId         uint
	Itinerary           Itinerary `gorm:"foreignKey:ItineraryId;references:Id"`
	AttractionId        uint
	TouristicAttraction TouristicAttraction `gorm:"foreignKey:AttractionId;references:Id"`
}
