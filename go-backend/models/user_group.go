package models

type UserGroup struct {
	BaseEntity
	UserId  string `gorm:"type:varchar(255);not null"`
	User    Users  `gorm:"foreignKey:UserId;references:Id"`
	GroupId uint   `gorm:"not null"`
	Group   Group  `gorm:"foreignKey:GroupId;references:Id"`
}
