package jiraApiFunctions

import (
	"fmt"
	"strings"
)

// Bulk Operations APIs
func BulkDeleteIssues(issuesUpdate interface{}) ([]byte, error) {
	return MakeJiraAPICall("POST", "/rest/api/3/bulk/issues/delete", issuesUpdate, nil)
}

func GetBulkIssueFields(issueIds []string, expand string) ([]byte, error) {
	params := map[string]string{
		"issueIds": strings.Join(issueIds, ","),
		"expand":   expand,
	}
	return MakeJiraAPICall("GET", "/rest/api/3/bulk/issues/fields", nil, params)
}

func BulkMoveIssues(moveRequest interface{}) ([]byte, error) {
	return MakeJiraAPICall("POST", "/rest/api/3/bulk/issues/move", moveRequest, nil)
}

func GetBulkIssueTransitions(issueIds []string, expand string) ([]byte, error) {
	params := map[string]string{
		"issueIds": strings.Join(issueIds, ","),
		"expand":   expand,
	}
	return MakeJiraAPICall("GET", "/rest/api/3/bulk/issues/transition", nil, params)
}

func BulkUnwatchIssues(unwatchRequest interface{}) ([]byte, error) {
	return MakeJiraAPICall("POST", "/rest/api/3/bulk/issues/unwatch", unwatchRequest, nil)
}

func BulkWatchIssues(watchRequest interface{}) ([]byte, error) {
	return MakeJiraAPICall("POST", "/rest/api/3/bulk/issues/watch", watchRequest, nil)
}

func GetBulkOperationStatus(taskId string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/bulk/queue/%s", taskId)
	return MakeJiraAPICall("GET", endpoint, nil, nil)
}