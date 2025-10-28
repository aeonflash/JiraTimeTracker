package jiraApiFunctions

import "fmt"

// Attachment APIs
func GetAttachmentContent(id string, redirect bool) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/attachment/content/%s", id)
	params := map[string]string{}
	if redirect {
		params["redirect"] = "true"
	}
	return MakeJiraAPICall("GET", endpoint, nil, params)
}

func GetAttachmentMeta() ([]byte, error) {
	return MakeJiraAPICall("GET", "/rest/api/3/attachment/meta", nil, nil)
}

func GetAttachmentThumbnail(id string, redirect, fallbackToDefault bool, width, height int) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/attachment/thumbnail/%s", id)
	params := map[string]string{}
	if redirect {
		params["redirect"] = "true"
	}
	if fallbackToDefault {
		params["fallbackToDefault"] = "true"
	}
	if width > 0 {
		params["width"] = fmt.Sprintf("%d", width)
	}
	if height > 0 {
		params["height"] = fmt.Sprintf("%d", height)
	}
	return MakeJiraAPICall("GET", endpoint, nil, params)
}

func DeleteAttachment(id string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/attachment/%s", id)
	return MakeJiraAPICall("DELETE", endpoint, nil, nil)
}

func GetAttachmentExpandHuman(id string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/attachment/%s/expand/human", id)
	return MakeJiraAPICall("GET", endpoint, nil, nil)
}

func GetAttachmentExpandRaw(id string) ([]byte, error) {
	endpoint := fmt.Sprintf("/rest/api/3/attachment/%s/expand/raw", id)
	return MakeJiraAPICall("GET", endpoint, nil, nil)
}