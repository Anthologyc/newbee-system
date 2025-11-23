package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `json:"password"` // 存加密后的 hash
	Role     string `json:"role"`     // "admin" 或 "user"
}