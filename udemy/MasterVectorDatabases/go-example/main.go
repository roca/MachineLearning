package main

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func setupOpenAITestServer() (client *openai.Client) {
	config := openai.DefaultConfig(os.Getenv("OPEN_API_KEY"))
	// config.BaseURL = "/v1"
	// config.BaseURL = "https://chatgpt.regeneron.com/chat"
	fmt.Println(os.Getenv("OPEN_API_KEY"))
	fmt.Println(config.BaseURL)

	client = openai.NewClientWithConfig(config)
	return
}

func main() {

	//client := openai.NewClient(os.Getenv("OPEN_API_KEY"))
	client := setupOpenAITestServer()
	resp, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: []string{
			"The food was delicious and the waiter",
			"Other examples of embedding request",
		},
		Model: openai.AdaEmbeddingV2,
	})

	if err != nil {
		fmt.Printf("CreateEmbeddings error: %v\n", err)
		return
	}

	fmt.Println(resp.Data[0].Embedding)
}
