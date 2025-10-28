# Jira GraphQL API Functions Summary

## Atlassian Studio Functions

```
function: atlassianStudio_userSiteContext
variables: cloudId - ID!
return object: AtlassianStudioUserSiteContextResult
description: Queries the products available on a site and user permissions to render Atlassian Studio experience for a given site
scopes: Not specified
```

## Customer Support Functions

```
function: customerSupport
variables: None
return object: SupportRequestCatalogQueryApi
description: This API is a wrapper for all CSP support Request queries
scopes: Not specified
```

## JSM Chat Functions

```
function: jsmChat
variables: None
return object: JsmChatQuery
description: Enable self governed onboarding of Jira GraphQL types to AGG (not available for OAuth authenticated requests)
scopes: Not available for OAuth requests
```

## Playbook Functions

```
function: playbook_jiraPlaybook
variables: playbookAri - ID!
return object: JiraPlaybookQueryPayload
description: Fetch playbook by playbook ARI
scopes: Requires '@optIn(to: "PlaybooksInJSM")' directive (EXPERIMENTAL)
```

```
function: playbook_jiraPlaybookInstancesForIssue
variables: after - String, cloudId - ID!, first - Int (default: 20), issueId - String!, projectKey - String!
return object: JiraPlaybookInstanceConnection
description: Used when user clicks "Playbooks" expandable section in issue view for Show output/Refresh/View Output/Browser reload
scopes: Requires '@optIn(to: "PlaybooksInJSM")' directive (EXPERIMENTAL)
```

```
function: playbook_jiraPlaybookLabelsForProject
variables: after - String, cloudId - ID!, filters - JiraPlaybookLabelFilter, first - Int (default: 10), projectKey - String!
return object: JiraPlaybookLabelConnection
description: Fetch all Playbook Labels for a project
scopes: Requires '@optIn(to: "PlaybooksInJSM")' directive (EXPERIMENTAL)
```

```
function: playbook_jiraPlaybookStepRunsForPlaybookInstance
variables: after - String, first - Int (default: 20), playbookInstanceAri - ID!
return object: JiraPlaybookStepRunConnection
description: Used when user clicks "Execution Output" tab in playbook
scopes: Requires '@optIn(to: "PlaybooksInJSM")' directive (EXPERIMENTAL)
```

```
function: playbook_jiraPlaybookStepRunsForProject
variables: after - String, cloudId - ID!, filters - JiraPlaybookExecutionFilter, first - Int (default: 20), projectKey - String!
return object: JiraPlaybookStepRunConnection
description: Used in "Execution Log" tab in Admin View
scopes: Requires '@optIn(to: "PlaybooksInJSM")' directive (EXPERIMENTAL)
```

```
function: playbook_jiraPlaybookStepUsageForProject
variables: after - String, cloudId - ID!, filters - JiraPlaybookStepUsageFilter, first - Int (default: 20), projectKey - String!
return object: JiraPlaybookStepUsageConnection
description: Used in Usage Tab in Admin view
scopes: Requires '@optIn(to: "PlaybooksInJSM")' directive (EXPERIMENTAL)
```

```
function: playbook_jiraPlaybooksForProject
variables: after - String, cloudId - ID!, filters - JiraPlaybookListFilter, first - Int (default: 20), projectKey - String!, sort - [JiraPlaybooksSortInput!] (default: [{by: NAME, order: ASC}])
return object: JiraPlaybookConnection
description: Used in List Playbook in Admin View
scopes: Requires '@optIn(to: "PlaybooksInJSM")' directive (EXPERIMENTAL)
```

## Glance Functions

```
function: glance_getVULNIssues
variables: None
return object: [GlanceUserInsights]
description: Get vulnerability issues insights
scopes: Not specified
```

```
function: glance_getPipelineEvents
variables: None
return object: [GlanceUserInsights]
description: Get pipeline events insights
scopes: Not specified
```

