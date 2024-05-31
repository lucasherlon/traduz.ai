package geminiapi

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// Generate the prompt to be translated
func GeneratePrompt(input string) string {
	erro := godotenv.Load(".env")
	if erro != nil {
		fmt.Println("Error loading .env file:", erro)
	}

	idiomaSaida := "português brasileiro"

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer func(client *genai.Client) {
		err := client.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client)

	model := client.GenerativeModel("gemini-1.5-flash-latest")
	fmt.Println(input)

	prompt := fmt.Sprintf("Por gentileza, traduza o seguinte texto para o %v: %v", idiomaSaida, input)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		log.Fatal(err)
	}
	response := showResponse(resp)
	return response
}

// Show the response generatad from the model
func showResponse(resp *genai.GenerateContentResponse) string {
	var response string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				response = response + fmt.Sprintf("%v", part)
			}
		}
	}
	return response
}
