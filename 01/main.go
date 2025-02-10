package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/kataras/iris/v12"
)

// Định nghĩa cấu trúc dữ liệu gửi lên Groq API
type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseBody struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Hàm gọi Groq API
func callGroqAPI(prompt string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		apiKey = "your_groq_api_key_here" // Thay bằng API Key của bạn
	}

	// Tạo request body
	requestData := RequestBody{
		Model: "llama3-8b-8192", // Model AI sử dụng
		Messages: []Message{
			{Role: "system", Content: "You are a helpful AI."},
			{Role: "user", Content: prompt},
		},
	}

	// Chuyển request thành JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// Gửi request đến Groq API
	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Gửi request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Đọc response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse JSON response
	var responseData ResponseBody
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return "", err
	}

	// Trả về nội dung phản hồi từ AI
	if len(responseData.Choices) > 0 {
		return responseData.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from Groq API")
}

func main() {
	app := iris.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Fatal("GROQ_API_KEY is not set in .env file")
	}

	// Route chính hiển thị giao diện HTML
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Chat với AI - Groq API</title>
			<style>
				body { font-family: Arial, sans-serif; text-align: center; padding: 20px; }
				.container { width: 50%; margin: auto; }
				.input-container { display: flex; gap: 10px; }
				textarea { flex: 1; height: 100px; padding: 10px; resize: none; }
				button { width: 100px; height: 50px; padding: 10px; }
				#answer { text-align: left; margin-top: 20px; padding: 10px; border: 1px solid #ddd; border-radius: 5px; background: #f9f9f9; }
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Chat với AI - Groq API</h2>
				<div class="input-container">
					<textarea id="question" placeholder="Nhập câu hỏi..."></textarea>
					<button onclick="sendQuestion()">Gửi câu hỏi</button>
				</div>
				<h3>Kết quả:</h3>
				<div id="answer"></div>
			</div>

			<script>
				async function sendQuestion() {
					let question = document.getElementById("question").value;
					if (!question.trim()) return alert("Nhập câu hỏi trước!");
					let response = await fetch("/ask", {
						method: "POST",
						headers: { "Content-Type": "application/json" },
						body: JSON.stringify({ question })
					});

					let data = await response.json();
					document.getElementById("answer").innerHTML = data.answer;
				}
			</script>
		</body>
		</html>
		`)
	})

	// API xử lý câu hỏi
	app.Post("/ask", func(ctx iris.Context) {
		var input struct {
			Question string `json:"question"`
		}

		if err := ctx.ReadJSON(&input); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"error": "Invalid request"})
			return
		}

		// Gọi Groq API
		answer, err := callGroqAPI(input.Question)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Error calling Groq API"})
			return
		}

		// Chuyển Markdown thành HTML
		htmlFlags := html.CommonFlags | html.HrefTargetBlank
		renderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})
		mdToHTML := markdown.ToHTML([]byte(answer), nil, renderer)

		ctx.JSON(iris.Map{"answer": string(mdToHTML)})
	})

	// Chạy server tại port 8080
	app.Listen(":8080")
}
