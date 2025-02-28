package services

import (
	"TechmasterProject/03/models"

	"gorm.io/gorm"
)

// Lấy danh sách từ vựng từ database bằng GORM

func GetWords(db *gorm.DB) ([]models.Word, error) {
	var words []models.Word
	if err := db.Find(&words).Error; err != nil {
		return nil, err
	}
	return words, nil
}

// Thêm từ mới vào database bằng GORM

func AddWord(db *gorm.DB, word models.Word) error {
	return db.Create(&word).Error
}
