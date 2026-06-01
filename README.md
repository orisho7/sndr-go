# SNDR Go SDK

A professional, production-ready Go client for the SNDR email platform. This SDK provides complete coverage of the SNDR API with a focus on type safety, concurrency, and developer experience.

## Installation

```bash
go get github.com/orisho7/sndr-go
```

## Quick Start

### 1. Secure Your API Key

Never hardcode your API key in source control. Store it in an environment variable or a `.env` file.

**Environment Variable (Shell):**
```bash
export SNDR_API_KEY="sndr_live_..."
```

**Environment Variable (PowerShell):**
```powershell
$env:SNDR_API_KEY="sndr_live_..."
```

**Using a .env file:**
Create a `.env` file in your project root:
```text
SNDR_API_KEY=sndr_live_your_key_here
```

### 2. Basic Usage

Initialize the client and send your first email. The SDK uses a nested resource pattern for intuitive navigation.

```go
package main

import (
    "context"
    "fmt"
    "os"
    "github.com/orisho7/sndr-go/pkg/sndr"
)

func main() {
    // Initialize the client using the environment variable
    apiKey := os.Getenv("SNDR_API_KEY")
    client := sndr.NewSndr(apiKey)

    // Send an email
    resp, err := client.Emails.Send(context.Background(), &sndr.SendRequest{
        From:    "hello@yourdomain.com", // Must be a verified sender
        To:      "customer@example.com",
        Subject: "Welcome aboard",
        HTML:    "<p>Thanks for joining us.</p>",
    })

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("Success! Message ID: %s\n", resp.ID)
}
```

## Features

### Emails Service

The Emails service handles all message delivery and retrieval operations.

#### Send Email (Idempotent)

Pass an optional idempotency key to prevent duplicate transmissions during retries.

```go
resp, err := client.Emails.Send(ctx, &sndr.SendRequest{
    From:       "hello@yourdomain.com",
    To:         []string{"user1@example.com", "user2@example.com"},
    Subject:    "Project Update",
    TemplateID: "tmpl_123",
    Variables:  map[string]any{"project": "Phoenix"},
}, "unique-request-id-123")
```

#### List and Retrieve

```go
// List sent emails with pagination
list, err := client.Emails.List(ctx, "optional_after_id")

// Get specific email details
email, err := client.Emails.Get(ctx, "em_123")
```

### Domains Service

Register and verify sending domains.

```go
// Register a new domain
domain, err := client.Domains.Create(ctx, &sndr.CreateDomainRequest{
    Name: "example.com",
})

// Trigger DNS verification
status, err := client.Domains.Verify(ctx, domain.ID)
```

### Templates and Resources

The SDK provides full CRUD support for all SNDR resources.

```go
// Create a versioned template
tmpl, err := client.Templates.Create(ctx, &sndr.CreateTemplateRequest{
    Name:    "Order Confirmation",
    Subject: "Your order #{{id}}",
    HTML:    "<h1>Confirmed</h1>",
})

// Generate a new API Key
key, err := client.Keys.Create(ctx, &sndr.CreateAPIKeyRequest{
    Name: "Monitoring Key",
})

// Register a Webhook
webhook, err := client.Webhooks.Create(ctx, &sndr.CreateWebhookRequest{
    URL:    "https://api.yourdomain.com/webhooks/sndr",
    Events: []string{"email.delivered", "email.bounced"},
})
```

## Error Handling

The SDK provides a structured `APIError` type for handling platform-specific errors.

```go
resp, err := client.Emails.Send(ctx, req)
if err != nil {
    if sndr.IsAPIError(err) {
        apiErr := err.(*sndr.APIError)
        // Handle specific error codes: INVALID_SENDER, RATE_LIMIT_EXCEEDED, etc.
        fmt.Printf("Code: %s, Message: %s\n", apiErr.Code, apiErr.Message)
    }
}
```

## Configuration

Configure the client using functional options.

```go
client := sndr.NewSndr(
    apiKey,
    sndr.WithTimeout(15 * time.Second),
    sndr.WithBaseURL("https://custom.sndr.sh"),
)
```

## Testing

### Unit Tests
The SDK includes a comprehensive test suite using mock servers.
```bash
go test -v ./pkg/sndr/...
```

### Production E2E Tests
To run tests against the live SNDR API:
```bash
export SNDR_API_KEY="your_live_key"
go test -v -tags=e2e ./pkg/sndr/e2e_test.go
```

## License

MIT
