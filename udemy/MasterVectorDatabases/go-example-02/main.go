package main

import (
	"context"
	"fmt"
	"log"
	"os"

	chroma "github.com/amikos-tech/chroma-go"
	openai "github.com/amikos-tech/chroma-go/openai"
	godotenv "github.com/joho/godotenv"
)

func main() {
	client := chroma.NewClient("http://localhost:8000")
	collectionName := "test-collection"
	metadata := map[string]interface{}{}
	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %s", err)
		return
	}
	embeddingFunction := openai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY")) //create a new OpenAI Embedding function
	distanceFunction := chroma.L2
	_, errRest := client.Reset(ctx) //reset the database
	if errRest != nil {
		log.Fatalf("Error resetting database: %s \n", errRest.Error())
	}
	col, err := client.CreateCollection(ctx, collectionName, metadata, true, embeddingFunction, distanceFunction)
	if err != nil {
		fmt.Printf("Error create collection: %s \n", err.Error())
		return
	}
	documents := []string{
		"This is a document about cats. Cats are great.",
		"this is a document about dogs. Dogs are great.",
	}
	ids := []string{
		"ID1",
		"ID2",
	}

	metadatas := []map[string]interface{}{
		{"key1": "value1"},
		{"key2": "value2"},
	}
	_, addError := col.Add(ctx, nil, metadatas, documents, ids)
	if addError != nil {
		log.Fatalf("Error adding documents: %s \n", addError)
	}
	countDocs, qrerr := col.Count(ctx)
	if qrerr != nil {
		log.Fatalf("Error counting documents: %s \n", qrerr)
	}
	fmt.Printf("countDocs: %v\n", countDocs) //this should result in 2
	qr, qrerr := col.Query(ctx, []string{"I love dogs"}, 5, nil, nil, nil)
	if qrerr != nil {
		log.Fatalf("Error querying documents: %s \n", qrerr)
	}
	fmt.Printf("qr: %v\n", qr.Documents[0][0]) //this should result in the document about dogs
}
