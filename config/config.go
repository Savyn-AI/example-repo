package config

import (
	"os"
)

// Config holds all configuration values
type Config struct {
	OpenAIAPIKey               string
	SQSInputQueueURL          string
	SQSOutputQueueURL         string
	S3BucketName              string
	MongoDBURI                string
	MongoDBDBName             string
	MongoDBCollectionName     string
	VisibilityTimeoutExtension int
	MaxPartLength             int
}

// New creates a new Config instance with values from environment variables
func New() *Config {
	return &Config{
		OpenAIAPIKey:               os.Getenv("OPENAI_API_KEY"),
		SQSInputQueueURL:          os.Getenv("SQS_INPUT_QUEUE_URL"),
		SQSOutputQueueURL:         os.Getenv("SQS_OUTPUT_QUEUE_URL"),
		S3BucketName:              os.Getenv("S3_BUCKET_NAME"),
		MongoDBURI:                os.Getenv("MONGODB_URI"),
		MongoDBDBName:             getEnvOrDefault("MONGODB_DB_NAME", "text_processing"),
		MongoDBCollectionName:     getEnvOrDefault("MONGODB_COLLECTION_NAME", "processed_texts"),
		VisibilityTimeoutExtension: 60,
		MaxPartLength:             500,
	}
}

// Helper function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 