package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Nickname string `gorm:"not null" json:"nickname"`
	Email    string `grom:"not null; unique" json:"email"`
}
