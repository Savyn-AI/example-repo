package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/nishantkr18/infra/migration/examples/analyse_code_nishant/config"
	"github.com/nishantkr18/infra/migration/examples/analyse_code_nishant/handlers"
	"github.com/sashabaranov/go-openai"
)

func main() {
	// Load configuration
	cfg := config.New()

	// Create handler
	handler, err := handlers.NewHandler(cfg)
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}
	defer handler.Close()

	openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Start Lambda handler
	lambda.Start(handler.Handle)
}
