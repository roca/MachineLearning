package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
)

func main() {

	prompt := prompts.NewPromptTemplate("What is the capital of {{.question}}?", []string{"question"})

	message, err := prompt.Format(map[string]any{
		"question": "France",
	})
	if err != nil {
		fmt.Println("Error formatting prompt:", err)
		return
	}

	fmt.Println("Formatted Prompt:", message)

	var options []ollama.Option
	options = append(options, ollama.WithModel("mistral"))
	llm, err := ollama.New(options...)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, message),
	}

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
