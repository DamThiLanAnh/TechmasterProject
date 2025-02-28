package database

import (
	"database/sql"
	"fmt"
	"log"

	"TechmasterProject/01/config"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	// Xây dựng DSN (Data Source Name)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
	)

	// Mở kết nối đến PostgreSQL
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Println("Lỗi khi mở kết nối database:", err)
		return nil, err
	}

	// Kiểm tra kết nối bằng Ping()
	if err := db.Ping(); err != nil {
		log.Println("Không thể kết nối đến database:", err)
		return nil, err
	}

	log.Println("✅ Kết nối database thành công!")
	return db, nil
}
