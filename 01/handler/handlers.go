package handlers

import (
	services "TechmasterProject/01/service"
	"github.com/kataras/iris/v12"
)

// ChatHandler xử lý yêu cầu chat
func ChatHandler(ctx iris.Context, apiKey string) {
	var request struct {
		Question string `json:"question"`
	}

	if err := ctx.ReadJSON(&request); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Dữ liệu không hợp lệ"})
		return
	}

	answer, err := services.AskGroq(request.Question, apiKey)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Không thể lấy dữ liệu từ API Groq"})
		return
	}

	ctx.JSON(iris.Map{"answer": answer})
}
