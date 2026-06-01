package sndr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.sndr.sh"
	headerAuth     = "Authorization"
	headerIdem     = "Idempotency-Key"
)

// Sndr handles communication with the SNDR API.
type Sndr struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	// Emails provides access to email-related operations.
	Emails *EmailsService
	// Domains provides access to domain-related operations.
	Domains *DomainsService
	// Templates provides access to template-related operations.
	Templates *TemplatesService
	// Keys provides access to API key operations.
	Keys *KeysService
	// Webhooks provides access to webhook operations.
	Webhooks *WebhooksService
}

// EmailsService provides methods for email operations.
type EmailsService struct {
	client *Sndr
}

// DomainsService provides methods for domain operations.
type DomainsService struct {
	client *Sndr
}

// TemplatesService provides methods for template operations.
type TemplatesService struct {
	client *Sndr
}

// KeysService provides methods for API key operations.
type KeysService struct {
	client *Sndr
}

// WebhooksService provides methods for webhook operations.
type WebhooksService struct {
	client *Sndr
}

// NewSndr returns a new SNDR client with the provided API key and options.
func NewSndr(apiKey string, opts ...Option) *Sndr {
	s := &Sndr{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	s.Emails = &EmailsService{client: s}
	s.Domains = &DomainsService{client: s}
	s.Templates = &TemplatesService{client: s}
	s.Keys = &KeysService{client: s}
	s.Webhooks = &WebhooksService{client: s}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// List retrieves a paginated list of sent emails.
func (s *EmailsService) List(ctx context.Context, after ...string) (*ListEmailsResponse, error) {
	url := fmt.Sprintf("%s/v1/emails", strings.TrimSuffix(s.client.baseURL, "/"))
	if len(after) > 0 && after[0] != "" {
		url = fmt.Sprintf("%s?after=%s", url, after[0])
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set(headerAuth, fmt.Sprintf("Bearer %s", s.client.apiKey))

	resp, err := s.client.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, s.client.handleError(resp)
	}

	var listResp ListEmailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &listResp, nil
}

// Send sends an email via the SNDR API.
func (s *EmailsService) Send(ctx context.Context, req *SendRequest, idempotencyKey ...string) (*SendResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("send request cannot be nil")
	}

	// Normalize 'To' field: ensure it's always a slice of strings for the API
	var to []string
	switch v := req.To.(type) {
	case string:
		to = []string{v}
	case []string:
		to = v
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok {
				to = append(to, str)
			}
		}
	default:
		return nil, fmt.Errorf("invalid 'to' field type: expected string or []string")
	}

	// Create a copy to avoid mutating the original request struct
	normalizedReq := *req
	normalizedReq.To = to

	body, err := json.Marshal(normalizedReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/send", strings.TrimSuffix(s.client.baseURL, "/"))
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set(headerAuth, fmt.Sprintf("Bearer %s", s.client.apiKey))

	if len(idempotencyKey) > 0 && idempotencyKey[0] != "" {
		httpReq.Header.Set(headerIdem, idempotencyKey[0])
	}

	resp, err := s.client.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, s.client.handleError(resp)
	}

	var sendResp SendResponse
	if err := json.NewDecoder(resp.Body).Decode(&sendResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &sendResp, nil
}

// Get retrieves a specific email by ID.
func (s *EmailsService) Get(ctx context.Context, id string) (*EmailResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("email id cannot be empty")
	}

	url := fmt.Sprintf("%s/v1/emails/%s", strings.TrimSuffix(s.client.baseURL, "/"), id)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set(headerAuth, fmt.Sprintf("Bearer %s", s.client.apiKey))

	resp, err := s.client.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, s.client.handleError(resp)
	}

	var emailResp EmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&emailResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &emailResp, nil
}

// Create adds a sending domain and returns required DNS records.
func (s *DomainsService) Create(ctx context.Context, req *CreateDomainRequest) (*DomainResponse, error) {
	return performPost[*CreateDomainRequest, *DomainResponse](ctx, s.client, "/v1/domains", req)
}

// Verify triggers DNS verification for a domain.
func (s *DomainsService) Verify(ctx context.Context, id string) (*DomainResponse, error) {
	return performPost[any, *DomainResponse](ctx, s.client, fmt.Sprintf("/v1/domains/%s/verify", id), nil)
}

// Create creates a new versioned template.
func (s *TemplatesService) Create(ctx context.Context, req *CreateTemplateRequest) (*TemplateResponse, error) {
	return performPost[*CreateTemplateRequest, *TemplateResponse](ctx, s.client, "/v1/templates", req)
}

// Create generates a new API key. The full key is returned exactly once.
func (s *KeysService) Create(ctx context.Context, req *CreateAPIKeyRequest) (*APIKeyResponse, error) {
	return performPost[*CreateAPIKeyRequest, *APIKeyResponse](ctx, s.client, "/v1/api-keys", req)
}

// Create registers a new webhook endpoint and returns the signing secret.
func (s *WebhooksService) Create(ctx context.Context, req *CreateWebhookRequest) (*WebhookResponse, error) {
	return performPost[*CreateWebhookRequest, *WebhookResponse](ctx, s.client, "/v1/webhooks", req)
}

// Generic helper for POST requests to reduce boilerplate
func performPost[T any, R any](ctx context.Context, client *Sndr, path string, payload T) (R, error) {
	var empty R
	body, err := json.Marshal(payload)
	if err != nil {
		return empty, fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s%s", strings.TrimSuffix(client.baseURL, "/"), path)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return empty, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set(headerAuth, fmt.Sprintf("Bearer %s", client.apiKey))

	resp, err := client.httpClient.Do(httpReq)
	if err != nil {
		return empty, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return empty, client.handleError(resp)
	}

	var result R
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return empty, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}

func (s *Sndr) handleError(resp *http.Response) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read error response body: %w", err)
	}

	var errRes ErrorResponse
	if err := json.Unmarshal(data, &errRes); err != nil {
		// Fallback for non-JSON errors
		return &APIError{
			StatusCode: resp.StatusCode,
			Code:       "UNKNOWN_ERROR",
			Message:    string(data),
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Code:       errRes.Code,
		Message:    errRes.Message,
		Fields:     errRes.Fields,
	}
}
