package main

import (
	"encoding/json"
	"fmt"
	"jiraTimeWidget/jiraApiFunctions"
	"log"
)

type RecentIssue struct {
	Key     string `json:"key"`
	Summary string `json:"summary"`
	Status  string `json:"status"`
}

func getRecentIssues(maxResults int) []RecentIssue {
	// Search for recently updated issues assigned to current user
	// Filter criteria:
	// - Assigned to current user
	// - Either not closed, OR closed but updated within last 7 days
	// - Exclude deferred and on hold issues (not actively worked on)
	// - Ordered by most recently updated first
	searchParams := map[string]string{
		"jql":        "assignee = currentUser() AND (status != Closed OR updated >= -7d) AND status NOT IN (Deferred, \"On Hold\") ORDER BY updated DESC",
		"maxResults": fmt.Sprintf("%d", maxResults),
		"fields":     "key,summary,status",
	}
	
	result, err := jiraApiFunctions.MakeJiraAPICall("GET", "/rest/api/3/search/jql", nil, searchParams)
	if err != nil {
		log.Printf("Error fetching recent issues: %v", err)
		return nil
	}
	
	var searchResult struct {
		Issues []struct {
			Key    string `json:"key"`
			Fields struct {
				Summary string `json:"summary"`
				Status  struct {
					Name string `json:"name"`
				} `json:"status"`
			} `json:"fields"`
		} `json:"issues"`
	}
	
	if err := json.Unmarshal(result, &searchResult); err != nil {
		log.Printf("Error parsing recent issues: %v", err)
		return nil
	}
	
	var recentIssues []RecentIssue
	for _, issue := range searchResult.Issues {
		recentIssues = append(recentIssues, RecentIssue{
			Key:     issue.Key,
			Summary: issue.Fields.Summary,
			Status:  issue.Fields.Status.Name,
		})
	}
	
	return recentIssues
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}