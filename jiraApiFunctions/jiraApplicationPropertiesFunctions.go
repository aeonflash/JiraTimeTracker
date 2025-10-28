package jiraApiFunctions

import "fmt"

// Application Properties APIs
func GetApplicationProperties(key, keyFilter, permissionLevel string) ([]byte, error) {
	params := map[string]string{
		"key":             key,
		"keyFilter":       keyFilter,
		"permissionLevel": permissionLevel,
	}
	return MakeJiraAPICall("GET", "/rest/api/3/application-properties", nil, params)
}

func GetAdvancedSettings() ([]byte, error) {
	return MakeJiraAPICall("GET", "/rest/api/3/application-properties/advanced-settings", nil, nil)
}

func SetApplicationProperty(id string, property interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/application-properties/%s", id)
	return MakeJiraAPICall("PUT", endpoint, property, nil)
}