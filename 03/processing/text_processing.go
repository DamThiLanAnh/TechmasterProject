package processing

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	_ "github.com/lib/pq"
)

type TranslatedWord struct {
	Vi string `json:"vi"`
	En string `json:"en"`
}

var stopWords = map[string]bool{
	"có": true, "thể": true, "tôi": true, "đang": true,
	"ở": true, "đến": true, "đây": true, "sẽ": true, "nhiều": true,
	"không": true, "gì": true, "vào": true, "là": true, "một": true,
}

var importantPhrases = []string{"đi thẳng", "rẽ phải", "cảm ơn", "không có gì", "chúc vui"}

var translationDictionary = map[string]string{
	"chào": "hello", "bạn": "you", "chỉ": "show", "đường": "road",
	"hồ": "lake", "phố": "street", "đi thẳng": "go straight",
	"rẽ phải": "turn right", "thấy": "see", "cảm ơn": "thank you",
	"không có gì": "you're welcome", "chúc": "wish", "vui": "happy",
}

func ExtractImportantWords(text string) []string {
	re := regexp.MustCompile(`(?i)(James:|Lan:)`)
	text = re.ReplaceAllString(text, "")

	text = strings.ToLower(text)
	text = strings.NewReplacer(".", "", ",", "", "!", "", "?", "").Replace(text)

	words := strings.Fields(text)
	wordSet := make(map[string]bool)

	for _, word := range words {
		if !stopWords[word] {
			wordSet[word] = true
		}
	}

	for _, phrase := range importantPhrases {
		if strings.Contains(text, phrase) {
			wordSet[phrase] = true
		}
	}

	var importantWords []string
	for word := range wordSet {
		importantWords = append(importantWords, word)
	}

	return importantWords
}

func TranslateWords(words []string) []TranslatedWord {
	var translatedWords []TranslatedWord
	for _, word := range words {
		if en, exists := translationDictionary[word]; exists {
			translatedWords = append(translatedWords, TranslatedWord{Vi: word, En: en})
		}
	}
	return translatedWords
}

func SaveToPostgres(words []string, originalText string, translatedWords []TranslatedWord) {
	
	if err != nil {
		log.Fatal("Lỗi kết nối DB:", err)
	}
	defer dbConn.Close()

	_, err = dbConn.Exec(`CREATE TABLE IF NOT EXISTS dialog (
		id BIGSERIAL PRIMARY KEY,
		lang VARCHAR(2) NOT NULL,
		content TEXT NOT NULL,
		json JSONB NOT NULL
	);`)
	if err != nil {
		log.Fatal("Lỗi tạo bảng dialog:", err)
	}

	_, err = dbConn.Exec(`CREATE TABLE IF NOT EXISTS word (
		id BIGSERIAL PRIMARY KEY,
		lang VARCHAR(2) NOT NULL,
		content TEXT NOT NULL,
		json JSONB NOT NULL,
		translate TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatal("Lỗi tạo bảng word:", err)
	}

	jsonWords, err := json.Marshal(words)
	if err != nil {
		log.Fatal("Lỗi chuyển JSON:", err)
	}

	_, err = dbConn.Exec("INSERT INTO dialog (lang, content, json) VALUES ($1, $2, $3)", "vi", originalText, jsonWords)
	if err != nil {
		log.Fatal("Lỗi lưu vào dialog:", err)
	}

	jsonTranslated, err := json.Marshal(translatedWords)
	if err != nil {
		log.Fatal("Lỗi chuyển JSON:", err)
	}

	for _, word := range translatedWords {
		_, err := dbConn.Exec("INSERT INTO word (lang, content, json, translate) VALUES ($1, $2, $3, $4)", "vi", word.Vi, jsonTranslated, word.En)
		if err != nil {
			log.Fatal("Lỗi lưu vào word:", err)
		}
	}

	fmt.Println("Dữ liệu đã được lưu vào PostgreSQL!", string(jsonTranslated))
}
