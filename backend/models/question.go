package models

import (
	"time"
	"gorm.io/gorm"
)

type Question struct {
	// ❌ 删除 gorm.Model
	// gorm.Model 

	// ✅ 手动定义这4个字段，并加上 json tag
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // json:"-" 表示不返回给前端

	// --- 下面是你原本的业务字段 ---

	// 类别
	Category string `json:"category" gorm:"type:varchar(100)"`

	// 题型
	QuestionType string `json:"question_type" gorm:"type:enum('single_choice','multiple_choice','judgment');default:'single_choice'"`

	// 题目内容
	QuestionText string `json:"question_text" gorm:"type:text"`

	// 选项
	Options map[string]string `json:"options" gorm:"type:json;serializer:json"`

	// 答案
	Answer []string `json:"answer" gorm:"type:json;serializer:json"`
}