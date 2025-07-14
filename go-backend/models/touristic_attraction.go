package models

type TouristicAttraction struct {
	BaseEntity
	Name        string `gorm:"type:varchar(255);not null;uniqueIndex:uniq_city_attraction"`
	Description string
	Category    string `gorm:"type:varchar(255)"`
	ImageUrl    string
	OpenHours   string
	Fee         float32
	Link        string
	CityId      uint `gorm:"not null;uniqueIndex:uniq_city_attraction"`
}
