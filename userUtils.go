package main

import (
	"bytes"
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
		jiraApiFunctions.JiraGraphQlBaseUri = jiraGraphQlBaseUri
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
	req.Header.Set("Authorization", "Bearer "+jiraApiKey)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error getting current user:", err)
		return nil
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var userResponse JiraResponse
		json.Unmarshal(body, &userResponse)
		return &userResponse
	}
	return nil
}