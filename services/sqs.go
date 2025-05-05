package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/nishantkr18/infra/migration/examples/analyse_code_nishant/models"
)

// SQSService handles interactions with AWS SQS
type SQSService struct {
	client          *sqs.Client
	inputQueueURL   string
	outputQueueURL  string
	visibilityTimeout int32
}

// NewSQSService creates a new SQSService instance
func NewSQSService(client *sqs.Client, inputQueueURL, outputQueueURL string, visibilityTimeout int32) *SQSService {
	return &SQSService{
		client:          client,
		inputQueueURL:   inputQueueURL,
		outputQueueURL:  outputQueueURL,
		visibilityTimeout: visibilityTimeout,
	}
}

// EnqueueProcessedText sends the processed text to the output SQS queue
func (s *SQSService) EnqueueProcessedText(processedText, documentID string) error {
	message := models.SQSMessage{
		DocumentID:    documentID,
		ProcessedText: processedText,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	_, err = s.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    &s.outputQueueURL,
		MessageBody: stringPtr(string(jsonMessage)),
	})
	if err != nil {
		return fmt.Errorf("error enqueueing to SQS: %v", err)
	}

	return nil
}

// ExtendMessageVisibility extends the visibility timeout of a message
func (s *SQSService) ExtendMessageVisibility(ctx context.Context, receiptHandle string) error {
	_, err := s.client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          &s.inputQueueURL,
		ReceiptHandle:     &receiptHandle,
		VisibilityTimeout: s.visibilityTimeout,
	})
	if err != nil {
		return fmt.Errorf("error extending visibility timeout: %v", err)
	}
	return nil
}