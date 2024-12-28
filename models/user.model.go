package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string  `json:"name"`
	Login        string  `json:"login" gorm:"index" `
	Password     string  `json:"password"`
	Avatar       *string `json:"avatar"`
	Email        string  `json:"email" gorm:"index"`
	RefreshToken string  `json:"refresh_token" gorm:"type:text"`
}
