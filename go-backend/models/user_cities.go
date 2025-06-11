package models

type UserCities struct {
	BaseEntity
	UserID string `gorm:"type:varchar(255);not null"`
	User   Users  `gorm:"foreignKey:UserID;references:Id"`
	CityID uint   `gorm:"not null"`
	City   City   `gorm:"foreignKey:CityID;references:Id"`
}
