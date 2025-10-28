package jiraApiFunctions

import (
	"fmt"
	"strings"
)

// Group APIs
func GetGroup(groupname, groupId, expand string) ([]byte, error) {
	params := map[string]string{
		"groupname": groupname,
		"groupId":   groupId,
		"expand":    expand,
	}
	return MakeJiraAPICall("GET", "/rest/api/3/group", nil, params)
}

func CreateGroup(groupData interface{}) ([]byte, error) {
	return MakeJiraAPICall("POST", "/rest/api/3/group", groupData, nil)
}

func DeleteGroup(groupname, groupId, swapGroup, swapGroupId string) ([]byte, error) {
	params := map[string]string{
		"groupname":   groupname,
		"groupId":     groupId,
		"swapGroup":   swapGroup,
		"swapGroupId": swapGroupId,
	}
	return MakeJiraAPICall("DELETE", "/rest/api/3/group", nil, params)
}

func FindGroups(query string, exclude []string, maxResults int, userName string) ([]byte, error) {
	params := map[string]string{
		"query":    query,
		"userName": userName,
	}
	if len(exclude) > 0 {
		params["exclude"] = strings.Join(exclude, ",")
	}
	if maxResults > 0 {
		params["maxResults"] = fmt.Sprintf("%d", maxResults)
	}
	return MakeJiraAPICall("GET", "/rest/api/3/groups/picker", nil, params)
}