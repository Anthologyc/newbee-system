package models

import "gorm.io/gorm"

type Mistake struct {
	gorm.Model
	UserID     uint   `json:"user_id"`
	QuestionID uint   `json:"question_id"`
	Question   Question `json:"question" gorm:"foreignKey:QuestionID"` // 关联题目详情
}