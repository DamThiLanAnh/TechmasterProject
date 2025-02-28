package models

import "gorm.io/gorm"

type Dialog struct {
	gorm.Model
	Lang    string `json:"lang"`
	Content string `json:"content"`
}
