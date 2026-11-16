package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	var options []openai.Option
	options = append(options, openai.WithModel("gpt-4"))
	options = append(options, openai.WithToken(os.Getenv("OPENAI_API_KEY"))) // You can also set your API key as an environment variable

	llm, err := openai.New(options...)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, "What is the capital of France?", llms.WithTemperature(0.7))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response from OpenAI:")
	fmt.Println(completion)
}
