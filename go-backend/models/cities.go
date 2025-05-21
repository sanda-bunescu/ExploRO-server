package models

type City struct {
	Id          uint `gorm:"primaryKey"`
	Name        string
	Description string
	ImageURL    string
	IsDefault   bool

	Users                []Users               `gorm:"many2many:user_cities;"`
	TouristicAttractions []TouristicAttraction `gorm:"foreignKey:city_id"`
	TripPlans            []TripPlan            `gorm:"foreignKey:city_id"`
}
