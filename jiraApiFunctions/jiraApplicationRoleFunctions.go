package jiraApiFunctions

import "fmt"

// Application Role APIs
func GetApplicationRoles() ([]byte, error) {
	return MakeJiraAPICall("GET", "/rest/api/3/applicationrole", nil, nil)
}

func GetApplicationRole(key string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/applicationrole/%s", key)
	return MakeJiraAPICall("GET", endpoint, nil, nil)
}