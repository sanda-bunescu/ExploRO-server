package models

type UserCities struct {
	BaseEntity
	UserID string `gorm:"type:varchar(255)"`
	User   Users  `gorm:"foreignKey:UserID;references:Id"`
	CityID uint
	City   City `gorm:"foreignKey:CityID;references:Id"`
}
