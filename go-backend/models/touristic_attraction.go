package models

type TouristicAttraction struct {
	BaseEntity
	Name        string
	Description string
	Category    string
	ImageUrl    string
	OpenHours   string
	Fee         float32
	Link        string
	CityId      uint
}
