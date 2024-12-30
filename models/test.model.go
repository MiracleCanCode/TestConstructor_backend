package models

import (
	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
	Name        string     `json:"name"`
	AuthorLogin string     `json:"author_login"`
	Questions   []Question `json:"questions" gorm:"foreignKey:TestID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Question struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Variants    []Variant `json:"variants" gorm:"foreignKey:QuestionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TestID      uint      `json:"test_id"`
}

type Variant struct {
	gorm.Model
	Name       string `json:"name"`
	QuestionID uint   `json:"question_id"`
	IsCorrect  bool   `json:"is_correct"`
}
