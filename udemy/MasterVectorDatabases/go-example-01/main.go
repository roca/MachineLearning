package main

import (
	"log"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

var client *openai.Client

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	azure()
	//client = openai.NewClient(os.Getenv("OPEN_API_KEY"))
	// config = openai.DefaultConfig(os.Getenv("OPEN_API_KEY"))
	// // config.BaseURL = "/v1"
	// // config.BaseURL = "https://chatgpt.regeneron.com/chat"
	// fmt.Println(os.Getenv("OPEN_API_KEY"))
	// fmt.Println(config.BaseURL)

	// client = openai.NewClientWithConfig(config)
	return
}

func main() {}

// func main() {

// 	resp, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
// 		Input: []string{
// 			"dog",
// 		},
// 		Model: openai.AdaEmbeddingV2,
// 	})

// 	if err != nil {
// 		fmt.Printf("CreateEmbeddings error: %v\n", err)
// 		return
// 	}

// 	fmt.Println("Embedding length for a 'dog': ", len(resp.Data[0].Embedding))

// 	v1 := []float64{2.8, 1.6, 3.8, 2.2}
// 	v2 := []float64{1.3, 3.5, 2.2, 0.9}
// 	v3 := convertFloat32ToFloat64(resp.Data[0].Embedding)

// 	sqlite()
// 	defer sqliteDatabase.Close()
// 	insertVector(v1)
// 	insertVector(v2)
// 	insertVector(v3)

// 	vectors := findVectors(v2)
// 	fmt.Println("Found vectors: ", vectors)
// }

func convertFloat32ToFloat64(fs []float32) []float64 {
	result := make([]float64, len(fs))
	for i, f := range fs {
		result[i] = float64(f)
	}
	return result
}
