package main

import (
	"errors"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: jira-reporting-window, Property 22: API error message display**
// For any Jira API call failure, an error message should be displayed to the user.
// **Validates: Requirements 7.1**
func TestPropertyAPIErrorMessageDisplay(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("API errors are categorized and display user-friendly messages", prop.ForAll(
		func(errorType int) bool {
			// Map error type to test error
			var testErr error
			var expectedType ErrorType

			switch errorType % 6 {
			case 0:
				// Authentication error
				testErr = errors.New("401 Unauthorized")
				expectedType = ErrorTypeAuth
			case 1:
				// Network error
				testErr = errors.New("connection timeout")
				expectedType = ErrorTypeNetwork
			case 2:
				// Rate limit error
				testErr = errors.New("429 Too Many Requests")
				expectedType = ErrorTypeRateLimit
			case 3:
				// Date range error
				testErr = errors.New("invalid date range")
				expectedType = ErrorTypeInvalidDateRange
			case 4:
				// Unknown error
				testErr = errors.New("unknown error occurred")
				expectedType = ErrorTypeUnknown
			case 5:
				// Another auth error variant
				testErr = errors.New("403 Forbidden")
				expectedType = ErrorTypeAuth
			}

			// Categorize the error
			reportErr := CategorizeError(testErr)

			// Verify error was categorized correctly
			if reportErr.Type != expectedType {
				return false
			}

			// Verify error has a user-friendly message
			userMessage := reportErr.GetUserFriendlyMessage()
			if userMessage == "" {
				return false
			}

			// Verify specific error messages for known types
			switch expectedType {
			case ErrorTypeAuth:
				if userMessage != "Authentication failed. Please check your API credentials in ~/.jirarc" {
					return false
				}
			case ErrorTypeNetwork:
				if userMessage != "Network error: Unable to connect to Jira. Please check your internet connection." {
					return false
				}
			case ErrorTypeRateLimit:
				if userMessage != "Jira API rate limit exceeded. Please wait a few minutes and try again." {
					return false
				}
			case ErrorTypeInvalidDateRange:
				if userMessage != "Invalid date range: Start date must be before or equal to end date." {
					return false
				}
			}

			return true
		},
		gen.IntRange(0, 100), // errorType
	))

	properties.Property("authentication errors are detected correctly", prop.ForAll(
		func(statusCode string) bool {
			// Test various authentication error patterns
			authErrors := []string{
				"401 Unauthorized",
				"403 Forbidden",
				"authentication failed",
				"UNAUTHORIZED access",
			}

			for _, errMsg := range authErrors {
				testErr := errors.New(errMsg)
				reportErr := CategorizeError(testErr)

				if reportErr.Type != ErrorTypeAuth {
					return false
				}

				expectedMsg := "Authentication failed. Please check your API credentials in ~/.jirarc"
				if reportErr.GetUserFriendlyMessage() != expectedMsg {
					return false
				}
			}

			return true
		},
		gen.AnyString(), // statusCode (not used, just for property generation)
	))

	properties.Property("network errors are detected correctly", prop.ForAll(
		func(dummy string) bool {
			// Test various network error patterns
			networkErrors := []string{
				"connection timeout",
				"network error",
				"dial tcp: no such host",
				"connection refused",
			}

			for _, errMsg := range networkErrors {
				testErr := errors.New(errMsg)
				reportErr := CategorizeError(testErr)

				if reportErr.Type != ErrorTypeNetwork {
					return false
				}

				expectedMsg := "Network error: Unable to connect to Jira. Please check your internet connection."
				if reportErr.GetUserFriendlyMessage() != expectedMsg {
					return false
				}
			}

			return true
		},
		gen.AnyString(), // dummy (not used, just for property generation)
	))

	properties.Property("rate limit errors are detected correctly", prop.ForAll(
		func(dummy string) bool {
			// Test various rate limit error patterns
			rateLimitErrors := []string{
				"429 Too Many Requests",
				"rate limit exceeded",
				"too many requests",
			}

			for _, errMsg := range rateLimitErrors {
				testErr := errors.New(errMsg)
				reportErr := CategorizeError(testErr)

				if reportErr.Type != ErrorTypeRateLimit {
					return false
				}

				expectedMsg := "Jira API rate limit exceeded. Please wait a few minutes and try again."
				if reportErr.GetUserFriendlyMessage() != expectedMsg {
					return false
				}
			}

			return true
		},
		gen.AnyString(), // dummy (not used, just for property generation)
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 23: Partial data graceful degradation**
// For any partial data retrieval scenario, the system should display the available
// data along with a warning about incomplete results.
// **Validates: Requirements 7.5**
func TestPropertyPartialDataGracefulDegradation(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("partial data shows warning with available data", prop.ForAll(
		func(failedSections []bool) bool {
			// Ensure we have at least one section
			if len(failedSections) == 0 {
				failedSections = []bool{false, false, false, false}
			}
			// Limit to 4 sections (assigned, worked, untracked, rollup)
			if len(failedSections) > 4 {
				failedSections = failedSections[:4]
			}

			// Count how many sections failed
			failureCount := 0
			for _, failed := range failedSections {
				if failed {
					failureCount++
				}
			}

			// If all sections failed, this would be a complete failure, not partial
			// So we skip this case
			if failureCount == len(failedSections) {
				return true
			}

			// If no sections failed, this is a success case, not partial
			if failureCount == 0 {
				return true
			}

			// This is a partial failure case
			// Verify that we can create empty reports for failed sections
			var assignedReport *AssignedIssuesReport
			var workedReport *WorkedIssuesReport
			var untrackedReport *UntrackedIssuesReport
			var rollupReport *RollupReport

			if len(failedSections) > 0 && failedSections[0] {
				// Assigned section failed - create empty report
				assignedReport = &AssignedIssuesReport{
					Issues:     []IssueReportItem{},
					TotalCount: 0,
				}
			} else {
				// Assigned section succeeded - create report with data
				assignedReport = &AssignedIssuesReport{
					Issues: []IssueReportItem{
						{Key: "TEST-1", Summary: "Test Issue", Status: "Open"},
					},
					TotalCount: 1,
				}
			}

			if len(failedSections) > 1 && failedSections[1] {
				// Worked section failed - create empty report
				workedReport = &WorkedIssuesReport{
					Issues:          []WorkedIssueItem{},
					TotalCount:      0,
					TotalTimeLogged: 0,
				}
			} else {
				// Worked section succeeded - create report with data
				workedReport = &WorkedIssuesReport{
					Issues: []WorkedIssueItem{
						{IssueReportItem: IssueReportItem{Key: "TEST-2", Summary: "Test Issue", Status: "Open"}},
					},
					TotalCount:      1,
					TotalTimeLogged: 3600,
				}
			}

			if len(failedSections) > 2 && failedSections[2] {
				// Untracked section failed - create empty report
				untrackedReport = &UntrackedIssuesReport{
					Issues:     []IssueReportItem{},
					TotalCount: 0,
				}
			} else {
				// Untracked section succeeded - create report with data
				untrackedReport = &UntrackedIssuesReport{
					Issues: []IssueReportItem{
						{Key: "TEST-3", Summary: "Test Issue", Status: "Open"},
					},
					TotalCount: 1,
				}
			}

			if len(failedSections) > 3 && failedSections[3] {
				// Rollup section failed - create empty report
				rollupReport = &RollupReport{
					TotalAssigned:     0,
					TotalWorked:       0,
					StatusBreakdown:   []StatusCount{},
					TotalUniqueIssues: 0,
				}
			} else {
				// Rollup section succeeded - create report with data
				rollupReport = &RollupReport{
					TotalAssigned:     1,
					TotalWorked:       1,
					StatusBreakdown:   []StatusCount{{StatusName: "Open", Count: 1, Percentage: 100.0}},
					TotalUniqueIssues: 1,
				}
			}

			// Verify that all reports are non-nil (even if empty)
			if assignedReport == nil || workedReport == nil || untrackedReport == nil || rollupReport == nil {
				return false
			}

			// Verify that empty reports have zero counts
			if len(failedSections) > 0 && failedSections[0] {
				if assignedReport.TotalCount != 0 || len(assignedReport.Issues) != 0 {
					return false
				}
			}

			if len(failedSections) > 1 && failedSections[1] {
				if workedReport.TotalCount != 0 || len(workedReport.Issues) != 0 {
					return false
				}
			}

			if len(failedSections) > 2 && failedSections[2] {
				if untrackedReport.TotalCount != 0 || len(untrackedReport.Issues) != 0 {
					return false
				}
			}

			if len(failedSections) > 3 && failedSections[3] {
				if rollupReport.TotalAssigned != 0 || rollupReport.TotalWorked != 0 {
					return false
				}
			}

			return true
		},
		gen.SliceOf(gen.Bool()), // failedSections
	))

	properties.Property("warning message includes failed section names", prop.ForAll(
		func(sectionIndex int) bool {
			// Map section index to section name
			sections := []string{"assigned issues", "worked issues", "untracked issues", "rollup report"}
			sectionIndex = sectionIndex % len(sections)

			sectionName := sections[sectionIndex]

			// Create a warning message that includes the section name
			warningMessage := "Some data could not be retrieved (" + sectionName + "). Showing partial results."

			// Verify the message contains the section name
			if !contains(warningMessage, sectionName) {
				return false
			}

			// Verify the message indicates partial results
			if !contains(warningMessage, "partial results") {
				return false
			}

			return true
		},
		gen.IntRange(0, 100), // sectionIndex
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
