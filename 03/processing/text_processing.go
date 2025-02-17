package processing

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"regexp"
	"strings"
)

var stopWords = map[string]bool{
	"có": true, "thể": true, "tôi": true, "đang": true,
	"ở": true, "đến": true, "đây": true, "sẽ": true, "nhiều": true,
	"không": true, "gì": true, "vào": true, "là": true, "một": true,
}

var importantPhrases = []string{"đi thẳng", "rẽ phải", "cảm ơn", "không có gì", "chúc vui"}

func extractImportantWords(text string) []string {
	// Loại bỏ tên người nói (James:, Lan:)
	re := regexp.MustCompile(`(?i)(James:|Lan:)`)
	text = re.ReplaceAllString(text, "")

	// Xóa dấu câu và chuẩn hóa chữ thường bằng thư viện go-text
	text = text.Sanitize(text, text.EscapeNone)
	text = strings.ToLower(text)

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

func saveToPostgres(words []string) {
	// Kết nối PostgreSQL
	connStr := "user=your_user password=your_password dbname=your_database host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Tạo bảng dialog nếu chưa có
	_, err = db.Exec(`
		CREATE TABLE dialog (
			id BIGSERIAL PRIMARY KEY,
			lang VARCHAR(2) NOT NULL, 	-- vi: Vietnamese, en: English
			content TEXT NOT NULL  		-- Lưu toàn bộ nội dung hội thoại
			json JSONB NOT NULL 		-- Luu data dang json
			);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Tạo bảng word nếu chưa có
	_, err = db.Exec(`
		CREATE TABLE word IF NOT EXISTS  (
			id BIGSERIAL PRIMARY KEY,
			lang VARCHAR(2) NOT NULL, 	-- vi: Vietnamese, en: English
			content TEXT NOT NULL, 		-- Lưu gốc
		    json JSONB NOT NULL,		-- Luu data dang json
			translate TEXT NOT NULL 	-- Lưu dịch ra tiếng Anh
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	// Chuyển slice thành JSON
	jsonData, err := json.Marshal(map[string][]string{"words": words})
	if err != nil {
		log.Fatal(err)
	}

	// Lưu vào PostgreSQL
	_, err = db.Exec("INSERT INTO word (json) VALUES ($1)", jsonData)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Dữ liệu đã được lưu vào PostgreSQL!", string(jsonData))
}
