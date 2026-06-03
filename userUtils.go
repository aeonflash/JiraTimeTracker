package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"jiraTimeWidget/jiraApiFunctions"
)

func loadJiraConfig() {
	homeDir, _ := os.UserHomeDir()
	data, err := os.ReadFile(homeDir + "/.jirarc")
	if err != nil {
		log.Println("Error reading .jirarc:", err)
		return
	}
	var config map[string]interface{}
	json.Unmarshal(data, &config)
	if jira, ok := config["jira"].(string); ok {
		jiraApiKey = jira
		jiraApiFunctions.JiraApiKey = jira
	}
	// Load Jira base URI
	if uri, ok := config["graphqlUri"].(string); ok {
		jiraGraphQlBaseUri = uri
		jiraApiFunctions.JiraGraphQlBaseUri = uri
	}
	// Load Cloud ID
	if cloudId, ok := config["cloudId"].(string); ok {
		jiraCloudId = cloudId
	}
	// Also check for email if provided
	if email, ok := config["email"].(string); ok {
		jiraApiFunctions.JiraEmail = email
	}
}

func getCurrentUser() *JiraResponse {
	query := map[string]string{"query": CurrentUserQuery}
	queryData, _ := json.Marshal(query)
	
	req, err := http.NewRequest("POST", jiraGraphQlBaseUri, bytes.NewBuffer(queryData))
	if err != nil {
		log.Println("Error building request:", err)
		return nil
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	// Use the same authentication method as the REST API
	if jiraApiFunctions.JiraEmail != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(jiraApiFunctions.JiraEmail + ":" + jiraApiKey))
		req.Header.Set("Authorization", "Basic "+auth)
	} else {
		req.Header.Set("Authorization", "Bearer "+jiraApiKey)
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error getting current user:", err)
		return nil
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode == http.StatusOK {
		var userResponse JiraResponse
		if err := json.Unmarshal(body, &userResponse); err != nil {
			log.Println("Error unmarshaling user response:", err)
			log.Println("Response body:", string(body))
			return nil
		}
		return &userResponse
	}
	
	log.Printf("Failed to get current user. Status: %d, Body: %s\n", resp.StatusCode, string(body))
	return nil
}