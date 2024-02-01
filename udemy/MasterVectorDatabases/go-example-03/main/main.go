package main

import (
	"context"
	"fmt"
	"go-example-03/milvus"
	"go-example-03/mongo"
	"go-example-03/proteins"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	defer milvus.CloseConnection(milvus.MilvusClient)
	if milvus.MilvusClient == nil {
		log.Fatal("MilvusClient is nil")
	}

	defer mongo.CloseConnection(mongo.MongoClient)
	if mongo.MongoClient == nil {
		log.Fatal("MongoClient is nil")
	}

	// List Collections
	names, err := milvus.ListAllCollection(context.Background(), milvus.MilvusClient)
	if err != nil {
		log.Fatal("failed to list collections:", err.Error())
	}
	fmt.Println("collections:", strings.Join(names, ", "))

	if slices.Contains(names, "proteins") {
		err = milvus.DropCollection(context.Background(), milvus.MilvusClient, "proteins")
		if err != nil {
			log.Fatal("failed to drop collection:", err.Error())
		}
	}

	// b := books.Books{CollectionName: "books"}
	p := proteins.Proteins{CollectionName: "proteins"}

	// Create a Collection
	err = p.CreateCollection()
	if err != nil {
		log.Println("failed to create Proteins collection:", err.Error())
		os.Exit(1)
	}

	names, err = milvus.ListAllCollection(context.Background(), milvus.MilvusClient)
	if err != nil {
		log.Println("failed to list Proteins collections:", err.Error())
		os.Exit(1)
	}
	fmt.Println("collections:", strings.Join(names, ", "))

	// err = milvus.RenameCollection(context.Background(), milvus.MilvusClient, "BProteins", "books")
	// if err != nil {
	// 	log.Println("failed to rename Proteins collection:", err.Error())
	// 	os.Exit(1)
	// }

	names, err = milvus.ListAllCollection(context.Background(), milvus.MilvusClient)
	if err != nil {
		log.Println("failed to list Proteins collections:", err.Error())
		os.Exit(1)
	}
	fmt.Println("collections:", strings.Join(names, ", "))

	recordCount, err := p.CreateProteins()
	if err != nil {
		log.Println("failed to create Proteins:", err.Error())
		os.Exit(1)
	}
	log.Println("Proteins created:", recordCount)

	// err = p.DeleteProteins("")
	// if err != nil {
	// 	log.Println("failed to delete Proteins:", err.Error())
	// 	os.Exit(1)
	// }
	// log.Println("1000 Proteins deleted from original", recordCount)

	err = p.BuildIndex()
	if err != nil {
		log.Println("failed to build index:", err.Error())
		os.Exit(1)
	}
	log.Println("Index built on Proteins collection")

	// err = milvus.DropCollection(context.Background(), milvus.MilvusClient, "books")
	// if err != nil {
	// 	log.Println("failed to drop Proteins collections:", err.Error())
	// 	os.Exit(1)
	// }

	names, err = milvus.ListAllCollection(context.Background(), milvus.MilvusClient)
	if err != nil {
		log.Println("failed to list Proteins collections:", err.Error())
		os.Exit(1)
	}
	fmt.Println("collections:", strings.Join(names, ", "))

}
