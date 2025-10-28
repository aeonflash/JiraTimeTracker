package jiraApiFunctions

import (
	"fmt"
	"strings"
)

// Project APIs
func GetProjects(expand string, recent int, properties []string) ([]byte, error) {
	params := map[string]string{
		"expand": expand,
	}
	if recent > 0 {
		params["recent"] = fmt.Sprintf("%d", recent)
	}
	if len(properties) > 0 {
		params["properties"] = strings.Join(properties, ",")
	}
	return MakeJiraAPICall("GET", "/rest/api/3/project", nil, params)
}

func CreateProject(projectData interface{}) ([]byte, error) {
	return MakeJiraAPICall("POST", "/rest/api/3/project", projectData, nil)
}

func GetProject(projectIdOrKey, expand string, properties []string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/project/%s", projectIdOrKey)
	params := map[string]string{
		"expand": expand,
	}
	if len(properties) > 0 {
		params["properties"] = strings.Join(properties, ",")
	}
	return MakeJiraAPICall("GET", endpoint, nil, params)
}

func UpdateProject(projectIdOrKey string, projectData interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/project/%s", projectIdOrKey)
	return MakeJiraAPICall("PUT", endpoint, projectData, nil)
}

func DeleteProject(projectIdOrKey string, enableUndo bool) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/project/%s", projectIdOrKey)
	params := map[string]string{}
	if enableUndo {
		params["enableUndo"] = "true"
	}
	return MakeJiraAPICall("DELETE", endpoint, nil, params)
}

func GetProjectComponents(projectIdOrKey string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/project/%s/components", projectIdOrKey)
	return MakeJiraAPICall("GET", endpoint, nil, nil)
}

func GetProjectVersions(projectIdOrKey, expand string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/project/%s/versions", projectIdOrKey)
	params := map[string]string{
		"expand": expand,
	}
	return MakeJiraAPICall("GET", endpoint, nil, params)
}