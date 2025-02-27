package services

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
)

// Cấu trúc phản hồi từ API Groq

type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Gửi câu hỏi đến API Groq

func AskGroq(question, apiKey string) (string, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": "llama3-8b-8192", // ✅ Dùng model của Groq
			"messages": []map[string]string{
				{"role": "user", "content": question},
			},
		}).
		Post("https://api.groq.com/openai/v1/chat/completions") // ✅ URL đúng

	if err != nil {
		return "", err
	}

	// Giải mã JSON phản hồi từ API
	var groqResponse GroqResponse
	if err := json.Unmarshal(resp.Body(), &groqResponse); err != nil {
		return "", err
	}

	// Kiểm tra xem API có trả lời không
	if len(groqResponse.Choices) == 0 {
		return "", errors.New("Không có phản hồi từ Groq API")
	}

	return groqResponse.Choices[0].Message.Content, nil
}
