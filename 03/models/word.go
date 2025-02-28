package models

import "gorm.io/gorm"

type Word struct {
	gorm.Model
	Lang      string `json:"lang"`
	Content   string `json:"content"`
	Translate string `json:"translate"`
}
