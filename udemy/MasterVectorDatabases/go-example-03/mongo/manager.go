package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ch := make(chan string, 1)

	// Create a context with a timeout of 5 seconds
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// Start the ConnectToMilvus function
	go ConnectToMongo(ctxTimeout, ch)

	select {
	case <-ctxTimeout.Done():
		fmt.Printf("Context cancelled: %v\n", ctxTimeout.Err())
	case result := <-ch:
		fmt.Printf("Received: %s\n", result)
	}

}

func ConnectToMongo(ctx context.Context, ch chan string) {

	username := os.Getenv("BIOREGISTRY_USER")
	password := os.Getenv("BIOREGISTRY_PASSWORD")
	host1 := os.Getenv("MONGOID_DATABASE_HOST1")
	host2 := os.Getenv("MONGOID_DATABASE_HOST2")
	host3 := os.Getenv("MONGOID_DATABASE_HOST3")
	db := os.Getenv("MONGOID_DATABASE")

	//uri := fmt.Sprintf("mongo:://%s:%s@%s/%s?maxPoolSize=20&w=majority", username, password, host1, db)
	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s?retryWrites=true&w=majority", username, password, host1, db)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri), options.Client().SetHosts([]string{host1, host2, host3}))
	if err != nil {
		log.Fatal("uri incorrect:", uri)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	MongoClient = client
	ch <- "connected to Mongo"
}

func CloseConnection(mongoClient *mongo.Client) {
	if mongoClient == nil {
		return
	}
	c := *mongoClient
	err := c.Disconnect(context.Background())
	if err != nil {
		log.Fatal("failed to close Milvus connection:", err.Error())
	}
}
