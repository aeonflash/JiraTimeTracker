package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateJiraURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
		errMsg  string
	}{
		{"empty string", "", true, "Jira URL is required"},
		{"no scheme", "example.atlassian.net", true, "URL must start with https://"},
		{"http not https", "http://example.atlassian.net", true, "URL must start with https://"},
		{"https no host", "https://", true, "URL must include a hostname"},
		{"valid atlassian URL", "https://myorg.atlassian.net", false, ""},
		{"valid with trailing slash", "https://myorg.atlassian.net/", false, ""},
		{"valid custom domain", "https://jira.mycompany.com", false, ""},
		{"just random text", "not a url at all", true, "URL must start with https://"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateJiraURL(tt.url)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errMsg)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %q", err.Error())
				}
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
		errMsg  string
	}{
		{"empty string", "", true, "email is required"},
		{"no at sign", "notanemail", true, "invalid email format"},
		{"no domain", "user@", true, "invalid email format"},
		{"no tld", "user@domain", true, "invalid email format"},
		{"double at", "user@@domain.com", true, "invalid email format"},
		{"valid simple", "user@example.com", false, ""},
		{"valid with dots", "first.last@company.co.uk", false, ""},
		{"valid with plus", "user+tag@example.com", false, ""},
		{"valid with dash", "user-name@my-company.com", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errMsg)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %q", err.Error())
				}
			}
		})
	}
}

func TestValidateAPIToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
		errMsg  string
	}{
		{"empty string", "", true, "API token is required"},
		{"too short", "abc", true, "API token seems too short — check that you pasted the full token"},
		{"9 chars", "123456789", true, "API token seems too short — check that you pasted the full token"},
		{"10 chars", "1234567890", false, ""},
		{"realistic token", "ATATT3xFfGF0ABCDEFGHIJ1234567890abcdefghij", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAPIToken(tt.token)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errMsg)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %q", err.Error())
				}
			}
		})
	}
}

func TestTestJiraConnection_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify it hits the right endpoint
		if r.URL.Path != "/rest/api/3/myself" {
			t.Errorf("expected path /rest/api/3/myself, got %s", r.URL.Path)
		}
		// Verify auth header is present
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"accountId":"123","displayName":"Test User"}`))
	}))
	defer server.Close()

	err := testJiraConnection(server.URL, "user@example.com", "valid-token-12345")
	if err != nil {
		t.Errorf("expected successful connection, got error: %v", err)
	}
}

func TestTestJiraConnection_AuthFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Unauthorized"}`))
	}))
	defer server.Close()

	err := testJiraConnection(server.URL, "user@example.com", "bad-token-12345")
	if err == nil {
		t.Error("expected error for 401, got nil")
		return
	}
	expected := "authentication failed — check your email and API token"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestTestJiraConnection_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"message":"Forbidden"}`))
	}))
	defer server.Close()

	err := testJiraConnection(server.URL, "user@example.com", "valid-token-12345")
	if err == nil {
		t.Error("expected error for 403, got nil")
		return
	}
	expected := "access forbidden — your token may lack permissions"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestTestJiraConnection_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`Not Found`))
	}))
	defer server.Close()

	err := testJiraConnection(server.URL, "user@example.com", "valid-token-12345")
	if err == nil {
		t.Error("expected error for 404, got nil")
		return
	}
	expected := "Jira API not found — check your URL is correct"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestTestJiraConnection_NetworkError(t *testing.T) {
	// Use a URL that won't connect
	err := testJiraConnection("https://localhost:1", "user@example.com", "valid-token-12345")
	if err == nil {
		t.Error("expected network error, got nil")
		return
	}
	expected := "connection failed — check your URL and internet connection"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestJiraConfigExists(t *testing.T) {
	// This test just confirms the function doesn't panic
	// Result depends on whether ~/.jirarc exists in the test environment
	_ = jiraConfigExists()
}
