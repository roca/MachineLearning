package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	chroma "github.com/amikos-tech/chroma-go"
	types "github.com/amikos-tech/chroma-go/types"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type OllamaEmbeddingFunction struct{}

func NewOllamaEmbeddingFunction() (*OllamaEmbeddingFunction, error) {
	return &OllamaEmbeddingFunction{}, nil
}

func (e *OllamaEmbeddingFunction) EmbedDocuments(ctx context.Context, documents []string) ([]*types.Embedding, error) {
	var embeddings []*types.Embedding

	for _, document := range documents {
		embedding, err := e.embedText(ctx, document)
		if err != nil {
			log.Fatalf("Error embedding document: %s \n", err)
		}

		embeddings = append(embeddings, embedding)
	}

	return embeddings, nil
}
func (e *OllamaEmbeddingFunction) EmbedQuery(ctx context.Context, document string) (*types.Embedding, error) {
	return e.embedText(ctx, document)
}

func (e *OllamaEmbeddingFunction) EmbedRecords(ctx context.Context, records []*types.Record, force bool) error {
	for _, record := range records {
		embedding, err := e.embedText(ctx, record.Document)
		if err != nil {
			log.Fatalf("Error embedding record: %s \n", err)
			return err
		}
		record.ID = uuid.UUID.String(uuid.New())
		record.Embedding = *embedding
	}
	return nil
}

func (e *OllamaEmbeddingFunction) embedText(ctx context.Context, document string) (*types.Embedding, error) {
	// Send POST REQUEST to Ollama API
	data := struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}{
		Model:  "mxbai-embed-large",
		Prompt: document,
	}
	jsonData, _ := json.Marshal(data)

	url := fmt.Sprintf("%s/api/embeddings", myEnv["OLLAMA_HOST"])

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalln("Can't connect to OLLAMA_HOST: ", err)
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		log.Fatalln("No response from OLLAMA_HOST: ", error)
		return nil, error
	}
	defer response.Body.Close()
	// log.Println("response Status:", response.Status)

	// Get the embeddings from the response

	type EmbeddingResponse struct {
		Embedding []float32 `json:"embedding"`
	}

	var embeddingResponse EmbeddingResponse
	json.NewDecoder(response.Body).Decode(&embeddingResponse)

	// log.Println("Embedding: ", embeddingResponse.Embedding)
	return &types.Embedding{
		ArrayOfFloat32: &embeddingResponse.Embedding,
	}, nil
}

type idGenerator struct{}

func (i *idGenerator) Generate(document string) string {
	return uuid.UUID.String(uuid.New())
}

var myEnv map[string]string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %s", err)
		return
	}

	myEnv, _ = godotenv.Read()
	fmt.Println("CHROMADB_TOKEN: ", myEnv["CHROMADB_TOKEN"])
	fmt.Println("CHROMADB_HOST: ", myEnv["CHROMADB_HOST"])
	fmt.Println("OLLAMA_HOST: ", myEnv["OLLAMA_HOST"])

}

