package handlers

import (
	"database/sql"
	"net/http"

	"TechmasterProject/03/models"
	"TechmasterProject/03/services"
	"github.com/kataras/iris/v12"
)

// ListWords - Lấy danh sách từ vựng từ DB
func ListWords(ctx iris.Context) {
	db := ctx.Values().Get("db").(*sql.DB)
	words, err := services.GetWords(db)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to fetch words"})
		return
	}
	ctx.JSON(iris.Map{"words": words})
}

// AddWord - Thêm một từ mới vào DB
func AddWord(ctx iris.Context) {
	var word models.Word
	if err := ctx.ReadJSON(&word); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid input"})
		return
	}

	db := ctx.Values().Get("db").(*sql.DB)
	err := services.AddWord(db, word)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to add word"})
		return
	}

	ctx.JSON(iris.Map{"message": "Word added successfully"})
}
