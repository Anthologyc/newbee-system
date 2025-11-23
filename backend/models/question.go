package models

import (
	"gorm.io/gorm"
)

type Question struct {
	gorm.Model
	// 类别
	Category string `json:"category" gorm:"type:varchar(100)"`

	// 🚀 改动 1: 使用 ENUM 类型
	// 限制只能是单选(single_choice)、多选(multiple_choice)或判断(judgment)
	QuestionType string `json:"question_type" gorm:"type:enum('single_choice','multiple_choice','judgment');default:'single_choice'"`

	// 🚀 改动 2: 使用 TEXT 类型
	// 这里的文本可能会很长
	QuestionText string `json:"question_text" gorm:"type:text"`

	// 🚀 改动 3: 显式指定为 JSON 类型
	// serializer:json 会让 GORM 自动把 Go 的 map 转换成 JSON 字符串存入
	// type:json 告诉 MySQL 使用原生的 JSON 数据类型
	Options map[string]string `json:"options" gorm:"type:json;serializer:json"`

	// 答案也存为 JSON 数组
	Answer []string `json:"answer" gorm:"type:json;serializer:json"`
}