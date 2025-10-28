package main

import (
	"log"
	"jiraTimeWidget/jiraApiFunctions"
)

func getJiraItem(jiraId string, token string) string {
	jiraItem, err := jiraApiFunctions.GetIssue(jiraId, "", "names,renderedFields")
	if err != nil {
		log.Println("Error getting jira item:", err)
		return ""
	}
	return string(jiraItem)
}