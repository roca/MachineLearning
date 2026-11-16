package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

/*
The Ollama LLM is a loaclly hosted.
So there is no need for an API key or any other authentication.
You must have the Ollama app installed and running on your machine for this to work.
You also have to pull the model by the command "ollama pull <model_name>" before you can use it.
*/

func main() {
	var options []ollama.Option
	options = append(options, ollama.WithModel("mistral"))

	llm, err := ollama.New(options...)
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
