package models

type TouristicAttraction struct {
	BaseEntity
	Name        string `gorm:"type:varchar(255);not null"`
	Description string
	Category    string `gorm:"type:varchar(255)"`
	ImageUrl    string
	OpenHours   string
	Fee         float32
	Link        string
	CityId      uint `gorm:"not null"`
}
