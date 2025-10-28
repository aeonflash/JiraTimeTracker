package jiraApiFunctions

import (
	"fmt"
	"strings"
)

// Search APIs
func SearchIssues(jql, expand string, fields []string, startAt, maxResults int, validateQuery bool) ([]byte, error) {
	params := map[string]string{
		"jql":    jql,
		"expand": expand,
	}
	if len(fields) > 0 {
		params["fields"] = strings.Join(fields, ",")
	}
	if startAt > 0 {
		params["startAt"] = fmt.Sprintf("%d", startAt)
	}
	if maxResults > 0 {
		params["maxResults"] = fmt.Sprintf("%d", maxResults)
	}
	if validateQuery {
		params["validateQuery"] = "true"
	}
	return MakeJiraAPICall("GET", "/rest/api/3/search", nil, params)
}

func SearchIssuesPost(searchRequest interface{}) ([]byte, error) {
	return MakeJiraAPICall("POST", "/rest/api/3/search", searchRequest, nil)
}