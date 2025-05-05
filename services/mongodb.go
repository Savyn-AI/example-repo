package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nishantkr18/infra/migration/examples/analyse_code_nishant/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBService handles interactions with MongoDB
type MongoDBService struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoDBService creates a new MongoDBService instance
func NewMongoDBService(uri, dbName, collectionName string) (*MongoDBService, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDBService{
		client:     client,
		collection: collection,
	}, nil
}

// StoreProcessedText stores the processed text in MongoDB
func (s *MongoDBService) StoreProcessedText(originalText, processedText string) (string, error) {
	document := models.ProcessedText{
		OriginalText:  originalText,
		ProcessedText: processedText,
		Timestamp:     time.Now().UTC(),
		Status:        "completed",
	}

	result, err := s.collection.InsertOne(context.TODO(), document)
	if err != nil {
		return "", fmt.Errorf("error storing in MongoDB: %v", err)
	}

	return result.InsertedID.(string), nil
}

// Close closes the MongoDB connection
func (s *MongoDBService) Close() error {
	return s.client.Disconnect(context.TODO())
} 