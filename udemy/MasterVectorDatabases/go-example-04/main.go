package main

import (
	"context"
	"fmt"
	"log"
	"os"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/joho/godotenv"
)

type OllamaEmbeddingFunction struct{}

func NewOllamaEmbeddingFunction() (*OllamaEmbeddingFunction, error) {
	return &OllamaEmbeddingFunction{}, nil
}

func (e *OllamaEmbeddingFunction) EmbedDocuments(ctx context.Context, documents []string) ([]*types.Embedding, error) {
	return []*types.Embedding{}, nil
}
func (e *OllamaEmbeddingFunction) EmbedQuery(ctx context.Context, document string) (*types.Embedding, error) {
	return &types.Embedding{}, nil
}
func (e *OllamaEmbeddingFunction) EmbedRecords(ctx context.Context, records []*types.Record, force bool) error {
	return nil
}

func main() {
	client, err := chroma.NewClient("http://localhost:8000")
	if err != nil {
		log.Fatalf("Error creating client: %s \n", err.Error())
	}

	collectionName := "test-collection"
	metadata := map[string]interface{}{}
	ctx := context.Background()
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %s", err)
		return
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	fmt.Println("openaiApiKey: ", apiKey)

	embeddingFunction, _ := NewOllamaEmbeddingFunction() //create a new Ollama Embedding function
	distanceFunction := types.L2
	_, errRest := client.Reset(ctx) //reset the database
	if errRest != nil {
		log.Fatalf("Error resetting database: %s \n", errRest.Error())
	}
	col, err := client.CreateCollection(ctx, collectionName, metadata, true, embeddingFunction, distanceFunction)
	if err != nil {
		fmt.Printf("Error create collection: %s \n", err.Error())
		return
	}
	//fmt.Println("col: ", col)
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
