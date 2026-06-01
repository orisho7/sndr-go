package sndr

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Send(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		request        *SendRequest
		idempotencyKey string
		wantErr        bool
		errCode        string
	}{
		{
			name: "Success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				// Verify headers
				if r.Header.Get("Authorization") != "Bearer test-key" {
					t.Errorf("expected bearer token, got %s", r.Header.Get("Authorization"))
				}
				if r.Header.Get("Idempotency-Key") != "idem-123" {
					t.Errorf("expected idempotency key, got %s", r.Header.Get("Idempotency-Key"))
				}

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(SendResponse{
					ID:     "msg_123",
					Status: "sent",
				})
			},
			request: &SendRequest{
				From:    "a@b.com",
				To:      []string{"c@d.com"},
				Subject: "test",
			},
			idempotencyKey: "idem-123",
			wantErr:        false,
		},
		{
			name: "API Error - 400 Bad Request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{
					Code:    "INVALID_EMAIL",
					Message: "The 'from' email is invalid",
				})
			},
			request: &SendRequest{From: "invalid"},
			wantErr: true,
			errCode: "INVALID_EMAIL",
		},
		{
			name: "Non-JSON Error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			},
			request: &SendRequest{},
			wantErr: true,
			errCode: "UNKNOWN_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Start mock server
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			// Initialize client pointing to mock server
			client := NewSndr("test-key", WithBaseURL(server.URL))

			resp, err := client.Emails.Send(context.Background(), tt.request, tt.idempotencyKey)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if apiErr, ok := err.(*APIError); ok {
					if apiErr.Code != tt.errCode {
						t.Errorf("expected error code %s, got %s", tt.errCode, apiErr.Code)
					}
				} else {
					t.Errorf("expected *APIError, got %T", err)
				}
				return
			}

			if resp.ID == "" {
				t.Error("expected message ID in response")
			}
		})
	}
}
