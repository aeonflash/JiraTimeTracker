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
	
	config := map[string]string{"jira": "test-api-key-123"}
	configData, _ := json.Marshal(config)
	os.WriteFile(configFile, configData, 0644)
	
	// Mock os.UserHomeDir
	originalHomeDir := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHomeDir)
	
	// Clear jiraApiKey and test
	originalKey := jiraApiKey
	jiraApiKey = ""
	defer func() { jiraApiKey = originalKey }()
	
	loadJiraConfig()
	
	if jiraApiKey != "test-api-key-123" {
		t.Errorf("Expected jiraApiKey 'test-api-key-123', got '%s'", jiraApiKey)
	}
}

func TestLoadJiraConfig_FileNotFound(t *testing.T) {
	// Use non-existent directory
	tempDir := t.TempDir()
	nonExistentDir := filepath.Join(tempDir, "nonexistent")
	
	originalHomeDir := os.Getenv("HOME")
	os.Setenv("HOME", nonExistentDir)
	defer os.Setenv("HOME", originalHomeDir)
	
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
	
	originalHomeDir := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHomeDir)
	
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