package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/orisho7/sndr-go/pkg/sndr"
)

func main() {
	// 1. Initialize Sndr with functional options
	apiKey := os.Getenv("SNDR_API_KEY")
	if apiKey == "" {
		apiKey = "test_key_123" // Placeholder for demonstration
	}

	sndrClient := sndr.NewSndr(
		apiKey,
		sndr.WithTimeout(10*time.Second),
	)

	// 2. Prepare Request
	req := &sndr.SendRequest{
		From:    "hello@example.com",
		To:      []string{"user@example.com"},
		Subject: "Production SDK Test",
		HTML:    "<h1>Success</h1><p>The SNDR Go SDK is working.</p>",
		Text:    "Success: The SNDR Go SDK is working.",
		Tags: map[string]string{
			"env":     "production",
			"version": "1.0.0",
		},
	}

	// 3. Execute with Context and Idempotency Key
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Using a unique ID for idempotency (prevents duplicate sends on retries)
	idempotencyKey := "req_abc_123"

	resp, err := sndrClient.Emails.Send(ctx, req, idempotencyKey)
	if err != nil {
		if sndr.IsAPIError(err) {
			apiErr := err.(*sndr.APIError)
			log.Fatalf("SNDR API Error: [%s] %s (Status: %d)", apiErr.Code, apiErr.Message, apiErr.StatusCode)
		}
		log.Fatalf("Network/Request Error: %v", err)
	}

	// 4. Output Results
	fmt.Printf("Email Sent Successfully!\n")
	fmt.Printf("Message ID: %s\n", resp.ID)
	fmt.Printf("Status:     %s\n", resp.Status)

	// 5. Demonstrate looking up the message
	fmt.Printf("\nRetrieving message details...\n")
	msg, err := sndrClient.Emails.Get(ctx, resp.ID)
	if err != nil {
		log.Printf("Warning: Could not retrieve message details: %v", err)
	} else {
		fmt.Printf("Retrieved Subject: %s\n", msg.Subject)
		fmt.Printf("Retrieved Status:  %s\n", msg.Status)
	}
}
