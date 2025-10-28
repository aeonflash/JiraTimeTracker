package jiraApiFunctions

import (
	"fmt"
	"strings"
)

// Field APIs
func GetFields() ([]byte, error) {
	return MakeJiraAPICall("GET", "/rest/api/3/field", nil, nil)
}

func CreateCustomField(fieldData interface{}) ([]byte, error) {
	return MakeJiraAPICall("POST", "/rest/api/3/field", fieldData, nil)
}

func SearchFields(expand string, startAt, maxResults int, types, ids []string, query, orderBy string) ([]byte, error) {
	params := map[string]string{
		"expand":  expand,
		"query":   query,
		"orderBy": orderBy,
	}
	if len(types) > 0 {
		params["type"] = strings.Join(types, ",")
	}
	if len(ids) > 0 {
		params["id"] = strings.Join(ids, ",")
	}
	if startAt > 0 {
		params["startAt"] = fmt.Sprintf("%d", startAt)
	}
	if maxResults > 0 {
		params["maxResults"] = fmt.Sprintf("%d", maxResults)
	}
	return MakeJiraAPICall("GET", "/rest/api/3/field/search", nil, params)
}

func GetField(fieldId string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/field/%s", fieldId)
	return MakeJiraAPICall("GET", endpoint, nil, nil)
}

func UpdateField(fieldId string, fieldData interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/field/%s", fieldId)
	return MakeJiraAPICall("PUT", endpoint, fieldData, nil)
}

func DeleteField(fieldId string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/field/%s", fieldId)
	return MakeJiraAPICall("DELETE", endpoint, nil, nil)
}