func main() {
	chomaDBToken := myEnv["CHROMADB_TOKEN"]

	clientOptions := []chroma.ClientOption{
		chroma.WithTenant(types.DefaultTenant),
		chroma.WithDatabase(types.DefaultDatabase),
		chroma.WithAuth(types.NewTokenAuthCredentialsProvider(chomaDBToken, types.AuthorizationTokenHeader)),
		chroma.WithAuth(types.NewTokenAuthCredentialsProvider(chomaDBToken, types.XChromaTokenHeader)),
	}

	client, err := chroma.NewClient(myEnv["CHROMADB_HOST"], clientOptions...)
	if err != nil {
		log.Fatalf("Error creating client: %s \n", err.Error())
	}
	_ = client

	collectioCount, _ := client.CountCollections(context.Background())

	fmt.Println("CollectioCount:", collectioCount)

	collectionName := "test-collection"
	metadata := map[string]interface{}{}
	ctx := context.Background()

	embeddingFunction, _ := NewOllamaEmbeddingFunction() //create a new Ollama Embedding function
	distanceFunction := types.L2
	col, err := client.GetCollection(ctx, collectionName, embeddingFunction)
	if err != nil {
		log.Fatalf("Error get collection: %s \n", err.Error())
	}

	if col == nil {
		col, err = client.CreateCollection(ctx, collectionName, metadata, true, embeddingFunction, distanceFunction)
		if err != nil {
			log.Fatalf("Error create collection: %s \n", err.Error())
		}
	} else if col != nil {
		count, err := col.Count(ctx)
		if err != nil {
			log.Fatalf("Error counting documents: %s \n", err)
		}
		if count > 0 {
			// deleteRecords(ctx, col, ids)
			err = deleteAllRecords(ctx, col)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println("col: ", col)
	documents := []string{
		"The Matrix is everywhere. It is all around us.",
		"Unfortunately, no one can be told what the Matrix is.",
		"You can see it when you look out your window or when you turn on your television.",
		"Like everyone else you were born into bondage. Into a prison that you cannot taste or see or touch. A prison for your mind.",
		"It is the world that has been pulled over your eyes to blind you from the truth.",
		"Human beings are a disease, a cancer of this planet. You're a plague and we are the cure.",
		"Look out that window. You've had your time. The future is our world, Morpheus. The future is our time.",
	}

	ids := make([]string, len(documents))
	for i, _ := range documents {
		ids[i] = uuid.UUID.String(uuid.New())
	}

	metadatas := []map[string]interface{}{
		{"category": "quote", "speaker": "Morpheus"},
		{"category": "quote", "speaker": "Morpheus"},
		{"category": "quote", "speaker": "Morpheus"},
		{"category": "quote", "speaker": "Morpheus"},
		{"category": "quote", "speaker": "Morpheus"},
		{"category": "quote", "speaker": "Agent Smith"},
		{"category": "quote", "speaker": "Agent Smith"},
	}

	// records := []*types.Record{
	// 	{Document: documents[0], Metadata: metadatas[0]},
	// 	{Document: documents[1], Metadata: metadatas[1]},
	// }

	// recordSet := types.RecordSet{
	// 	Records:           records,
	// 	IDGenerator:       new(idGenerator),
	// 	EmbeddingFunction: embeddingFunction,
	// }

	// _, err := col.AddRecords(ctx, &recordSet)

	_, addError := col.Add(ctx, nil, metadatas, documents, ids)
	if addError != nil {
		log.Fatalf("Error adding documents: %s \n", addError)
	}

	countDocs, qrerr := col.Count(ctx)
	if qrerr != nil {
		log.Fatalf("Error counting documents: %s \n", qrerr)
	}

	fmt.Printf("countDocs: %v\n", countDocs) //this should result in 2
	where :=  make(map[string]interface{})
	where["speaker"] = "Morpheus"

	qr, qrerr := col.Query(ctx, []string{"Am I in a prison"}, 5,where, nil, nil)
	if qrerr != nil {
		log.Fatalf("Error querying documents: %s \n", qrerr)
	}

	fmt.Printf("qr: %v\n", qr.Documents[0][0]) //this should result in the document about dogs

}

func deleteAllRecords(ctx context.Context, col *chroma.Collection) error {
	results, err := col.Get(ctx, nil, nil, nil, nil)
	if err != nil {
		log.Fatalf("Error getting documents: %s \n", err)
	}
	ids := results.Ids

	idsDelete, err := col.Delete(ctx, ids, nil, nil)
	if err != nil {
		return fmt.Errorf("Error deleting documents: %s \n", err)
	}
	fmt.Printf("idsDeleted: %v\n", idsDelete) //this should result in 2
	countDocs, err := col.Count(ctx)
	if err != nil {
		return fmt.Errorf("Error counting documents: %s \n", err)

	}
	fmt.Printf("After deletion countDocs: %v\n", countDocs) //this should result in 2
	return nil

}
