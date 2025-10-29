package main

import (
	"log"
	"time"
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

func logWorkToJira(jiraId string, timeSpent string, comment string, startTime time.Time) error {
	worklogData := map[string]interface{}{
		"timeSpent": timeSpent,
		"comment": map[string]interface{}{
			"type": "doc",
			"version": 1,
			"content": []map[string]interface{}{
				{
					"type": "paragraph",
					"content": []map[string]interface{}{
						{
							"text": comment,
							"type": "text",
						},
					},
				},
			},
		},
		"started": startTime.Format("2006-01-02T15:04:05.000-0700"),
	}
	
	_, err := jiraApiFunctions.AddWorklog(jiraId, worklogData)
	return err
}