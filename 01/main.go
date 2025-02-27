package main

import (
	"TechmasterProject/01/config"
	"TechmasterProject/01/handler"
	"fmt"

	"github.com/kataras/iris/v12"
)

func main() {
	cfg := config.LoadConfig()
	app := iris.New()

	// Cấu hình route API
	app.Post("/ask", func(ctx iris.Context) {
		handlers.ChatHandler(ctx, cfg.GroqAPIKey)
	})

	// Phục vụ tệp tĩnh
	app.HandleDir("/", "./static")

	fmt.Println("Server đang chạy trên http://localhost:8080")
	app.Listen(":8080")
}
