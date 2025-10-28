package jiraApiFunctions

import (
	"fmt"
)

// Issue APIs
func GetIssue(issueIdOrKey string, fields, expand string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s", issueIdOrKey)
	params := map[string]string{
		"fields": fields,
		"expand": expand,
	}
	return MakeJiraAPICall("GET", endpoint, nil, params)
}

func CreateIssue(issueData interface{}) ([]byte, error) {
	return MakeJiraAPICall("POST", "/rest/api/3/issue", issueData, nil)
}

func UpdateIssue(issueIdOrKey string, issueData interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s", issueIdOrKey)
	return MakeJiraAPICall("PUT", endpoint, issueData, nil)
}

func DeleteIssue(issueIdOrKey, deleteSubtasks string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s", issueIdOrKey)
	params := map[string]string{
		"deleteSubtasks": deleteSubtasks,
	}
	return MakeJiraAPICall("DELETE", endpoint, nil, params)
}

func GetIssueTransitions(issueIdOrKey, expand string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueIdOrKey)
	params := map[string]string{
		"expand": expand,
	}
	return MakeJiraAPICall("GET", endpoint, nil, params)
}

func TransitionIssue(issueIdOrKey string, transitionData interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueIdOrKey)
	return MakeJiraAPICall("POST", endpoint, transitionData, nil)
}

func GetIssueComments(issueIdOrKey string, startAt, maxResults int, orderBy, expand string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/comment", issueIdOrKey)
	params := map[string]string{
		"orderBy": orderBy,
		"expand":  expand,
	}
	if startAt > 0 {
		params["startAt"] = fmt.Sprintf("%d", startAt)
	}
	if maxResults > 0 {
		params["maxResults"] = fmt.Sprintf("%d", maxResults)
	}
	return MakeJiraAPICall("GET", endpoint, nil, params)
}

func AddComment(issueIdOrKey string, commentData interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/comment", issueIdOrKey)
	return MakeJiraAPICall("POST", endpoint, commentData, nil)
}

func UpdateComment(issueIdOrKey, commentId string, commentData interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/comment/%s", issueIdOrKey, commentId)
	return MakeJiraAPICall("PUT", endpoint, commentData, nil)
}

func DeleteComment(issueIdOrKey, commentId string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/comment/%s", issueIdOrKey, commentId)
	return MakeJiraAPICall("DELETE", endpoint, nil, nil)
}

func GetIssueWatchers(issueIdOrKey string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/watchers", issueIdOrKey)
	return MakeJiraAPICall("GET", endpoint, nil, nil)
}

func AddWatcher(issueIdOrKey, accountId string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/watchers", issueIdOrKey)
	return MakeJiraAPICall("POST", endpoint, accountId, nil)
}

func RemoveWatcher(issueIdOrKey, accountId string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/watchers", issueIdOrKey)
	params := map[string]string{
		"accountId": accountId,
	}
	return MakeJiraAPICall("DELETE", endpoint, nil, params)
}

func GetIssueWorklog(issueIdOrKey string, startAt, maxResults int, expand string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/worklog", issueIdOrKey)
	params := map[string]string{
		"expand": expand,
	}
	if startAt > 0 {
		params["startAt"] = fmt.Sprintf("%d", startAt)
	}
	if maxResults > 0 {
		params["maxResults"] = fmt.Sprintf("%d", maxResults)
	}
	return MakeJiraAPICall("GET", endpoint, nil, params)
}

func AddWorklog(issueIdOrKey string, worklogData interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/worklog", issueIdOrKey)
	return MakeJiraAPICall("POST", endpoint, worklogData, nil)
}

func UpdateWorklog(issueIdOrKey, worklogId string, worklogData interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s", issueIdOrKey, worklogId)
	return MakeJiraAPICall("PUT", endpoint, worklogData, nil)
}

func DeleteWorklog(issueIdOrKey, worklogId string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s", issueIdOrKey, worklogId)
	return MakeJiraAPICall("DELETE", endpoint, nil, nil)
}