# Jira REST API Available Functions

## Announcement Banner Functions

### GetAnnouncementBanner
**Description:** Get current announcement banner configuration  
**Required Params:** None  
**Expected Return:** `[]byte` - AnnouncementBannerConfiguration JSON  
**Example:**
```go
banner, err := GetAnnouncementBanner()
```

### SetAnnouncementBanner
**Description:** Update announcement banner configuration  
**Required Params:** `config AnnouncementBannerConfigurationUpdate`  
**Expected Return:** `[]byte` - Empty response (204)  
**Example:**
```go
config := AnnouncementBannerConfigurationUpdate{
    IsEnabled: true,
    Message: "System maintenance tonight",
}
result, err := SetAnnouncementBanner(config)
```

## Application Properties Functions

### GetApplicationProperties
**Description:** Get application properties with optional filtering  
**Required Params:** `key, keyFilter, permissionLevel string` (all optional)  
**Expected Return:** `[]byte` - Array of ApplicationProperty JSON  
**Example:**
```go
props, err := GetApplicationProperties("", "", "")
```

### GetAdvancedSettings
**Description:** Get advanced application settings  
**Required Params:** None  
**Expected Return:** `[]byte` - Array of ApplicationProperty JSON  
**Example:**
```go
settings, err := GetAdvancedSettings()
```

### SetApplicationProperty
**Description:** Update a specific application property  
**Required Params:** `id string, property ApplicationProperty`  
**Expected Return:** `[]byte` - ApplicationProperty JSON  
**Example:**
```go
prop := ApplicationProperty{Value: "new-value"}
result, err := SetApplicationProperty("jira.title", prop)
```

## Application Role Functions

### GetApplicationRoles
**Description:** Get all application roles  
**Required Params:** None  
**Expected Return:** `[]byte` - Array of ApplicationRole JSON  
**Example:**
```go
roles, err := GetApplicationRoles()
```

### GetApplicationRole
**Description:** Get specific application role by key  
**Required Params:** `key string`  
**Expected Return:** `[]byte` - ApplicationRole JSON  
**Example:**
```go
role, err := GetApplicationRole("jira-users")
```

## Attachment Functions

### GetAttachmentContent
**Description:** Download attachment file content  
**Required Params:** `id string, redirect bool`  
**Expected Return:** `[]byte` - Binary file content  
**Example:**
```go
content, err := GetAttachmentContent("12345", false)
```

### GetAttachmentMeta
**Description:** Get attachment system settings  
**Required Params:** None  
**Expected Return:** `[]byte` - AttachmentSettings JSON  
**Example:**
```go
meta, err := GetAttachmentMeta()
```

### GetAttachmentThumbnail
**Description:** Get attachment thumbnail image  
**Required Params:** `id string, redirect, fallbackToDefault bool, width, height int`  
**Expected Return:** `[]byte` - Binary image data  
**Example:**
```go
thumb, err := GetAttachmentThumbnail("12345", false, true, 150, 150)
```

### DeleteAttachment
**Description:** Delete an attachment from an issue  
**Required Params:** `id string`  
**Expected Return:** `[]byte` - Empty response  
**Example:**
```go
result, err := DeleteAttachment("12345")
```

## Issue Functions

### GetIssue
**Description:** Get issue details by ID or key  
**Required Params:** `issueIdOrKey string, fields, expand string` (fields/expand optional)  
**Expected Return:** `[]byte` - Issue JSON object  
**Example:**
```go
issue, err := GetIssue("PROJ-123", "summary,status", "transitions")
```

### CreateIssue
**Description:** Create a new issue  
**Required Params:** `issueData interface{}`  
**Expected Return:** `[]byte` - Created issue JSON  
**Example:**
```go
issueData := map[string]interface{}{
    "fields": map[string]interface{}{
        "project": map[string]string{"key": "PROJ"},
        "summary": "New issue",
        "issuetype": map[string]string{"name": "Task"},
    },
}
issue, err := CreateIssue(issueData)
```

### UpdateIssue
**Description:** Update an existing issue  
**Required Params:** `issueIdOrKey string, issueData interface{}`  
**Expected Return:** `[]byte` - Empty response (204)  
**Example:**
```go
updateData := map[string]interface{}{
    "fields": map[string]interface{}{
        "summary": "Updated summary",
    },
}
result, err := UpdateIssue("PROJ-123", updateData)
```

