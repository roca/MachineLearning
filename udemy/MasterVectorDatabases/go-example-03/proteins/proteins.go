package proteins

import (
	"context"
	"errors"
	"fmt"
	"go-example-03/milvus"
	mg "go-example-03/mongo"
	"log"
	"os"
	"time"

	"slices"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

var (
	collectionName = "proteins"
)
var schema = &entity.Schema{
	CollectionName: collectionName,
	Description:    "Test protein search",
	Fields: []*entity.Field{
		{
			Name:       "protein_id",
			DataType:   entity.FieldTypeInt64,
			PrimaryKey: true,
			AutoID:     false,
		},
		{
			Name:     "name",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length":    "32",
				"default_value": "Unknown",
			},
			PrimaryKey: false,
			AutoID:     false,
		},
		{
			Name:     "char_count",
			DataType: entity.FieldTypeInt64,
			TypeParams: map[string]string{
				"default_value": "9999",
			},
		},
		{
			Name:     "protein_json",
			DataType: entity.FieldTypeJSON,
			TypeParams: map[string]string{
				"max_length": "16000",
			},
		},
		{
			Name:     "protein_vector",
			DataType: entity.FieldTypeFloatVector,
			TypeParams: map[string]string{
				"dim": "1536",
			},
		},
	},
	EnableDynamicField: true,
}

type Proteins struct {
	Schema          *entity.Schema
	CollectionName  string
	ProteinIDs      []int64
	Names           []string
	CharCounts      []int64
	ProteinJson     [][]byte
	ProteinVector   [][]float32
	MongoDbProteins *mongo.Cursor
}

func (p *Proteins) CreateCollection() error {
	//defer milvus.CloseConnection(milvus.MilvusClient)
	schema.CollectionName = p.CollectionName
	p.Schema = schema

	collectionNames, _ := milvus.GetCollectionNames(context.Background(), milvus.MilvusClient)
	if slices.Contains(collectionNames, collectionName) {
		return errors.New("Proteins collection already exists!")
	}

	err := milvus.CreateCollection(context.Background(), milvus.MilvusClient, p.Schema)
	if err != nil {
		return err
	}

	milvus.CreatePartition(context.Background(), milvus.MilvusClient, p.CollectionName, "Bioregistry")
	return nil
}

// Create  100 proteins from the Mongo database
func (p *Proteins) CreateProteins() (int, error) {
	//defer milvus.CloseConnection(milvus.MilvusClient)
	proteinIDs := make([]int64, 0, 100)
	names := make([]string, 0, 100)
	charCounts := make([]int64, 0, 100)
	proteinJson := make([][]byte, 0, 100)
	proteinVector := make([][]float32, 0, 100)

	// TODO: Get 100 proteins from the Mongo database
	proteinCount, err := p.GetProteins()
	if err != nil {
		return -1, err
	}
	fmt.Println("Protein count:", proteinCount)

	// proteinIDsColumn := entity.NewColumnInt64("protein_id", proteinIDs)
	// namesColumn := entity.NewColumnVarChar("name", names)
	// charCountsColumn := entity.NewColumnInt64("char_count", charCounts)
	// proteinJsonColumn := entity.NewColumnJSONBytes("protein_json", proteinJson)
	// proteinVectorColumn := entity.NewColumnFloatVector("protein_vector", 1536, proteinVector)

	p.ProteinIDs = proteinIDs
	p.Names = names
	p.CharCounts = charCounts
	p.ProteinJson = proteinJson
	p.ProteinVector = proteinVector

	// column, err := (*milvus.MilvusClient).Insert(
	// 	context.Background(), // ctx
	// 	"proteins",           // CollectionName
	// 	"",                   // partitionName
	// 	proteinIDsColumn,     // columnarData
	// 	namesColumn,          // columnarData
	// 	charCountsColumn,     // columnarData
	// 	proteinJsonColumn,    // columnarData
	// 	proteinVectorColumn,  // columnarData
	// )
	// if err != nil {
	// 	return -1, err
	// }

	_, err = (*milvus.MilvusClient).ManualCompaction(context.Background(), collectionName, time.Second*30)
	if err != nil {
		return -1, err
	}

	// return column.Len(), nil
	return 0, nil
}

// Delete items base on a expression
func (p *Proteins) DeleteProteins(expr string) error {
	//defer milvus.CloseConnection(milvus.MilvusClient)
	err := (*milvus.MilvusClient).Delete(
		context.Background(), // ctx
		"proteins",           // CollectionName
		"Bioregistry",        // partitionName
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
func (p *Proteins) BuildIndex() error {
	//defer milvus.CloseConnection(milvus.MilvusClient)
	idx, err := entity.NewIndexIvfFlat( // NewIndex func
		entity.L2, // metricType
		1024,      // ConstructParams
	)
	if err != nil {
		return err
	}

	err = (*milvus.MilvusClient).CreateIndex(
		context.Background(), // ctx
		"proteins",           // CollectionName
		"protein_vector",     // fieldName
		idx,                  // entity.Index
		false,                // async
		client.WithIndexName("protein_vector_index"),
	)
	if err != nil {
		return err
	}

	return nil

}

// Get 100 proteins from the Mongo database
func (p *Proteins) GetProteins() (int64, error) {

	db := os.Getenv("MONGOID_DATABASE")

	coll := mg.MongoClient.Database(db).Collection("proteins")

	// { "$and" : [ { "descriptive_name" : { "$ne" : None } }, { "descriptive_name" : { "$exists" : "true" } } ] }
	andFilter := bson.M{
		"$and": []bson.M{ // you can try this in []interface
			bson.M{"descriptive_name": bson.M{"$ne": "None"}},
			bson.M{"descriptive_name": bson.M{"$exists": "true"}},
		},
	}
	_ = andFilter

	count, err := coll.CountDocuments(context.TODO(), andFilter)
	if err != nil {
		log.Fatal(err)
	}

	return count, nil

}
