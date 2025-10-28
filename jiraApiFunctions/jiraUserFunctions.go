package jiraApiFunctions

import "fmt"

// User APIs
func GetCurrentUser(expand string) ([]byte, error) {
	params := map[string]string{
		"expand": expand,
	}
	return MakeJiraAPICall("GET", "/rest/api/3/myself", nil, params)
}

func GetUser(accountId, username, key, expand string) ([]byte, error) {
	params := map[string]string{
		"accountId": accountId,
		"username":  username,
		"key":       key,
		"expand":    expand,
	}
	return MakeJiraAPICall("GET", "/rest/api/3/user", nil, params)
}

func CreateUser(userData interface{}) ([]byte, error) {
	return MakeJiraAPICall("POST", "/rest/api/3/user", userData, nil)
}

func DeleteUser(accountId, username, key string) ([]byte, error) {
	params := map[string]string{
		"accountId": accountId,
		"username":  username,
		"key":       key,
	}
	return MakeJiraAPICall("DELETE", "/rest/api/3/user", nil, params)
}

func FindUsers(query string, startAt, maxResults int, property string) ([]byte, error) {
	params := map[string]string{
		"query":    query,
		"property": property,
	}
	if startAt > 0 {
		params["startAt"] = fmt.Sprintf("%d", startAt)
	}
	if maxResults > 0 {
		params["maxResults"] = fmt.Sprintf("%d", maxResults)
	}
	return MakeJiraAPICall("GET", "/rest/api/3/user/search", nil, params)
}

func GetUserGroups(accountId, username, key string) ([]byte, error) {
	params := map[string]string{
		"accountId": accountId,
		"username":  username,
		"key":       key,
	}
	return MakeJiraAPICall("GET", "/rest/api/3/user/groups", nil, params)
}