package handler

import (
	"TechmasterProject/01/service"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/kataras/iris/v12"
	"log"
	"os"
)

type Input struct {
	Question string `json:"question"`
}

func ServeHomePage(ctx iris.Context) {
	htmlContent, err := os.ReadFile("templates/index.html")
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString("Error loading page")
		return
	}
	ctx.HTML(string(htmlContent))
}

func HandleAsk(ctx iris.Context) {
	var input Input
	if err := ctx.ReadJSON(&input); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request"})
		return
	}

	answer, err := service.CallGroqAPI(input.Question)
	if err != nil {
		log.Println("Error calling Groq API:", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Error calling Groq API"})
		return
	}

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	renderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})
	mdToHTML := markdown.ToHTML([]byte(answer), nil, renderer)

	ctx.JSON(iris.Map{"answer": string(mdToHTML)})
}
