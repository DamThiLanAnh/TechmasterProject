package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"TechmasterProject/03/config"
)

type GroqRequest struct {
	Input string `json:"input"`
}

type GroqResponse struct {
	Output string `json:"output"`
}

func CallGroqAPI(input string) (string, error) {
	url := "https://api.groq.com/openai/v1/chat/completions"
	requestBody, _ := json.Marshal(GroqRequest{Input: input})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.AppConfig.GroqAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var groqResponse GroqResponse
	json.Unmarshal(body, &groqResponse)

	return groqResponse.Output, nil
}
