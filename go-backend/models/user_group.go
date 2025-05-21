package models

type UserGroup struct {
	BaseEntity
	UserId  string `gorm:"type:varchar(255)"`
	User    Users  `gorm:"foreignKey:UserId;references:Id"`
	GroupId uint
	Group   Group `gorm:"foreignKey:GroupId;references:Id"`
}
