package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
	Name      string         `json:"name"`
	CreatedAt datatypes.Date `json:"createdAt"`
	UpdatedAt datatypes.Date `json:"updatedAt"`
	Questions []Question     `json:"questions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Question struct {
	gorm.Model
	Name        string    `json:"name"`
	QuestionId  int       `json:"question_id"`
	Description string    `json:"description"`
	Variants    []Variant `json:"variants" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TestID      uint      `json:"test_id"`
}

type Variant struct {
	gorm.Model
	Name       string `json:"name"`
	VariantId  int    `json:"variant_id"`
	QuestionID uint   `json:"question_id"`
}
