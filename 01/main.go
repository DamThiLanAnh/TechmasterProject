package main

import (
	"TechmasterProject/01/config"
	"TechmasterProject/01/handler"
	"github.com/kataras/iris/v12"
	"log"
)

func main() {
	config.LoadEnv()

	app := iris.New()
	app.Get("/", handler.ServeHomePage)
	app.Post("/ask", handler.HandleAsk)

	log.Println("Server is running on :8080")
	app.Listen(":8080")
}
