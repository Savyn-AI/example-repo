package models

import "time"

// SQSEvent represents the SQS event structure
type SQSEvent struct {
	Records []struct {
		Body          string `json:"body"`
		ReceiptHandle string `json:"receiptHandle"`
	} `json:"Records"`
}

// ProcessedText represents the structure of processed text data
type ProcessedText struct {
	DocumentID    string    `bson:"_id,omitempty" json:"document_id"`
	OriginalText  string    `bson:"original_text" json:"original_text"`
	ProcessedText string    `bson:"processed_text" json:"processed_text"`
	Timestamp     time.Time `bson:"timestamp" json:"timestamp"`
	Status        string    `bson:"status" json:"status"`
}

// SQSMessage represents the structure of messages sent to SQS
type SQSMessage struct {
	DocumentID    string `json:"document_id"`
	ProcessedText string `json:"processed_text"`
	Timestamp     string `json:"timestamp"`
} 