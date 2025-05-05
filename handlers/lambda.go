package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/nishantkr18/infra/migration/examples/analyse_code_nishant/config"
	"github.com/nishantkr18/infra/migration/examples/analyse_code_nishant/models"
	"github.com/nishantkr18/infra/migration/examples/analyse_code_nishant/services"
)

// Handler handles SQS events, processes text with OpenAI in parts,
// stores results in MongoDB and S3, and enqueues to output SQS
type Handler struct {
	cfg          *config.Config
	openaiSvc    *services.OpenAIService
	mongoSvc     *services.MongoDBService
	s3Svc        *services.S3Service
	sqsSvc       *services.SQSService
}

// NewHandler creates a new Handler instance
func NewHandler(cfg *config.Config) (*Handler, error) {
	// Initialize AWS clients
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(awsCfg)
	s3Client := s3.NewFromConfig(awsCfg)

	// Initialize services
	openaiSvc := services.NewOpenAIService(cfg.OpenAIAPIKey)
	
	mongoSvc, err := services.NewMongoDBService(cfg.MongoDBURI, cfg.MongoDBDBName, cfg.MongoDBCollectionName)
	if err != nil {
		return nil, err
	}

	s3Svc := services.NewS3Service(s3Client, cfg.S3BucketName)
	sqsSvc := services.NewSQSService(sqsClient, cfg.SQSInputQueueURL, cfg.SQSOutputQueueURL, int32(cfg.VisibilityTimeoutExtension))

	return &Handler{
		cfg:          cfg,
		openaiSvc:    openaiSvc,
		mongoSvc:     mongoSvc,
		s3Svc:        s3Svc,
		sqsSvc:       sqsSvc,
	}, nil
}

// Handle processes the SQS event
func (h *Handler) Handle(ctx context.Context, event models.SQSEvent) (map[string]interface{}, error) {
	for _, record := range event.Records {
		var message map[string]interface{}
		if err := json.Unmarshal([]byte(record.Body), &message); err != nil {
			return nil, fmt.Errorf("error unmarshaling message: %v", err)
		}

		textToProcess, ok := message["text"].(string)
		if !ok || textToProcess == "" {
			log.Println("No 'text' key found in the SQS message.")
			continue
		}

		receiptHandle := record.ReceiptHandle
		textParts := h.openaiSvc.SplitText(textToProcess, h.cfg.MaxPartLength)
		var results []string

		for i, part := range textParts {
			log.Printf("Processing part %d/%d", i+1, len(textParts))
			openaiResult, err := h.openaiSvc.ProcessText(part)
			if err != nil {
				log.Printf("Failed to process part %d. Retrying later: %v", i+1, err)
				return nil, fmt.Errorf("OpenAI processing failed for a part: %v", err)
			}
			results = append(results, openaiResult)

			// Extend visibility timeout
			log.Printf("Extending visibility timeout for message: %s", receiptHandle)
			if err := h.sqsSvc.ExtendMessageVisibility(ctx, receiptHandle); err != nil {
				log.Printf("Error extending visibility timeout: %v", err)
			} else {
				log.Printf("Visibility timeout extended by %d seconds.", h.cfg.VisibilityTimeoutExtension)
			}
			time.Sleep(1 * time.Second) // Be mindful of API rate limits
		}

		finalResult := strings.Join(results, " ")
		log.Printf("Successfully processed all parts. Final result: %s", finalResult)

		// Store in MongoDB
		documentID, err := h.mongoSvc.StoreProcessedText(textToProcess, finalResult)
		if err != nil {
			return nil, err
		}
		log.Printf("Stored in MongoDB with ID: %s", documentID)

		// Store in S3
		s3Key, err := h.s3Svc.StoreProcessedText(textToProcess, finalResult, documentID)
		if err != nil {
			return nil, err
		}
		log.Printf("Stored in S3 with key: %s", s3Key)

		// Enqueue to output SQS
		if err := h.sqsSvc.EnqueueProcessedText(finalResult, documentID); err != nil {
			return nil, err
		}
		log.Println("Enqueued to output SQS queue")
	}

	return map[string]interface{}{
		"statusCode": 200,
		"body":       map[string]string{"message": "Successfully processed SQS messages"},
	}, nil
}

// Close closes all service connections
func (h *Handler) Close() error {
	return h.mongoSvc.Close()
} 