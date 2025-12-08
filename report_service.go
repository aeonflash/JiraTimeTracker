package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"jiraTimeWidget/jiraApiFunctions"
)

// ReportService handles data fetching and report generation
type ReportService struct {
	Cache       *ReportCache
	CurrentUser *User
}

// NewReportService creates a new ReportService instance
func NewReportService(user *User) *ReportService {
	return &ReportService{
		Cache:       NewReportCache(),
		CurrentUser: user,
	}
}

// JiraIssue represents a Jira issue from the API
type JiraIssue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary     string `json:"summary"`
		Status      struct {
			Name string `json:"name"`
		} `json:"status"`
		Assignee *struct {
			AccountId string `json:"accountId"`
		} `json:"assignee"`
		Created      string  `json:"created"`
		Updated      string  `json:"updated"`
		Resolutiondate *string `json:"resolutiondate"`
	} `json:"fields"`
}

// JiraSearchResponse represents the response from Jira search API
type JiraSearchResponse struct {
	Issues     []JiraIssue `json:"issues"`
	Total      int         `json:"total"`
	StartAt    int         `json:"startAt"`
	MaxResults int         `json:"maxResults"`
}

// JiraWorklog represents a worklog entry from the API
type JiraWorklog struct {
	ID              string `json:"id"`
	Author          struct {
		AccountId string `json:"accountId"`
	} `json:"author"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
	Started          string `json:"started"`
	Created          string `json:"created"`
}

// JiraWorklogResponse represents the response from worklog API
type JiraWorklogResponse struct {
	Worklogs   []JiraWorklog `json:"worklogs"`
	Total      int           `json:"total"`
	StartAt    int           `json:"startAt"`
	MaxResults int           `json:"maxResults"`
}

// GetAssignedIssues retrieves all issues assigned to the current user within the date range
func (rs *ReportService) GetAssignedIssues(dateRange DateRange) (*AssignedIssuesReport, error) {
	log.Printf("GetAssignedIssues called for date range: %s to %s", dateRange.StartDate.Format("2006-01-02"), dateRange.EndDate.Format("2006-01-02"))
	
	// Check cache first
	if rs.Cache.IsValid(dateRange) {
		if cached, exists := rs.Cache.Get("assigned_issues"); exists {
			if report, ok := cached.(*AssignedIssuesReport); ok {
				log.Printf("Returning cached assigned issues: %d items", len(report.Issues))
				return report, nil
			}
		}
	}

	// Build JQL query
	jql := fmt.Sprintf(
		"assignee = currentUser() AND updated >= \"%s\" AND updated <= \"%s\" ORDER BY key ASC",
		dateRange.StartDate.Format("2006-01-02"),
		dateRange.EndDate.Format("2006-01-02"),
	)
	log.Printf("JQL Query: %s", jql)

	// Fetch issues from Jira
	issues, err := rs.fetchAllIssues(jql)
	if err != nil {
		log.Printf("Error fetching assigned issues: %v", err)
		return nil, fmt.Errorf("failed to fetch assigned issues: %w", err)
	}
	log.Printf("Fetched %d assigned issues from Jira", len(issues))

	// Convert to report items
	reportItems := make([]IssueReportItem, 0, len(issues))
	for _, issue := range issues {
		item, err := rs.convertToReportItem(issue)
		if err != nil {
			log.Printf("Warning: failed to convert issue %s: %v", issue.Key, err)
			continue
		}
		reportItems = append(reportItems, item)
	}

	report := &AssignedIssuesReport{
		Issues:     reportItems,
		TotalCount: len(reportItems),
	}

	// Cache the result
	rs.Cache.Store("assigned_issues", report, dateRange)

	return report, nil
}

// GetWorkedIssues retrieves all issues the user has worked on within the date range
func (rs *ReportService) GetWorkedIssues(dateRange DateRange) (*WorkedIssuesReport, error) {
	// Check cache first
	if rs.Cache.IsValid(dateRange) {
		if cached, exists := rs.Cache.Get("worked_issues"); exists {
			if report, ok := cached.(*WorkedIssuesReport); ok {
				return report, nil
			}
		}
	}

	// Build JQL query for issues with worklogs by current user
	jql := fmt.Sprintf(
		"worklogAuthor = currentUser() AND worklogDate >= \"%s\" AND worklogDate <= \"%s\" ORDER BY key ASC",
		dateRange.StartDate.Format("2006-01-02"),
		dateRange.EndDate.Format("2006-01-02"),
	)

	// Fetch issues from Jira
	issues, err := rs.fetchAllIssues(jql)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch worked issues: %w", err)
	}

	// Convert to worked issue items with time aggregation
	reportItems := make([]WorkedIssueItem, 0, len(issues))
	totalTimeLogged := time.Duration(0)

	for _, issue := range issues {
		item, err := rs.convertToWorkedIssueItem(issue, dateRange)
		if err != nil {
			log.Printf("Warning: failed to convert issue %s: %v", issue.Key, err)
			continue
		}
		reportItems = append(reportItems, item)
		totalTimeLogged += item.TimeLogged
	}

	report := &WorkedIssuesReport{
		Issues:          reportItems,
		TotalCount:      len(reportItems),
		TotalTimeLogged: totalTimeLogged,
	}

	// Cache the result
	rs.Cache.Store("worked_issues", report, dateRange)

	return report, nil
}

// GetUntrackedIssues retrieves assigned issues without worklogs
func (rs *ReportService) GetUntrackedIssues(dateRange DateRange) (*UntrackedIssuesReport, error) {
	// Check cache first
	if rs.Cache.IsValid(dateRange) {
		if cached, exists := rs.Cache.Get("untracked_issues"); exists {
			if report, ok := cached.(*UntrackedIssuesReport); ok {
				return report, nil
			}
		}
	}

	// Get assigned and worked issues
	assignedReport, err := rs.GetAssignedIssues(dateRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned issues: %w", err)
	}

	workedReport, err := rs.GetWorkedIssues(dateRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get worked issues: %w", err)
	}

	// Create a set of worked issue keys
	workedKeys := make(map[string]bool)
	for _, issue := range workedReport.Issues {
		workedKeys[issue.Key] = true
	}

	// Find untracked issues (assigned but not worked)
	untrackedItems := make([]IssueReportItem, 0)
	for _, issue := range assignedReport.Issues {
		// Check if issue was open during the date range
		if !IsIssueOpenDuringRange(issue, dateRange) {
			continue
		}

		// Check if issue has worklogs in the date range
		if !workedKeys[issue.Key] {
			// Also check if issue has worklogs outside the range
			hasWorklogsInRange := rs.hasWorklogsInDateRange(issue.Key, dateRange)
			if !hasWorklogsInRange {
				untrackedItems = append(untrackedItems, issue)
			}
		}
	}

	report := &UntrackedIssuesReport{
		Issues:     untrackedItems,
		TotalCount: len(untrackedItems),
	}

	// Cache the result
	rs.Cache.Store("untracked_issues", report, dateRange)

	return report, nil
}

// GetRollupReport generates summary statistics
func (rs *ReportService) GetRollupReport(dateRange DateRange) (*RollupReport, error) {
	// Check cache first
	if rs.Cache.IsValid(dateRange) {
		if cached, exists := rs.Cache.Get("rollup_report"); exists {
			if report, ok := cached.(*RollupReport); ok {
				return report, nil
			}
		}
	}

	// Get all report types
	assignedReport, err := rs.GetAssignedIssues(dateRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned issues: %w", err)
	}

	workedReport, err := rs.GetWorkedIssues(dateRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get worked issues: %w", err)
	}

	// Calculate unique issues across all categories
	uniqueIssues := make(map[string]string) // key -> status
	for _, issue := range assignedReport.Issues {
		uniqueIssues[issue.Key] = issue.Status
	}
	for _, issue := range workedReport.Issues {
		uniqueIssues[issue.Key] = issue.Status
	}

	// Calculate status breakdown
	statusCounts := make(map[string]int)
	for _, status := range uniqueIssues {
		statusCounts[status]++
	}

	// Convert to StatusCount slice with percentages
	statusBreakdown := make([]StatusCount, 0, len(statusCounts))
	totalUniqueIssues := len(uniqueIssues)
	for statusName, count := range statusCounts {
		percentage := 0.0
		if totalUniqueIssues > 0 {
			percentage = float64(count) / float64(totalUniqueIssues) * 100.0
		}
		statusBreakdown = append(statusBreakdown, StatusCount{
			StatusName: statusName,
			Count:      count,
			Percentage: percentage,
		})
	}

	// Sort by count descending
	sort.Slice(statusBreakdown, func(i, j int) bool {
		return statusBreakdown[i].Count > statusBreakdown[j].Count
	})

	report := &RollupReport{
		TotalAssigned:     assignedReport.TotalCount,
		TotalWorked:       workedReport.TotalCount,
		StatusBreakdown:   statusBreakdown,
		TotalUniqueIssues: totalUniqueIssues,
	}

	log.Printf("Rollup Report: TotalAssigned=%d, TotalWorked=%d, UniqueIssues=%d, StatusBreakdown=%d items",
		report.TotalAssigned, report.TotalWorked, report.TotalUniqueIssues, len(report.StatusBreakdown))
	for _, status := range report.StatusBreakdown {
		log.Printf("  Status: %s, Count: %d, Percentage: %.1f%%", status.StatusName, status.Count, status.Percentage)
	}

	// Cache the result
	rs.Cache.Store("rollup_report", report, dateRange)

	return report, nil
}

// fetchAllIssues fetches all issues matching the JQL query with pagination
func (rs *ReportService) fetchAllIssues(jql string) ([]JiraIssue, error) {
	allIssues := make([]JiraIssue, 0)
	startAt := 0
	maxResults := 100

	for {
		fields := []string{"summary", "status", "assignee", "created", "updated", "resolutiondate"}
		
		// Use GET endpoint with query parameters (the old SearchIssues function)
		response, err := jiraApiFunctions.SearchIssues(jql, "", fields, startAt, maxResults, false)
		if err != nil {
			log.Printf("SearchIssues API error: %v", err)
			return nil, err
		}

		log.Printf("API Response (first 500 chars): %s", string(response[:min(500, len(response))]))

		var searchResponse JiraSearchResponse
		if err := json.Unmarshal(response, &searchResponse); err != nil {
			log.Printf("Failed to unmarshal response: %v", err)
			return nil, fmt.Errorf("failed to parse search response: %w", err)
		}

		log.Printf("Parsed response: Total=%d, StartAt=%d, Issues=%d", searchResponse.Total, searchResponse.StartAt, len(searchResponse.Issues))

		allIssues = append(allIssues, searchResponse.Issues...)

		// Check if we've fetched all issues
		if startAt+len(searchResponse.Issues) >= searchResponse.Total {
			break
		}

		startAt += maxResults
	}

	return allIssues, nil
}

// parseJiraDate parses a Jira date string which includes milliseconds
func parseJiraDate(dateStr string) (time.Time, error) {
	// Jira dates are in format: 2025-11-24T13:38:08.046-0700
	// Try RFC3339Nano first (includes milliseconds)
	t, err := time.Parse(time.RFC3339Nano, dateStr)
	if err == nil {
		return t, nil
	}
	
	// Fallback to RFC3339
	t, err = time.Parse(time.RFC3339, dateStr)
	if err == nil {
		return t, nil
	}
	
	// Try custom format with milliseconds
	t, err = time.Parse("2006-01-02T15:04:05.999-0700", dateStr)
	return t, err
}

// convertToReportItem converts a JiraIssue to an IssueReportItem
func (rs *ReportService) convertToReportItem(issue JiraIssue) (IssueReportItem, error) {
	createdDate, err := parseJiraDate(issue.Fields.Created)
	if err != nil {
		return IssueReportItem{}, fmt.Errorf("failed to parse created date: %w", err)
	}

	var resolvedDate *time.Time
	if issue.Fields.Resolutiondate != nil && *issue.Fields.Resolutiondate != "" {
		parsed, err := parseJiraDate(*issue.Fields.Resolutiondate)
		if err != nil {
			log.Printf("Warning: failed to parse resolution date for %s: %v", issue.Key, err)
		} else {
			resolvedDate = &parsed
		}
	}

	return IssueReportItem{
		Key:          issue.Key,
		Summary:      issue.Fields.Summary,
		Status:       issue.Fields.Status.Name,
		AssignedDate: createdDate, // Using created as assigned date
		CreatedDate:  createdDate,
		ResolvedDate: resolvedDate,
	}, nil
}

// convertToWorkedIssueItem converts a JiraIssue to a WorkedIssueItem with time aggregation
func (rs *ReportService) convertToWorkedIssueItem(issue JiraIssue, dateRange DateRange) (WorkedIssueItem, error) {
	baseItem, err := rs.convertToReportItem(issue)
	if err != nil {
		return WorkedIssueItem{}, err
	}

	// Fetch worklogs for this issue
	worklogs, err := rs.fetchWorklogsForIssue(issue.Key, dateRange)
	if err != nil {
		return WorkedIssueItem{}, fmt.Errorf("failed to fetch worklogs: %w", err)
	}

	// Aggregate time logged
	totalSeconds := 0
	worklogCount := 0
	for _, worklog := range worklogs {
		if worklog.Author.AccountId == rs.CurrentUser.AccountId {
			totalSeconds += worklog.TimeSpentSeconds
			worklogCount++
		}
	}

	return WorkedIssueItem{
		IssueReportItem: baseItem,
		TimeLogged:      time.Duration(totalSeconds) * time.Second,
		WorklogCount:    worklogCount,
	}, nil
}

// fetchWorklogsForIssue fetches all worklogs for an issue within the date range
func (rs *ReportService) fetchWorklogsForIssue(issueKey string, dateRange DateRange) ([]JiraWorklog, error) {
	allWorklogs := make([]JiraWorklog, 0)
	startAt := 0
	maxResults := 100

	for {
		response, err := jiraApiFunctions.GetIssueWorklog(issueKey, startAt, maxResults, "")
		if err != nil {
			return nil, err
		}

		var worklogResponse JiraWorklogResponse
		if err := json.Unmarshal(response, &worklogResponse); err != nil {
			return nil, fmt.Errorf("failed to parse worklog response: %w", err)
		}

		// Filter worklogs by date range
		for _, worklog := range worklogResponse.Worklogs {
			startedTime, err := parseJiraDate(worklog.Started)
			if err != nil {
				log.Printf("Warning: failed to parse worklog started time: %v", err)
				continue
			}

			// Check if worklog is within date range
			if !startedTime.Before(dateRange.StartDate) && !startedTime.After(dateRange.EndDate) {
				allWorklogs = append(allWorklogs, worklog)
			}
		}

		// Check if we've fetched all worklogs
		if startAt+len(worklogResponse.Worklogs) >= worklogResponse.Total {
			break
		}

		startAt += maxResults
	}

	return allWorklogs, nil
}

// hasWorklogsInDateRange checks if an issue has any worklogs by the current user in the date range
func (rs *ReportService) hasWorklogsInDateRange(issueKey string, dateRange DateRange) bool {
	worklogs, err := rs.fetchWorklogsForIssue(issueKey, dateRange)
	if err != nil {
		log.Printf("Warning: failed to check worklogs for %s: %v", issueKey, err)
		return false
	}

	for _, worklog := range worklogs {
		if worklog.Author.AccountId == rs.CurrentUser.AccountId {
			return true
		}
	}

	return false
}

// InvalidateCache clears all cached report data
func (rs *ReportService) InvalidateCache() {
	rs.Cache.Invalidate()
}
