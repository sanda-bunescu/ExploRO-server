package models

import "time"

type BaseEntity struct {
	Id        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	DeletedAt *time.Time
}
