package models

type StopPoint struct {
	BaseEntity
	ItineraryStopNr     int                 `gorm:"not null;uniqueIndex:uniq_itinerary_stop"`
	ItineraryId         uint                `gorm:"not null;uniqueIndex:uniq_itinerary_stop;uniqueIndex:uniq_itinerary_attraction"`
	Itinerary           Itinerary           `gorm:"foreignKey:ItineraryId;references:Id"`
	AttractionId        uint                `gorm:"not null;uniqueIndex:uniq_itinerary_attraction"`
	TouristicAttraction TouristicAttraction `gorm:"foreignKey:AttractionId;references:Id"`
}
