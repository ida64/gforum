package database

import "gorm.io/gorm"

type ActionModel struct {
	gorm.Model

	Method string `gorm:"not null"`
	Url    string `gorm:"not null"`
	IP     string `gorm:"not null"`

	UserID uint
}
