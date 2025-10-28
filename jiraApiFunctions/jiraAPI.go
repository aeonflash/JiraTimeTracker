package jiraApiFunctions

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	JiraGraphQlBaseUri string
	JiraApiKey         string
)

// Generic API call function
func MakeJiraAPICall(method, endpoint string, body interface{}, queryParams map[string]string) ([]byte, error) {
	baseURL := strings.Replace(JiraGraphQlBaseUri, "/gateway/api/graphql", "", 1)
	fullURL := baseURL + endpoint
	
	// Add query parameters
	if len(queryParams) > 0 {
		params := url.Values{}
		for k, v := range queryParams {
			if v != "" {
				params.Add(k, v)
			}
		}
		if len(params) > 0 {
			fullURL += "?" + params.Encode()
		}
	}
	
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}
	
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+JiraApiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	return io.ReadAll(resp.Body)
}