```
function: glance_getCurrentUserSettings
variables: None
return object: UserSettings
description: Get current user settings
scopes: Not specified
```

## Catchup Functions

```
function: catchupEditMetadataForContent
variables: cloudId - String, contentId - ID!, contentType - CatchupContentType!, endTimeMs - Long!, updateType - CatchupOverviewUpdateType
return object: CatchupEditMetadataForContent
description: Get edit metadata for content with time-based filtering
scopes: confluence:atlassian-external
```

```
function: catchupGetLastViewedTime
variables: cloudId - String, contentId - ID!, contentType - CatchupContentType!
return object: CatchupLastViewedTimeResponse
description: Get the last viewed time for specific content
scopes: confluence:atlassian-external
```

```
function: catchupVersionDiffMetadataForContent
variables: cloudId - String, contentId - ID!, contentType - CatchupContentType!, originalContentVersion - Int!, revisedContentVersion - Int!
return object: CatchupVersionDiffMetadataResponse
description: Get version diff metadata between two content versions
scopes: confluence:atlassian-external
```

## Confluence Functions

```
function: confluence_contentAISummaries
variables: contentAris - [ID]!, objectType - KnowledgeGraphObjectType!
return object: [ConfluenceContentAISummaryResponse]
description: Get AI summaries for Confluence content
scopes: confluence:atlassian-external
```

```
function: confluence_latestKnowledgeGraphObjectV2
variables: cloudId - String!, contentId - ID!, contentType - KnowledgeGraphContentType!, objectType - KnowledgeGraphObjectType!
return object: KnowledgeGraphObjectResponseV2
description: Get latest knowledge graph object for content
scopes: confluence:atlassian-external
```

## Feed Functions

```
function: feed
variables: after - String, cloudId - String, first - Int (default: 25), sortBy - String
return object: PaginatedFeed
description: Get paginated feed content
scopes: confluence:atlassian-external
```

```
function: forYouFeed
variables: after - String, cloudId - String, first - Int (default: 5)
return object: ForYouPaginatedFeed
description: Get personalized "For You" feed content
scopes: confluence:atlassian-external
```

## AI and Smart Features Functions

```
function: getAIConfig
variables: cloudId - String, product - Product!
return object: AIConfigResponse
description: Get AI configuration for a specific product
scopes: confluence:atlassian-external
```

```
function: getCommentReplySuggestions
variables: cloudId - String, commentId - ID!, language - String
return object: CommentReplySuggestions
description: Get AI-powered reply suggestions for comments
scopes: confluence:atlassian-external
```

```
function: getCommentsSummary
variables: cloudId - String, commentsType - CommentsType!, contentId - ID!, contentType - SummaryType!, language - String
return object: SmartFeaturesCommentsSummary
description: Get AI summary of comments for content
scopes: confluence:atlassian-external
```

## Feed Configuration Functions

```
function: getFeedUserConfig
variables: cloudId - String
return object: FollowingFeedGetUserConfig
description: Get user configuration for following feed
scopes: confluence:atlassian-external
```

```
function: getRecommendedFeedUserConfig
variables: cloudId - String
return object: RecommendedFeedUserConfig
description: Get user configuration for recommended feed
scopes: confluence:atlassian-external
```

## Recommendation Functions

```
function: getRecommendedLabels
variables: cloudId - String, entityId - ID!, entityType - String!, first - Int, spaceId - ID!
return object: RecommendedLabels
description: Get recommended labels for an entity in a space
scopes: confluence:atlassian-external
```

```
function: getRecommendedPages
variables: cloudId - String, entityId - ID!, entityType - String!, experience - String!
return object: RecommendedPages
description: Get recommended pages based on entity and experience context
scopes: confluence:atlassian-external
```

```
function: getRecommendedPagesSpaceStatus
variables: cloudId - String (incomplete in source)
return object: Not specified (incomplete in source)
description: Get recommended pages space status (incomplete function definition)
scopes: confluence:atlassian-external
```