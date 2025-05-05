package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Service handles interactions with AWS S3
type S3Service struct {
	client     *s3.Client
	bucketName string
}

// NewS3Service creates a new S3Service instance
func NewS3Service(client *s3.Client, bucketName string) *S3Service {
	return &S3Service{
		client:     client,
		bucketName: bucketName,
	}
}

// StoreProcessedText stores the processed text in S3
func (s *S3Service) StoreProcessedText(originalText, processedText, documentID string) (string, error) {
	timestamp := time.Now().UTC().Format("20060102_150405")
	key := fmt.Sprintf("processed_texts/%s_%s.json", documentID, timestamp)

	content := map[string]interface{}{
		"document_id":    documentID,
		"original_text":  originalText,
		"processed_text": processedText,
		"timestamp":      timestamp,
		"status":         "completed",
	}

	jsonContent, err := json.Marshal(content)
	if err != nil {
		return "", fmt.Errorf("error marshaling content: %v", err)
	}

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &s.bucketName,
		Key:         &key,
		Body:        strings.NewReader(string(jsonContent)),
		ContentType: stringPtr("application/json"),
	})
	if err != nil {
		return "", fmt.Errorf("error storing in S3: %v", err)
	}

	return key, nil
}

// Helper function to convert string to pointer
func stringPtr(s string) *string {
	return &s
} 