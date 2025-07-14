package models

type UserGroup struct {
	BaseEntity
	UserId  string `gorm:"type:varchar(255);not null;uniqueIndex:uniq_user_group"`
	User    Users  `gorm:"foreignKey:UserId;references:Id"`
	GroupId uint   `gorm:"not null;uniqueIndex:uniq_user_group"`
	Group   Group  `gorm:"foreignKey:GroupId;references:Id"`
}
