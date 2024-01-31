package main

import (
	"context"
	"fmt"
	"go-example-03/collections"
	"go-example-03/proteins"
	"log"
	"os"
	"slices"
	"strings"
)

func main() {
	defer collections.CloseConnection(collections.MilvusClient)
	if collections.MilvusClient == nil {
		log.Fatal("MilvusClient is nil")
	}

	// List Collections
	names, err := collections.ListAllCollection(context.Background(), collections.MilvusClient)
	if err != nil {
		log.Fatal("failed to list collections:", err.Error())
	}
	fmt.Println("collections:", strings.Join(names, ", "))

	if slices.Contains(names, "proteins") {
		err = collections.DropCollection(context.Background(), collections.MilvusClient, "proteins")
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

	names, err = collections.ListAllCollection(context.Background(), collections.MilvusClient)
	if err != nil {
		log.Println("failed to list Proteins collections:", err.Error())
		os.Exit(1)
	}
	fmt.Println("collections:", strings.Join(names, ", "))

	// err = collections.RenameCollection(context.Background(), collections.MilvusClient, "BProteins", "books")
	// if err != nil {
	// 	log.Println("failed to rename Proteins collection:", err.Error())
	// 	os.Exit(1)
	// }

	names, err = collections.ListAllCollection(context.Background(), collections.MilvusClient)
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

	// err = collections.DropCollection(context.Background(), collections.MilvusClient, "books")
	// if err != nil {
	// 	log.Println("failed to drop Proteins collections:", err.Error())
	// 	os.Exit(1)
	// }

	names, err = collections.ListAllCollection(context.Background(), collections.MilvusClient)
	if err != nil {
		log.Println("failed to list Proteins collections:", err.Error())
		os.Exit(1)
	}
	fmt.Println("collections:", strings.Join(names, ", "))

}
