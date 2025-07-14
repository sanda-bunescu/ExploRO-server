package models

type UserCities struct {
	BaseEntity
	UserID string `gorm:"type:varchar(255);not null;uniqueIndex:uniq_user_city"`
	User   Users  `gorm:"foreignKey:UserID;references:Id"`
	CityID uint   `gorm:"not null;uniqueIndex:uniq_user_city"`
	City   City   `gorm:"foreignKey:CityID;references:Id"`
}
