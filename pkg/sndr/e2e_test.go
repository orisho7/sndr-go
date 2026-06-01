// go:build e2e
package sndr_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/orisho7/sndr-go/pkg/sndr"
)

// TestLiveSend attempts a real network request to the SNDR production API.
// Run with: $env:SNDR_API_KEY="your_key"; go test -v -tags=e2e ./pkg/sndr/e2e_test.go
func TestLiveSend(t *testing.T) {
	apiKey := os.Getenv("SNDR_API_KEY")
	if apiKey == "" {
		t.Skip("SNDR_API_KEY not set, skipping E2E test")
	}

	client := sndr.NewSndr(apiKey, sndr.WithTimeout(15*time.Second))

	req := &sndr.SendRequest{
		From:    "test@yourdomain.com", // Must be a verified sender in your SNDR dashboard
		To:      []string{"recipient@example.com"},
		Subject: "E2E Production SDK Test",
		Text:    "This is a live test from the Go SDK.",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.Emails.Send(ctx, req)
	if err != nil {
		t.Fatalf("Live request failed: %v", err)
	}

	if resp.ID == "" {
		t.Error("expected message ID from production API")
	}

	t.Logf("Success! Message ID: %s", resp.ID)
}
