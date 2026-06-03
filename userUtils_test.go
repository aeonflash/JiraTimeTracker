package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadJiraConfig_Success(t *testing.T) {
	// Create temp directory and file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, ".jirarc")
	
	config := map[string]string{
		"jira":       "test-api-key-123",
		"email":      "test@example.com",
		"graphqlUri": "https://test.atlassian.net/gateway/api/graphql",
		"cloudId":    "test-cloud-id-456",
	}
	configData, _ := json.Marshal(config)
	os.WriteFile(configFile, configData, 0644)
	
	// Mock os.UserHomeDir by setting both HOME and USERPROFILE
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("HOME", tempDir)
	os.Setenv("USERPROFILE", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	}()
	
	// Clear globals and test
	originalKey := jiraApiKey
	originalUri := jiraGraphQlBaseUri
	originalCloudId := jiraCloudId
	jiraApiKey = ""
	jiraGraphQlBaseUri = ""
	jiraCloudId = ""
	defer func() {
		jiraApiKey = originalKey
		jiraGraphQlBaseUri = originalUri
		jiraCloudId = originalCloudId
	}()
	
	loadJiraConfig()
	
	if jiraApiKey != "test-api-key-123" {
		t.Errorf("Expected jiraApiKey 'test-api-key-123', got '%s'", jiraApiKey)
	}
	if jiraGraphQlBaseUri != "https://test.atlassian.net/gateway/api/graphql" {
		t.Errorf("Expected jiraGraphQlBaseUri 'https://test.atlassian.net/gateway/api/graphql', got '%s'", jiraGraphQlBaseUri)
	}
	if jiraCloudId != "test-cloud-id-456" {
		t.Errorf("Expected jiraCloudId 'test-cloud-id-456', got '%s'", jiraCloudId)
	}
}

func TestLoadJiraConfig_FileNotFound(t *testing.T) {
	// Use non-existent directory
	tempDir := t.TempDir()
	nonExistentDir := filepath.Join(tempDir, "nonexistent")
	
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("HOME", nonExistentDir)
	os.Setenv("USERPROFILE", nonExistentDir)
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	}()
	
	originalKey := jiraApiKey
	jiraApiKey = "original-key"
	defer func() { jiraApiKey = originalKey }()
	
	loadJiraConfig()
	
	// Should remain unchanged when file not found
	if jiraApiKey != "original-key" {
		t.Errorf("Expected jiraApiKey to remain 'original-key', got '%s'", jiraApiKey)
	}
}

func TestLoadJiraConfig_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, ".jirarc")
	
	// Write invalid JSON
	os.WriteFile(configFile, []byte("{invalid json"), 0644)
	
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")
	os.Setenv("HOME", tempDir)
	os.Setenv("USERPROFILE", tempDir)
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	}()
	
	originalKey := jiraApiKey
	jiraApiKey = "original-key"
	defer func() { jiraApiKey = originalKey }()
	
	loadJiraConfig()
	
	// Should remain unchanged when JSON is invalid
	if jiraApiKey != "original-key" {
		t.Errorf("Expected jiraApiKey to remain 'original-key', got '%s'", jiraApiKey)
	}
}

func TestGetCurrentUser_Success(t *testing.T) {
	mockResponse := JiraResponse{
		Data: Data{
			Me: Me{
				User: User{
					AccountId: "test-account-id",
					Name:      "Test User",
				},
			},
		},
	}
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()
	
	// Override globals for test
	originalUri := jiraGraphQlBaseUri
	originalKey := jiraApiKey
	jiraGraphQlBaseUri = server.URL
	jiraApiKey = "test-key"
	defer func() {
		jiraGraphQlBaseUri = originalUri
		jiraApiKey = originalKey
	}()
	
	result := getCurrentUser()
	
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if result.Data.Me.User.AccountId != "test-account-id" {
		t.Errorf("Expected AccountId 'test-account-id', got '%s'", result.Data.Me.User.AccountId)
	}
	if result.Data.Me.User.Name != "Test User" {
		t.Errorf("Expected Name 'Test User', got '%s'", result.Data.Me.User.Name)
	}
}

func TestGetCurrentUser_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()
	
	originalUri := jiraGraphQlBaseUri
	originalKey := jiraApiKey
	jiraGraphQlBaseUri = server.URL
	jiraApiKey = "invalid-key"
	defer func() {
		jiraGraphQlBaseUri = originalUri
		jiraApiKey = originalKey
	}()
	
	result := getCurrentUser()
	
	if result != nil {
		t.Errorf("Expected nil result for unauthorized, got %v", result)
	}
}