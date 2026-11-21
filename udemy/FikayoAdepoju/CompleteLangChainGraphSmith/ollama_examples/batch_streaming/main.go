package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	var options []ollama.Option
	options = append(options, ollama.WithModel("mistral"))

	llm, err := ollama.New(options...)
	if err != nil {
		log.Fatal(err)
	}

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Respond to each of the following questions in French."),
		llms.TextParts(llms.ChatMessageTypeHuman, "Who is the president of the United States?"),
		llms.TextParts(llms.ChatMessageTypeHuman, "What is the capital of the United States?"),
		llms.TextParts(llms.ChatMessageTypeHuman, "What is the most popular sport in the United States?"),
	}

	ctx := context.Background()
	fmt.Println("Response from Ollama:")
	completion, err := llm.GenerateContent(ctx, messages, llms.WithTemperature(0.7))
	if err != nil {
		log.Fatal(err)
	}

	for _, choice := range completion.Choices {
		fmt.Println(choice.GenerationInfo)
		fmt.Println(choice.Content)
	}

	fmt.Println("\n\nDone streaming responses.")

}
