package model

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title    string `gorm:"not null" json:"title"`
	Content  string `gorm:"longtext" json:"content"`
	Status   int    `gorm:"default:0" json:"status"`
	Category string `json:"category"`
	UserId   uint   `gorm:"not null" json:"user_id"`
	User     *User  `gorm:"foreignKey:UserId" json:"user,omitempty"`
}
