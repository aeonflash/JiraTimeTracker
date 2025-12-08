package jiraApiFunctions

import (
	"fmt"
	"strings"
)

// Search APIs
func SearchIssues(jql, expand string, fields []string, startAt, maxResults int, validateQuery bool) ([]byte, error) {
	params := map[string]string{
		"jql": jql,
	}
	if expand != "" {
		params["expand"] = expand
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
	// Use the new search/jql endpoint with GET
	return MakeJiraAPICall("GET", "/rest/api/3/search/jql", nil, params)
}

func SearchIssuesPost(searchRequest interface{}) ([]byte, error) {
	// Use the standard search endpoint with POST
	return MakeJiraAPICall("POST", "/rest/api/2/search", searchRequest, nil)
}