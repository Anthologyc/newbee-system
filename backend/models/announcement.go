package models

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	// 手动定义基础字段，确保 json tag 为小写
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 公告字段
	Title   string `json:"title" gorm:"type:varchar(200);not null"` // 标题
	Content string `json:"content" gorm:"type:text;not null"`       // 内容
	Status  bool   `json:"status" gorm:"default:true"`              // 状态 (true: 发布, false: 草稿) - 预留字段
}