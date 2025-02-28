package handlers

import (
	"TechmasterProject/03/services"
	"net/http"

	"github.com/kataras/iris/v12"
)

func GenerateResponse(ctx iris.Context) {
	var input struct {
		Text string `json:"text"`
	}

	if err := ctx.ReadJSON(&input); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid input"})
		return
	}

	response, err := services.CallGroqAPI(input.Text)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to call Groq API"})
		return
	}

	ctx.JSON(iris.Map{"response": response})
}
