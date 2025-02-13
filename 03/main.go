package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

// Định nghĩa cấu trúc cho mỗi đối tượng hội thoại

type Dialog struct {
	Lang    string `json:"lang"`
	Content string `json:"content"`
}

// Định nghĩa cấu trúc cho danh sách từ quan trọng

type Words struct {
	Words []string `json:"words"`
}

func main() {
	// Kết nối đến cơ sở dữ liệu PostgreSQL
	connStr := "user=myuser dbname=mydatabase password=mypassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Lỗi khi kết nối DB:", err)
	}
	defer db.Close()

	// Kiểm tra kết nối cơ sở dữ liệu
	if err = db.Ping(); err != nil {
		log.Fatal("Lỗi khi ping DB:", err)
	}

	// Tạo bảng nếu chưa tồn tại
	createTables(db)

	// Đọc nội dung file text
	fileContent, err := ioutil.ReadFile("conversation.txt")
	if err != nil {
		log.Fatal("Lỗi khi đọc file:", err)
	}

	// Chia nội dung thành các dòng
	lines := strings.Split(string(fileContent), "\n")

	// Danh sách cụm từ quan trọng cần giữ lại
	importantPhrases := []string{
		"chào", "bạn", "chỉ", "đường", "hồ", "đang", "phố",
		"đi thẳng", "rẽ phải", "thấy", "cảm ơn", "không có gì",
		"chúc", "vui",
	}

	// Danh sách từ không quan trọng (danh từ riêng)
	stopWords := map[string]bool{
		"Hoàn Kiếm":      true,
		"Tràng Tiền":     true,
		"Đinh Tiên Hoàng": true,
		"Lan":             true,
		"James":           true,
	}

	// Mảng để chứa từ quan trọng
	importantWords := make(map[string]bool)

	// Duyệt qua từng dòng để lọc từ quan trọng
	for _, line := range lines {
		if line == "" {
			continue
		}

		// Phân tách lang và content
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		lang := strings.TrimSpace(parts[0])
		content := normalizeText(strings.TrimSpace(parts[1]))

		// Lọc các từ quan trọng
		for _, phrase := range importantPhrases {
			if strings.Contains(content, phrase) && !stopWords[phrase] {
				importantWords[phrase] = true
			}
		}

		// Chèn vào bảng dialog
		dialogID := insertDialog(db, lang, content)
		updateDialogJSON(db, dialogID, lang, content)
	}

	// Lưu danh sách từ vào database
	saveWordsToDB(db, importantWords)

	fmt.Println("Dữ liệu đã được lưu vào database thành công!")
}

// Hàm chuẩn hóa văn bản (loại bỏ dấu câu, viết thường)
func normalizeText(text string) string {
	// Xử lý chuyển về chữ thường và loại bỏ dấu câu không cần thiết
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, ".", "")
	text = strings.ReplaceAll(text, ",", "")
	text = strings.ReplaceAll(text, "!", "")
	text = strings.ReplaceAll(text, "?", "")
	return text
}

// Hàm tạo bảng nếu chưa tồn tại
func createTables(db *sql.DB) {
	createDialogTable := `
		CREATE TABLE IF NOT EXISTS dialog (
			id SERIAL PRIMARY KEY,
			lang TEXT NOT NULL,
			content TEXT NOT NULL,
			json JSONB
		);
	`
	createWordTable := `
		CREATE TABLE IF NOT EXISTS word (
			id SERIAL PRIMARY KEY,
			lang TEXT NOT NULL,
			json JSONB,
			content TEXT,
			translate TEXT
		);
	`
	if _, err := db.Exec(createDialogTable); err != nil {
		log.Fatal("Lỗi khi tạo bảng dialog:", err)
	}
	if _, err := db.Exec(createWordTable); err != nil {
		log.Fatal("Lỗi khi tạo bảng word:", err)
	}
}

// Hàm chèn dữ liệu vào bảng dialog
func insertDialog(db *sql.DB, lang, content string) int {
	query := `INSERT INTO dialog (lang, content) VALUES ($1, $2) RETURNING id;`
	var id int
	err := db.QueryRow(query, lang, content).Scan(&id)
	if err != nil {
		log.Fatal("Lỗi khi chèn dữ liệu vào dialog:", err)
	}
	return id
}

// Hàm cập nhật JSON vào bảng dialog
func updateDialogJSON(db *sql.DB, id int, lang, content string) {
	dialogJSON, err := json.Marshal(Dialog{Lang: lang, Content: content})
	if err != nil {
		log.Fatal("Lỗi khi tạo JSON:", err)
	}
	query := `UPDATE dialog SET json = $1 WHERE id = $2;`
	_, err = db.Exec(query, dialogJSON, id)
	if err != nil {
		log.Fatal("Lỗi khi cập nhật JSON vào dialog:", err)
	}
}

// Hàm lưu danh sách từ quan trọng vào database
func saveWordsToDB(db *sql.DB, importantWords map[string]bool) {
	// Chuyển map thành mảng
	wordList := make([]string, 0, len(importantWords))
	for word := range importantWords {
		wordList = append(wordList, word)
	}

	// Chuyển thành JSON
	wordsJSON, err := json.Marshal(Words{Words: wordList})
	if err != nil {
		log.Fatal("Lỗi khi tạo JSON từ danh sách từ:", err)
	}

	// Ghi dữ liệu vào bảng word
	query := `
		INSERT INTO word (id, lang, json, content, translate)
		VALUES (1, $1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE 
		SET lang = EXCLUDED.lang, json = EXCLUDED.json, content = EXCLUDED.content, translate = EXCLUDED.translate;
	`
	_, err = db.Exec(query, "vi", wordsJSON, "Nội dung mặc định", "Dịch mặc định")
	if err != nil {
		log.Fatal("Lỗi khi ghi dữ liệu vào bảng word:", err)
	}

	fmt.Println("Danh sách từ quan trọng đã được lưu vào database: ", string(wordsJSON))
}
