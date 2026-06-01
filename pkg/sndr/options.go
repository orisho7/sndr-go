package sndr

import (
	"net/http"
	"time"
)

// Option defines a functional configuration for the SNDR client.
type Option func(*Sndr)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(s *Sndr) {
		if httpClient != nil {
			s.httpClient = httpClient
		}
	}
}

// WithBaseURL overrides the default SNDR API base URL.
func WithBaseURL(url string) Option {
	return func(s *Sndr) {
		s.baseURL = url
	}
}

// WithTimeout sets a default timeout for all requests.
func WithTimeout(timeout time.Duration) Option {
	return func(s *Sndr) {
		s.httpClient.Timeout = timeout
	}
}