### DeleteIssue
**Description:** Delete an issue  
**Required Params:** `issueIdOrKey, deleteSubtasks string`  
**Expected Return:** `[]byte` - Empty response  
**Example:**
```go
result, err := DeleteIssue("PROJ-123", "true")
```

### GetIssueTransitions
**Description:** Get available transitions for an issue  
**Required Params:** `issueIdOrKey, expand string`  
**Expected Return:** `[]byte` - Transitions JSON array  
**Example:**
```go
transitions, err := GetIssueTransitions("PROJ-123", "")
```

### TransitionIssue
**Description:** Transition an issue to a new status  
**Required Params:** `issueIdOrKey string, transitionData interface{}`  
**Expected Return:** `[]byte` - Empty response (204)  
**Example:**
```go
transitionData := map[string]interface{}{
    "transition": map[string]string{"id": "31"},
}
result, err := TransitionIssue("PROJ-123", transitionData)
```

### GetIssueComments
**Description:** Get comments for an issue  
**Required Params:** `issueIdOrKey string, startAt, maxResults int, orderBy, expand string`  
**Expected Return:** `[]byte` - Comments JSON array  
**Example:**
```go
comments, err := GetIssueComments("PROJ-123", 0, 50, "created", "")
```

### AddComment
**Description:** Add a comment to an issue  
**Required Params:** `issueIdOrKey string, commentData interface{}`  
**Expected Return:** `[]byte` - Created comment JSON  
**Example:**
```go
commentData := map[string]interface{}{
    "body": "This is a comment",
}
comment, err := AddComment("PROJ-123", commentData)
```

### UpdateComment
**Description:** Update an existing comment  
**Required Params:** `issueIdOrKey, commentId string, commentData interface{}`  
**Expected Return:** `[]byte` - Updated comment JSON  
**Example:**
```go
commentData := map[string]interface{}{
    "body": "Updated comment text",
}
comment, err := UpdateComment("PROJ-123", "67890", commentData)
```

### DeleteComment
**Description:** Delete a comment from an issue  
**Required Params:** `issueIdOrKey, commentId string`  
**Expected Return:** `[]byte` - Empty response  
**Example:**
```go
result, err := DeleteComment("PROJ-123", "67890")
```

### AddWorklog
**Description:** Log work time on an issue  
**Required Params:** `issueIdOrKey string, worklogData interface{}`  
**Expected Return:** `[]byte` - Created worklog JSON  
**Example:**
```go
worklogData := map[string]interface{}{
    "timeSpent": "2h",
    "comment": "Fixed the bug",
}
worklog, err := AddWorklog("PROJ-123", worklogData)
```

## Project Functions

### GetProjects
**Description:** Get list of all projects  
**Required Params:** `expand string, recent int, properties []string`  
**Expected Return:** `[]byte` - Projects JSON array  
**Example:**
```go
projects, err := GetProjects("", 0, []string{})
```

### CreateProject
**Description:** Create a new project  
**Required Params:** `projectData interface{}`  
**Expected Return:** `[]byte` - Created project JSON  
**Example:**
```go
projectData := map[string]interface{}{
    "key": "TEST",
    "name": "Test Project",
    "projectTypeKey": "software",
}
project, err := CreateProject(projectData)
```

### GetProject
**Description:** Get specific project details  
**Required Params:** `projectIdOrKey, expand string, properties []string`  
**Expected Return:** `[]byte` - Project JSON object  
**Example:**
```go
project, err := GetProject("PROJ", "", []string{})
```

### UpdateProject
**Description:** Update project details  
**Required Params:** `projectIdOrKey string, projectData interface{}`  
**Expected Return:** `[]byte` - Updated project JSON  
**Example:**
```go
updateData := map[string]interface{}{
    "name": "Updated Project Name",
}
project, err := UpdateProject("PROJ", updateData)
```

### DeleteProject
**Description:** Delete a project  
**Required Params:** `projectIdOrKey string, enableUndo bool`  
**Expected Return:** `[]byte` - Empty response  
**Example:**
```go
result, err := DeleteProject("PROJ", true)
```

## User Functions

