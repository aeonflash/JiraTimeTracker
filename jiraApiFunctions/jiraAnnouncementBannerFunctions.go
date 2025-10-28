package jiraApiFunctions

// Announcement Banner APIs
func GetAnnouncementBanner() ([]byte, error) {
	return MakeJiraAPICall("GET", "/rest/api/3/announcementBanner", nil, nil)
}

func SetAnnouncementBanner(config interface{}) ([]byte, error) {
	return MakeJiraAPICall("PUT", "/rest/api/3/announcementBanner", config, nil)
}