package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string  `json:"name"`
	Login        string  `json:"login" gorm:"index,unique"`
	Password     string  `json:"password"`
	Avatar       *string `json:"avatar"`
	Email        string  ` json:"email" gorm:"index,unique"`
	RefreshToken string  `json:"refresh_token" gorm:"type:text"`
	Tests        []Test  `json:"tests" gorm:"foreignKey:UserID"`
}
