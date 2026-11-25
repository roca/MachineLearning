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
	var options []ollama.Option
	options = append(options, ollama.WithModel("mistral"))
	llm, err := ollama.New(options...)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	prompt1 := prompts.NewPromptTemplate("What is the capital of {{.question}}?", []string{"question"})
	prompt2 := prompts.NewPromptTemplate("What is the population of {{.question}}?", []string{"question"})

	fmt.Println("Response from Ollama:")

	choices, err := invoke(ctx, llm, map[string]any{"question": "France"}, prompt1, prompt2)
	if err != nil {
		log.Fatal(err)
	}

	for _, choice := range choices {
		fmt.Println(choice.GenerationInfo)
		fmt.Println(choice.Content)
	}

	fmt.Println("\n\nDone streaming responses.")
}

func invoke(ctx context.Context, llm *ollama.LLM, values map[string]any, prompts ...prompts.PromptTemplate) ([]*llms.ContentChoice, error) {
	var messsagesContent []llms.MessageContent

	if len(prompts) == 0 {
		return nil, fmt.Errorf("No prompts provided")
	}

	for _, prompt := range prompts {
		message, err := generateMessageFromPrompt(prompt, values)
		if err != nil {
			return nil, fmt.Errorf("Error formatting prompt: %v", err)
		}

		messsagesContent = append(messsagesContent, llms.TextParts(llms.ChatMessageTypeHuman, message))
	}

	completion, err := llm.GenerateContent(ctx, messsagesContent, llms.WithTemperature(0.7))
	if err != nil {
		return nil, fmt.Errorf("Error generating content: %v", err)
	}

	return completion.Choices, nil

}

func generateMessageFromPrompt(prompt prompts.PromptTemplate, input map[string]any) (string, error) {
	message, err := prompt.Format(input)
	if err != nil {
		return "", fmt.Errorf("Error formatting prompt: %v", err)
	}

	return message, nil
}
