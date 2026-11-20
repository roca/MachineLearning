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
		llms.TextParts(llms.ChatMessageTypeHuman, "Who is the president of the United States?"),
		llms.TextParts(llms.ChatMessageTypeHuman, "What is the capital of the United States?"),
		llms.TextParts(llms.ChatMessageTypeHuman, "What is the most popular sport in the United States?"),
		llms.TextParts(llms.ChatMessageTypeHuman, "Respond to each of the above questions a seperate line."),
	}

	ctx := context.Background()
	fmt.Println("Response from Ollama:")
	completion, err := llm.GenerateContent(ctx, messages, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}), llms.WithTemperature(0.7))
	if err != nil {
		log.Fatal(err)
	}

	_ = completion
	fmt.Println("\n\nDone streaming responses.")

}
