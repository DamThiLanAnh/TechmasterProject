package processing

import (
	"TechmasterProject/03/db"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	_ "github.com/lib/pq"
)

// Danh sách stop words
var stopWords = map[string]bool{
	"có": true, "thể": true, "tôi": true, "đang": true,
	"ở": true, "đến": true, "đây": true, "sẽ": true, "nhiều": true,
	"không": true, "gì": true, "vào": true, "là": true, "một": true,
}

// Cụm từ quan trọng
var importantPhrases = []string{"đi thẳng", "rẽ phải", "cảm ơn", "không có gì", "chúc vui"}

// Xử lý văn bản, trích xuất từ quan trọng

func ExtractImportantWords(text string) []string {
	// Loại bỏ tên người nói
	re := regexp.MustCompile(`(?i)(James:|Lan:)`)
	text = re.ReplaceAllString(text, "")

	// Chuyển chữ thường và xóa dấu câu
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, ".", "")
	text = strings.ReplaceAll(text, ",", "")
	text = strings.ReplaceAll(text, "!", "")
	text = strings.ReplaceAll(text, "?", "")

	// Tách từ
	words := strings.Fields(text)

	// Lọc bỏ stopwords
	var importantWords []string
	wordSet := make(map[string]bool)

	for _, word := range words {
		if !stopWords[word] {
			wordSet[word] = true
		}
	}

	// Thêm cụm từ quan trọng
	for _, phrase := range importantPhrases {
		if strings.Contains(text, phrase) {
			wordSet[phrase] = true
		}
	}

	// Chuyển map thành slice
	for word := range wordSet {
		importantWords = append(importantWords, word)
	}

	return importantWords
}

// Lưu dữ liệu vào PostgreSQL

func SaveToPostgres(words []string, originalText string) {
	// Kết nối đến PostgreSQL
	dbConn, err := db.ConnectToDB()
	if err != nil {
		log.Fatal("Lỗi kết nối DB:", err)
	}
	defer dbConn.Close()

	// Tạo bảng dialog nếu chưa có
	_, err = dbConn.Exec(`
		CREATE TABLE IF NOT EXISTS dialog (
			id BIGSERIAL PRIMARY KEY,
			lang VARCHAR(2) NOT NULL,
			content TEXT NOT NULL,
			json JSONB NOT NULL
		);
	`)
	if err != nil {
		log.Fatal("Lỗi tạo bảng dialog:", err)
	}

	// Tạo bảng word nếu chưa có
	_, err = dbConn.Exec(`
		CREATE TABLE IF NOT EXISTS word (
			id BIGSERIAL PRIMARY KEY,
			lang VARCHAR(2) NOT NULL,
			content TEXT NOT NULL,
		    json JSONB NOT NULL,
			translate TEXT NOT NULL
		);
	`)
	if err != nil {
		log.Fatal("Lỗi tạo bảng word:", err)
	}

	// Chuyển danh sách từ quan trọng thành JSON
	jsonData, err := json.Marshal(map[string][]string{"words": words})
	if err != nil {
		log.Fatal("Lỗi chuyển JSON:", err)
	}

	// Lưu nội dung hội thoại vào bảng dialog
	_, err = dbConn.Exec("INSERT INTO dialog (lang, content, json) VALUES ($1, $2, $3)", "vi", originalText, jsonData)
	if err != nil {
		log.Fatal("Lỗi lưu vào dialog:", err)
	}

	// Lưu danh sách từ vào bảng word
	for _, word := range words {
		_, err := dbConn.Exec("INSERT INTO word (lang, content, json, translate) VALUES ($1, $2, $3, $4)", "vi", word, jsonData, "")
		if err != nil {
			log.Fatal("Lỗi lưu vào word:", err)
		}
	}

	fmt.Println("Dữ liệu đã được lưu vào PostgreSQL!", string(jsonData))
}
