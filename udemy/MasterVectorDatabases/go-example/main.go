package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	openai "github.com/0x9ef/openai-go"
)

func main() {
	fmt.Println(os.Getenv("OPEN_API_KEY"))

	e := openai.New(os.Getenv("OPEN_API_KEY"))
	r, err := e.Completion(context.Background(), &openai.CompletionOptions{
		// Choose model, you can see list of available models in models.go file
		Model: openai.ModelGPT3TextDavince,
		// Text to completion
		Prompt: []string{"Write a little bit of Wikipedia. What is that?"},
	})
	if err != nil {
		log.Fatal(err)
	}

	if b, err := json.MarshalIndent(r, "", "  "); err != nil {
		panic(err)
	} else {
		fmt.Println(string(b))
	}

	// Wikipedia is a free online encyclopedia, created and edited by volunteers.
	fmt.Println("What is the Wikipedia?", r.Choices[0].Text)
}