### GetCurrentUser
**Description:** Get current authenticated user details  
**Required Params:** `expand string`  
**Expected Return:** `[]byte` - User JSON object  
**Example:**
```go
user, err := GetCurrentUser("groups")
```

### GetUser
**Description:** Get specific user details  
**Required Params:** `accountId, username, key, expand string`  
**Expected Return:** `[]byte` - User JSON object  
**Example:**
```go
user, err := GetUser("5b10a2844c20165700ede21g", "", "", "")
```

### CreateUser
**Description:** Create a new user account  
**Required Params:** `userData interface{}`  
**Expected Return:** `[]byte` - Created user JSON  
**Example:**
```go
userData := map[string]interface{}{
    "name": "testuser",
    "emailAddress": "test@example.com",
    "displayName": "Test User",
}
user, err := CreateUser(userData)
```

### FindUsers
**Description:** Search for users by query  
**Required Params:** `query string, startAt, maxResults int, property string`  
**Expected Return:** `[]byte` - Users JSON array  
**Example:**
```go
users, err := FindUsers("john", 0, 50, "")
```

## Search Functions

### SearchIssues
**Description:** Search issues using JQL query  
**Required Params:** `jql, expand string, fields []string, startAt, maxResults int, validateQuery bool`  
**Expected Return:** `[]byte` - Search results JSON  
**Example:**
```go
results, err := SearchIssues("project = PROJ AND status = Open", "", []string{"summary", "status"}, 0, 50, true)
```

### SearchIssuesPost
**Description:** Search issues using POST with complex query  
**Required Params:** `searchRequest interface{}`  
**Expected Return:** `[]byte` - Search results JSON  
**Example:**
```go
searchRequest := map[string]interface{}{
    "jql": "project = PROJ",
    "fields": []string{"summary", "status"},
    "maxResults": 50,
}
results, err := SearchIssuesPost(searchRequest)
```

## Field Functions

### GetFields
**Description:** Get all system and custom fields  
**Required Params:** None  
**Expected Return:** `[]byte` - Fields JSON array  
**Example:**
```go
fields, err := GetFields()
```

### CreateCustomField
**Description:** Create a new custom field  
**Required Params:** `fieldData interface{}`  
**Expected Return:** `[]byte` - Created field JSON  
**Example:**
```go
fieldData := map[string]interface{}{
    "name": "My Custom Field",
    "type": "com.atlassian.jira.plugin.system.customfieldtypes:textfield",
}
field, err := CreateCustomField(fieldData)
```

### SearchFields
**Description:** Search fields with filters  
**Required Params:** `expand string, startAt, maxResults int, types, ids []string, query, orderBy string`  
**Expected Return:** `[]byte` - Fields JSON array  
**Example:**
```go
fields, err := SearchFields("", 0, 50, []string{}, []string{}, "custom", "name")
```

## Group Functions

### GetGroup
**Description:** Get group details  
**Required Params:** `groupname, groupId, expand string`  
**Expected Return:** `[]byte` - Group JSON object  
**Example:**
```go
group, err := GetGroup("jira-users", "", "users")
```

### CreateGroup
**Description:** Create a new group  
**Required Params:** `groupData interface{}`  
**Expected Return:** `[]byte` - Created group JSON  
**Example:**
```go
groupData := map[string]interface{}{
    "name": "new-group",
}
group, err := CreateGroup(groupData)
```

### FindGroups
**Description:** Search for groups  
**Required Params:** `query string, exclude []string, maxResults int, userName string`  
**Expected Return:** `[]byte` - Groups JSON array  
**Example:**
```go
groups, err := FindGroups("admin", []string{}, 50, "")
```

## Bulk Operations Functions

### BulkDeleteIssues
**Description:** Delete multiple issues at once  
**Required Params:** `issuesUpdate interface{}`  
**Expected Return:** `[]byte` - Bulk operation result JSON  
**Example:**
```go
deleteData := map[string]interface{}{
    "issueIds": []string{"10001", "10002"},
}
result, err := BulkDeleteIssues(deleteData)
```

### BulkWatchIssues
**Description:** Watch multiple issues at once  
**Required Params:** `watchRequest interface{}`  
**Expected Return:** `[]byte` - Empty response  
**Example:**
```go
watchData := map[string]interface{}{
    "issueIds": []string{"10001", "10002"},
}
result, err := BulkWatchIssues(watchData)
```