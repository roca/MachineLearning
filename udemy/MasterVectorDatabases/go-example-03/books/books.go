package books

import (
	"context"
	"errors"
	"go-example-03/collections"
	"math/rand"

	"slices"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

var (
	collectionName = "books"
)
var schema = &entity.Schema{
	CollectionName: collectionName,
	Description:    "Test book search",
	Fields: []*entity.Field{
		{
			Name:       "book_id",
			DataType:   entity.FieldTypeInt64,
			PrimaryKey: true,
			AutoID:     false,
		},
		{
			Name:       "word_count",
			DataType:   entity.FieldTypeInt64,
			PrimaryKey: false,
			AutoID:     false,
		},
		{
			Name:     "book_intro",
			DataType: entity.FieldTypeFloatVector,
			TypeParams: map[string]string{
				"dim": "2",
			},
		},
	},
	EnableDynamicField: true,
}

type Books struct {
	Schema         *entity.Schema
	CollectionName string
	BookIDs        []int64
	WordCounts     []int64
	BookIntros     [][]float32
}

func (b *Books) CreateCollection() error {
	//defer collections.CloseConnection(collections.MilvusClient)
	schema.CollectionName = b.CollectionName
	b.Schema = schema

	collectionNames, _ := collections.GetCollectionNames(context.Background(), collections.MilvusClient)
	if slices.Contains(collectionNames, collectionName) {
		return errors.New("Books collection already exists!")
	}

	err := collections.CreateCollection(context.Background(), collections.MilvusClient, b.Schema)
	if err != nil {
		return err
	}
	return nil
}

// Create 2000 random books
func (b *Books) CreateBooks() (int, error) {
	//defer collections.CloseConnection(collections.MilvusClient)
	bookIDs := make([]int64, 0, 2000)
	wordCounts := make([]int64, 0, 2000)
	bookIntros := make([][]float32, 0, 2000)
	for i := 0; i < 2000; i++ {
		bookIDs = append(bookIDs, int64(i))
		wordCounts = append(wordCounts, int64(i+10000))
		v := make([]float32, 0, 2)
		for j := 0; j < 2; j++ {
			v = append(v, rand.Float32())
		}
		bookIntros = append(bookIntros, v)
	}
	idColumn := entity.NewColumnInt64("book_id", bookIDs)
	wordColumn := entity.NewColumnInt64("word_count", wordCounts)
	introColumn := entity.NewColumnFloatVector("book_intro", 2, bookIntros)

	b.BookIDs = bookIDs
	b.WordCounts = wordCounts
	b.BookIntros = bookIntros

	column, err := (*collections.MilvusClient).Insert(
		context.Background(), // ctx
		"books",              // CollectionName
		"",                   // partitionName
		idColumn,             // columnarData
		wordColumn,           // columnarData
		introColumn,          // columnarData
	)
	if err != nil {
		return -1, err
	}

	return column.Len(), nil
}

// Delete items base on a expression
func (b *Books) DeleteBooks(expr string) error {
	//defer collections.CloseConnection(collections.MilvusClient)
	err := (*collections.MilvusClient).Delete(
		context.Background(), // ctx
		"books",              // CollectionName
		"",                   // partitionName
		expr,                 // expr
	)
	if err != nil {
		return err
	}

	// Compact collection
	// This function is under active development on the GO client.

	return nil
}

// build index on the book_intro field
func (b *Books) BuildIndex() error {
	//defer collections.CloseConnection(collections.MilvusClient)
	idx, err := entity.NewIndexIvfFlat( // NewIndex func
		entity.L2, // metricType
		1024,      // ConstructParams
	)
	if err != nil {
		return err
	}

	err = (*collections.MilvusClient).CreateIndex(
		context.Background(), // ctx
		"books",              // CollectionName
		"book_intro",         // fieldName
		idx,                  // entity.Index
		false,                // async
	)
	if err != nil {
		return err
	}

	return nil

}
