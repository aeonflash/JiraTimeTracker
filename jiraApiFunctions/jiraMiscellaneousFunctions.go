package jiraApiFunctions

import (
	"fmt"
	"strings"
)

// Auditing APIs
func GetAuditRecords(offset, limit int, filter, from, to string) ([]byte, error) {
	params := map[string]string{
		"filter": filter,
		"from":   from,
		"to":     to,
	}
	if offset > 0 {
		params["offset"] = fmt.Sprintf("%d", offset)
	}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	return MakeJiraAPICall("GET", "/rest/api/3/auditing/record", nil, params)
}

// Avatar APIs
func GetSystemAvatars(avatarType string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/avatar/%s/system", avatarType)
	return MakeJiraAPICall("GET", endpoint, nil, nil)
}

// Changelog APIs
func GetChangelogsBulk(changelogIds interface{}) ([]byte, error) {
	return MakeJiraAPICall("POST", "/rest/api/3/changelog/bulkfetch", changelogIds, nil)
}

// Classification APIs
func GetClassificationLevels(status []string, orderBy string) ([]byte, error) {
	params := map[string]string{
		"orderBy": orderBy,
	}
	if len(status) > 0 {
		params["status"] = strings.Join(status, ",")
	}
	return MakeJiraAPICall("GET", "/rest/api/3/classification-levels", nil, params)
}

// Comment APIs
func GetCommentsList(commentRequest interface{}) ([]byte, error) {
	return MakeJiraAPICall("POST", "/rest/api/3/comment/list", commentRequest, nil)
}

func GetCommentProperties(commentId string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/comment/%s/properties", commentId)
	return MakeJiraAPICall("GET", endpoint, nil, nil)
}

func DeleteCommentProperty(commentId, propertyKey string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/comment/%s/properties/%s", commentId, propertyKey)
	return MakeJiraAPICall("DELETE", endpoint, nil, nil)
}

// Component APIs
func GetComponents(query, projectIdOrKey, orderBy string, maxResults int) ([]byte, error) {
	params := map[string]string{
		"query":          query,
		"projectIdOrKey": projectIdOrKey,
		"orderBy":        orderBy,
	}
	if maxResults > 0 {
		params["maxResults"] = fmt.Sprintf("%d", maxResults)
	}
	return MakeJiraAPICall("GET", "/rest/api/3/component", nil, params)
}

func DeleteComponent(id, moveIssuesTo string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/component/%s", id)
	params := map[string]string{
		"moveIssuesTo": moveIssuesTo,
	}
	return MakeJiraAPICall("DELETE", endpoint, nil, params)
}

func GetComponentRelatedIssueCounts(id string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/component/%s/relatedIssueCounts", id)
	return MakeJiraAPICall("GET", endpoint, nil, nil)
}

// Configuration APIs
func GetConfiguration() ([]byte, error) {
	return MakeJiraAPICall("GET", "/rest/api/3/configuration", nil, nil)
}

func GetTimeTrackingConfiguration() ([]byte, error) {
	return MakeJiraAPICall("GET", "/rest/api/3/configuration/timetracking", nil, nil)
}

func GetTimeTrackingProviders() ([]byte, error) {
	return MakeJiraAPICall("GET", "/rest/api/3/configuration/timetracking/list", nil, nil)
}

func GetTimeTrackingOptions() ([]byte, error) {
	return MakeJiraAPICall("GET", "/rest/api/3/configuration/timetracking/options", nil, nil)
}

// Custom Field APIs
func GetCustomFieldOption(id string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/customFieldOption/%s", id)
	return MakeJiraAPICall("GET", endpoint, nil, nil)
}

// Dashboard APIs
func GetDashboards(filter string, startAt, maxResults int) ([]byte, error) {
	params := map[string]string{
		"filter": filter,
	}
	if startAt > 0 {
		params["startAt"] = fmt.Sprintf("%d", startAt)
	}
	if maxResults > 0 {
		params["maxResults"] = fmt.Sprintf("%d", maxResults)
	}
	return MakeJiraAPICall("GET", "/rest/api/3/dashboard", nil, params)
}

func BulkEditDashboards(editRequest interface{}) ([]byte, error) {
	return MakeJiraAPICall("PUT", "/rest/api/3/dashboard/bulk/edit", editRequest, nil)
}

func GetAvailableGadgets(moduleKey, uri []string, gadgetId []int) ([]byte, error) {
	params := map[string]string{}
	if len(moduleKey) > 0 {
		params["moduleKey"] = strings.Join(moduleKey, ",")
	}
	if len(uri) > 0 {
		params["uri"] = strings.Join(uri, ",")
	}
	if len(gadgetId) > 0 {
		gadgetIds := make([]string, len(gadgetId))
		for i, id := range gadgetId {
			gadgetIds[i] = fmt.Sprintf("%d", id)
		}
		params["gadgetId"] = strings.Join(gadgetIds, ",")
	}
	return MakeJiraAPICall("GET", "/rest/api/3/dashboard/gadgets", nil, params)
}

func SearchDashboards(dashboardName, accountId, owner, groupname, groupId string, projectId int, orderBy, status, expand string, startAt, maxResults int) ([]byte, error) {
	params := map[string]string{
		"dashboardName": dashboardName,
		"accountId":     accountId,
		"owner":         owner,
		"groupname":     groupname,
		"groupId":       groupId,
		"orderBy":       orderBy,
		"status":        status,
		"expand":        expand,
	}
	if projectId > 0 {
		params["projectId"] = fmt.Sprintf("%d", projectId)
	}
	if startAt > 0 {
		params["startAt"] = fmt.Sprintf("%d", startAt)
	}
	if maxResults > 0 {
		params["maxResults"] = fmt.Sprintf("%d", maxResults)
	}
	return MakeJiraAPICall("GET", "/rest/api/3/dashboard/search", nil, params)
}

// Data Policy APIs
func GetDataPolicy() ([]byte, error) {
	return MakeJiraAPICall("GET", "/rest/api/3/data-policy", nil, nil)
}

func GetProjectDataPolicy(ids string) ([]byte, error) {
	params := map[string]string{
		"ids": ids,
	}
	return MakeJiraAPICall("GET", "/rest/api/3/data-policy/project", nil, params)
}