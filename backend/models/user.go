package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"-"` // 存储加密后的密码，不返回给前端
	Role     string `json:"role" gorm:"default:'user'"` // admin 或 user
}