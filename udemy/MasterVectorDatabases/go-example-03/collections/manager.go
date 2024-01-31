package collections

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

var MilvusClient *client.Client

func init() {

	ch := make(chan string, 1)

	// Create a context with a timeout of 5 seconds
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// Start the ConnectToMilvus function
	go ConnectToMilvus(ctxTimeout, ch)

	select {
	case <-ctxTimeout.Done():
		fmt.Printf("Context cancelled: %v\n", ctxTimeout.Err())
	case result := <-ch:
		fmt.Printf("Received: %s\n", result)
	}

}

func ConnectToMilvus(ctx context.Context, ch chan string) {
	// NewGrpcClient
	milvusClient, err := client.NewGrpcClient(
		ctx,               // ctx
		"localhost:19530", // addr
	)
	if err != nil {
		ch <- fmt.Sprint("failed to connect to Milvus:", err.Error())
	}

	milvusClient.UsingDatabase(ctx, "bioregistry")

	MilvusClient = &milvusClient
	ch <- "connected to Milvus"
}

func CloseConnection(milvusClient *client.Client) {
	if milvusClient == nil {
		return
	}
	c := *milvusClient
	err := c.Close()
	if err != nil {
		log.Fatal("failed to close Milvus connection:", err.Error())
	}
}

func ListAllCollection(ctx context.Context, client *client.Client) ([]string, error) {
	collections, err := (*client).ListCollections(ctx)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, collection := range collections {
		names = append(names, collection.Name)
	}

	return names, nil
}

func GetCollectionNames(ctx context.Context, client *client.Client) ([]string, error) {
	collections, err := (*client).ListCollections(ctx)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, collection := range collections {
		names = append(names, collection.Name)
	}
	return names, nil
}

func CreateCollection(ctx context.Context, client *client.Client, schema *entity.Schema) error {
	err := (*client).CreateCollection(
		ctx, // ctx
		schema,
		2, // shardNum
	)
	if err != nil {
		return err
	}
	return nil
}

func DropCollection(ctx context.Context, client *client.Client, collectionName string) error {
	err := (*client).DropCollection(ctx, collectionName)
	if err != nil {
		return err
	}
	return nil
}

func RenameCollection(ctx context.Context, client *client.Client, collectionName string, newName string) error {
	err := (*client).RenameCollection(ctx, collectionName, newName)
	if err != nil {
		return err
	}
	return nil
}

func CreatePartition(ctx context.Context, client *client.Client, collectionName string, partitionName string) error {
	err := (*client).CreatePartition(ctx, collectionName, partitionName)
	if err != nil {
		return err
	}
	return nil
}

func CompactCollection(ctx context.Context, client *client.Client, collectionName string) error {
	_, err := (*client).ManualCompaction(ctx, collectionName, time.Second*30)
	if err != nil {
		return err
	}
	return nil
}
