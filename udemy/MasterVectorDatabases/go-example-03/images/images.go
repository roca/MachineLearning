package images

import (
	"context"
	"errors"
	"math/rand"
	"slices"
	"go-example-03/collections"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

var (
	collectionName = "images"
)

var schema = &entity.Schema{
	CollectionName: collectionName,
	Description:    "Image book search",
	Fields: []*entity.Field{
		{
			Name:       "image_id",
			DataType:   entity.FieldTypeInt64,
			PrimaryKey: true,
			AutoID:     false,
		},
		{
			Name:     "image",
			DataType: entity.FieldTypeFloatVector,
			TypeParams: map[string]string{
				"dim": "2",
			},
		},
	},
	EnableDynamicField: true,
}

type Images struct {
	Schema         *entity.Schema
	CollectionName string
	ImageIDs       []int64
	Images         [][]float32
}

func (i *Images) CreateCollection() error {
	//defer collections.CloseConnection(collections.MilvusClient)
	schema.CollectionName = i.CollectionName
	i.Schema = schema

	collectionNames, _ := collections.GetCollectionNames(context.Background(), collections.MilvusClient)
	if slices.Contains(collectionNames, collectionName) {
		return errors.New("Images collection already exists!")
	}

	err := collections.CreateCollection(context.Background(), collections.MilvusClient, i.Schema)
	if err != nil {
		return err
	}
	return nil
}

func (i *Images) CreateImages() (int, error) {
	//defer collections.CloseConnection(collections.MilvusClient)
	imageIDs := make([]int64, 0, 2000)
	images := make([][]float32, 0, 2000)
	for i := 0; i < 2000; i++ {
		imageIDs = append(imageIDs, int64(i))
		v := make([]float32, 0, 2)
		for j := 0; j < 2; j++ {
			v = append(v, rand.Float32())
		}
		images = append(images, v)
	}
	idColumn := entity.NewColumnInt64("image_id", imageIDs)
	imageColumn := entity.NewColumnFloatVector("image", 2, images)

	i.ImageIDs = imageIDs
	i.Images = images

	column, err := (*collections.MilvusClient).Insert(
		context.Background(), // ctx
		"images",             // CollectionName
		"",                   // partitionName
		idColumn,             // columnarData
		imageColumn,          // columnarData
	)
	if err != nil {
		return -1, err
	}

	return column.Len(), nil
}

// Delete items base on a expression
func (i *Images) DeleteBooks(expr string) error {
	//defer collections.CloseConnection(collections.MilvusClient)
	err := (*collections.MilvusClient).Delete(
		context.Background(), // ctx
		"images",             // CollectionName
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
func (i *Images) BuildIndex() error {
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
		"images",             // CollectionName
		"image",              // fieldName
		idx,                  // entity.Index
		false,                // async
	)
	if err != nil {
		return err
	}

	return nil

}
