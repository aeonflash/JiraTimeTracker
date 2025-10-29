package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"jiraTimeWidget/jiraApiFunctions"
)

func TestGetJiraItem_Success(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if !strings.Contains(r.URL.Path, "/issue/TEST-123") {
			t.Errorf("Expected path to contain /issue/TEST-123, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", r.Header.Get("Authorization"))
		}
		
		// Return mock response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key":"TEST-123","summary":"Test Issue"}`))
	}))
	defer server.Close()

	// Override global variables for test - need to import jiraApiFunctions
	originalUri := jiraGraphQlBaseUri
	jiraGraphQlBaseUri = server.URL + "/gateway/api/graphql"
	defer func() { jiraGraphQlBaseUri = originalUri }()
	
	// Also set the package variable
	originalPackageUri := jiraApiFunctions.JiraGraphQlBaseUri
	jiraApiFunctions.JiraGraphQlBaseUri = server.URL + "/gateway/api/graphql"
	defer func() { jiraApiFunctions.JiraGraphQlBaseUri = originalPackageUri }()
	
	// Set API key
	originalKey := jiraApiFunctions.JiraApiKey
	jiraApiFunctions.JiraApiKey = "test-token"
	defer func() { jiraApiFunctions.JiraApiKey = originalKey }()

	result := getJiraItem("TEST-123", "test-token")
	
	expected := `{"key":"TEST-123","summary":"Test Issue"}`
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestGetJiraItem_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	originalUri := jiraGraphQlBaseUri
	jiraGraphQlBaseUri = server.URL + "/gateway/api/graphql"
	defer func() { jiraGraphQlBaseUri = originalUri }()
	
	originalPackageUri := jiraApiFunctions.JiraGraphQlBaseUri
	jiraApiFunctions.JiraGraphQlBaseUri = server.URL + "/gateway/api/graphql"
	defer func() { jiraApiFunctions.JiraGraphQlBaseUri = originalPackageUri }()
	
	originalKey := jiraApiFunctions.JiraApiKey
	jiraApiFunctions.JiraApiKey = "test-token"
	defer func() { jiraApiFunctions.JiraApiKey = originalKey }()

	result := getJiraItem("INVALID-123", "test-token")
	
	if result != "" {
		t.Errorf("Expected empty string for 404, got %s", result)
	}
}

func TestGetJiraItem_EmptyToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		// When token is empty, it should be "Bearer " (with space) or just "Bearer"
		if auth != "Bearer " && auth != "Bearer" {
			t.Errorf("Expected Authorization header 'Bearer ' or 'Bearer', got '%s'", auth)
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	originalUri := jiraGraphQlBaseUri
	jiraGraphQlBaseUri = server.URL + "/gateway/api/graphql"
	defer func() { jiraGraphQlBaseUri = originalUri }()
	
	originalPackageUri := jiraApiFunctions.JiraGraphQlBaseUri
	jiraApiFunctions.JiraGraphQlBaseUri = server.URL + "/gateway/api/graphql"
	defer func() { jiraApiFunctions.JiraGraphQlBaseUri = originalPackageUri }()
	
	originalKey := jiraApiFunctions.JiraApiKey
	jiraApiFunctions.JiraApiKey = ""
	defer func() { jiraApiFunctions.JiraApiKey = originalKey }()

	result := getJiraItem("TEST-123", "")
	
	if result != "" {
		t.Errorf("Expected empty string for unauthorized, got %s", result)
	}
}