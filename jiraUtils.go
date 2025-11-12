package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"jiraTimeWidget/jiraApiFunctions"
)

func getJiraItem(jiraId string, token string) string {
	jiraItem, err := jiraApiFunctions.GetIssue(jiraId, "", "names,renderedFields")
	if err != nil {
		log.Println("Error getting jira item:", err)
		return ""
	}
	return string(jiraItem)
}

func logWorkToJira(jiraId string, timeSpent string, comment string, startTime time.Time) error {
	worklogData := map[string]interface{}{
		"timeSpent": timeSpent,
		"comment": map[string]interface{}{
			"type": "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{
							"text": comment,
							"type": "text",
						},
					},
				},
			},
		},
		"started": startTime.Format("2006-01-02T15:04:05.000-0700"),
	}
	
	_, err := jiraApiFunctions.AddWorklog(jiraId, worklogData)
	return err
}

// GetIssueStatus retrieves the current status of a Jira issue
func GetIssueStatus(issueKey string) (*StatusInfo, error) {
	// Call GetIssue API with fields parameter set to "status"
	response, err := jiraApiFunctions.GetIssue(issueKey, "status", "")
	if err != nil {
		log.Printf("Error fetching issue status for %s: %v", issueKey, err)
		return nil, err
	}

	// Parse JSON response to extract status information
	var issueData struct {
		Fields struct {
			Status struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
				StatusCategory struct {
					Key string `json:"key"`
				} `json:"statusCategory"`
			} `json:"status"`
		} `json:"fields"`
	}

	if err := json.Unmarshal(response, &issueData); err != nil {
		log.Printf("Error parsing status response for %s: %v", issueKey, err)
		return nil, err
	}

	// Check if status data is present
	if issueData.Fields.Status.ID == "" {
		log.Printf("No status found in response for %s", issueKey)
		return nil, fmt.Errorf("no status found for issue %s", issueKey)
	}

	// Create and return StatusInfo struct
	statusInfo := &StatusInfo{
		ID:          issueData.Fields.Status.ID,
		Name:        issueData.Fields.Status.Name,
		Description: issueData.Fields.Status.Description,
		Category:    issueData.Fields.Status.StatusCategory.Key,
	}

	return statusInfo, nil
}

// GetAvailableTransitions retrieves all available status transitions for a Jira issue
func GetAvailableTransitions(issueKey string) ([]Transition, error) {
	// Get current status to determine direction
	currentStatus, err := GetIssueStatus(issueKey)
	if err != nil {
		log.Printf("Error fetching current status for %s: %v", issueKey, err)
		return nil, err
	}

	// Call GetIssueTransitions API
	response, err := jiraApiFunctions.GetIssueTransitions(issueKey, "")
	if err != nil {
		log.Printf("Error fetching transitions for %s: %v", issueKey, err)
		return nil, err
	}

	// Parse JSON response to extract available transitions
	var transitionsData struct {
		Transitions []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			To   struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
				StatusCategory struct {
					Key string `json:"key"`
				} `json:"statusCategory"`
			} `json:"to"`
		} `json:"transitions"`
	}

	if err := json.Unmarshal(response, &transitionsData); err != nil {
		log.Printf("Error parsing transitions response for %s: %v", issueKey, err)
		return nil, err
	}

	// Convert to Transition structs with IsForward field populated
	transitions := make([]Transition, 0, len(transitionsData.Transitions))
	for _, t := range transitionsData.Transitions {
		transition := Transition{
			ID:   t.ID,
			Name: t.Name,
			To: StatusInfo{
				ID:          t.To.ID,
				Name:        t.To.Name,
				Description: t.To.Description,
				Category:    t.To.StatusCategory.Key,
			},
			IsForward: determineTransitionDirection(currentStatus.Category, t.To.StatusCategory.Key, t.To.Name),
		}
		transitions = append(transitions, transition)
	}

	return transitions, nil
}

// determineTransitionDirection determines if a transition is forward or backward
func determineTransitionDirection(currentCategory, targetCategory, targetStatusName string) bool {
	// Check for impediment/backward status names first
	impedimentStatuses := map[string]bool{
		"on hold":  true,
		"blocked":  true,
		"hold":     true,
		"paused":   true,
		"waiting":  true,
		"deferred": true,
	}

	// Normalize status name to lowercase for comparison
	normalizedName := strings.ToLower(targetStatusName)
	if impedimentStatuses[normalizedName] {
		return false // Impediment states show left arrow
	}

	// Jira uses these standard status category keys
	categoryOrder := map[string]int{
		"new":           1, // To Do / Backlog
		"indeterminate": 2, // In Progress / In Development
		"done":          3, // Done / Closed / Completed
	}

	currentOrder, currentExists := categoryOrder[currentCategory]
	targetOrder, targetExists := categoryOrder[targetCategory]

	// Log categories for debugging
	log.Printf("Transition direction: current='%s' (order=%d, exists=%v) -> target='%s' (order=%d, exists=%v) status='%s'",
		currentCategory, currentOrder, currentExists, targetCategory, targetOrder, targetExists, targetStatusName)

	// If both categories exist in the map, compare their order
	if currentExists && targetExists {
		// If same category, treat as forward (lateral move)
		if targetOrder == currentOrder {
			return true
		}
		return targetOrder > currentOrder
	}

	// If categories are not in the map, try to infer direction
	// Treat unknown categories as "indeterminate" (middle state)
	if !currentExists {
		currentOrder = 2
	}
	if !targetExists {
		targetOrder = 2
	}

	return targetOrder > currentOrder
}

// ExecuteStatusTransition executes a status transition for a Jira issue
func ExecuteStatusTransition(issueKey string, transitionID string) error {
	// Format request body as {"transition": {"id": "transitionID"}}
	transitionData := map[string]interface{}{
		"transition": map[string]interface{}{
			"id": transitionID,
		},
	}

	// Call TransitionIssue API with transition data
	response, err := jiraApiFunctions.TransitionIssue(issueKey, transitionData)
	if err != nil {
		log.Printf("Error executing transition %s for issue %s: %v", transitionID, issueKey, err)
		return fmt.Errorf("failed to execute status transition: %w", err)
	}

	// Check if response indicates an error (non-empty response body might contain error details)
	if len(response) > 0 {
		var errorResponse struct {
			ErrorMessages []string          `json:"errorMessages"`
			Errors        map[string]string `json:"errors"`
		}
		
		if err := json.Unmarshal(response, &errorResponse); err == nil {
			if len(errorResponse.ErrorMessages) > 0 {
				log.Printf("Transition failed for %s: %v", issueKey, errorResponse.ErrorMessages)
				return fmt.Errorf("transition failed: %s", errorResponse.ErrorMessages[0])
			}
			if len(errorResponse.Errors) > 0 {
				log.Printf("Transition failed for %s: %v", issueKey, errorResponse.Errors)
				for field, msg := range errorResponse.Errors {
					return fmt.Errorf("transition failed: %s - %s", field, msg)
				}
			}
		}
	}

	log.Printf("Successfully executed transition %s for issue %s", transitionID, issueKey)
	return nil
}