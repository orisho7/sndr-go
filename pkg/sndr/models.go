package sndr

// SendRequest defines the payload for the /v1/emails endpoint.
type SendRequest struct {
	From        string            `json:"from"`
	To          any               `json:"to"` // Supports string or []string
	Subject     string            `json:"subject"`
	HTML        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	TemplateID  string            `json:"template_id,omitempty"`
	Variables   map[string]any    `json:"variables,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
}

// Attachment defines an email attachment.
type Attachment struct {
	Filename string `json:"filename"`
	Content  string `json:"content"` // Base64 encoded
	Type     string `json:"type"`    // MIME type
}

// SendResponse defines the successful response from the /v1/emails endpoint.
type SendResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// ErrorResponse defines the structured error format from the SNDR API.
type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// EmailResponse defines the detailed information for a single email.
type EmailResponse struct {
	ID        string            `json:"id"`
	From      string            `json:"from"`
	To        []string          `json:"to"`
	Subject   string            `json:"subject"`
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Tags      map[string]string `json:"tags,omitempty"`
}

// ListEmailsResponse defines the paginated list of sent emails.
type ListEmailsResponse struct {
	Data []EmailResponse `json:"data"`
	Meta struct {
		HasMore bool   `json:"has_more"`
		Next    string `json:"next,omitempty"`
	} `json:"meta"`
}

// CreateDomainRequest defines the payload to add a domain.
type CreateDomainRequest struct {
	Name string `json:"name"`
}

// DNSRecord defines a DNS record returned by the domain creation.
type DNSRecord struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Priority int    `json:"priority,omitempty"`
	Status   string `json:"status"`
}

// DomainResponse defines the response for a created domain.
type DomainResponse struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	DNSRecords []DNSRecord `json:"dns_records"`
	Status     string      `json:"status"`
}

// CreateTemplateRequest defines the payload to create a template.
type CreateTemplateRequest struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	HTML    string `json:"html"`
	Text    string `json:"text,omitempty"`
}

// TemplateResponse defines the response for a created template.
type TemplateResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version int    `json:"version"`
}

// CreateAPIKeyRequest defines the payload to create an API key.
type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

// APIKeyResponse defines the response for a created API key.
type APIKeyResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"` // Full key returned once
}

// CreateWebhookRequest defines the payload to register a webhook.
type CreateWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

// WebhookResponse defines the response for a registered webhook.
type WebhookResponse struct {
	ID            string   `json:"id"`
	URL           string   `json:"url"`
	Events        []string `json:"events"`
	SigningSecret string   `json:"signing_secret"`
}
