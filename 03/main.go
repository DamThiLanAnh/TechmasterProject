package main

import (
	"TechmasterProject/03/processing"
	"fmt"
	"log"
	"os"
)

func main() {
	// Đoạn hội thoại mẫu
	conversation := `
	James: Chào bạn! Bạn có thể chỉ tôi đường đến hồ Hoàn Kiếm không?
	Lan: Chào bạn! Bạn đang ở đâu?
	James: Tôi đang ở phố Tràng Tiền.
	Lan: Từ đây, bạn đi thẳng, rẽ phải vào đường Đinh Tiên Hoàng, sẽ thấy hồ.
	James: Cảm ơn bạn nhiều!
	Lan: Không có gì, chúc bạn đi vui!
	`

	// Lưu đoạn hội thoại vào file
	fileName := "conversation.txt"
	err := os.WriteFile(fileName, []byte(conversation), 0644)
	if err != nil {
		log.Fatal("Lỗi khi ghi file:", err)
	}
	fmt.Println("✅ Đã lưu đoạn hội thoại vào file", fileName)

	// Đọc nội dung từ file
	content, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("Lỗi khi đọc file:", err)
	}

	// Lọc từ quan trọng
	importantWords := processing.ExtractImportantWords(string(content))

	// Dịch các từ quan trọng
	translatedWords := processing.TranslateWords(importantWords)

	// Lưu vào PostgreSQL
	processing.SaveToPostgres(importantWords, string(content), translatedWords)
}
