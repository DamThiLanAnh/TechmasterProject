package routes

import (
	"gorm.io/gorm"

	"TechmasterProject/03/handlers"
	"github.com/kataras/iris/v12"
)

func RegisterRoutes(app *iris.Application, db *gorm.DB) {
	// Thêm database vào context để dùng trong handlers
	app.Use(func(ctx iris.Context) {
		ctx.Values().Set("db", db)
		ctx.Next()
	})

	app.Post("/groq", handlers.GenerateResponse)
	app.Get("/words", handlers.ListWords)
	app.Post("/words", handlers.AddWord)
}
