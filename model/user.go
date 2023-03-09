package model

import "gorm.io/gorm"

// User struct
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Verified bool   `gorm:"not null" json:"verified"`
	IsAdmin  bool   `gorm:"not null" json:"isAdmin"`
}
