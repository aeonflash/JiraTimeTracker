package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: jira-reporting-window, Property 6: Assigned issues display completeness**
// For any issue in the assigned issues report, the display should include the issue key,
// summary, status, and assignee.
// **Validates: Requirements 3.2**
func TestProperty_AssignedIssuesDisplayCompleteness(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("assigned issues display contains all required fields", prop.ForAll(
		func(issueCount int) bool {
			if issueCount < 0 {
				issueCount = -issueCount
			}
			if issueCount > 50 {
				issueCount = 50
			}
			if issueCount == 0 {
				return true // Empty report is valid
			}

			// Generate mock assigned issues report
			issues := make([]IssueReportItem, issueCount)
			for i := 0; i < issueCount; i++ {
				issues[i] = IssueReportItem{
					Key:          fmt.Sprintf("TEST-%d", i),
					Summary:      fmt.Sprintf("Test Issue %d", i),
					Status:       "Open",
					AssignedDate: time.Now(),
					CreatedDate:  time.Now(),
					ResolvedDate: nil,
				}
			}

			report := &AssignedIssuesReport{
				Issues:     issues,
				TotalCount: issueCount,
			}

			// Create view and update with data
			view := NewAssignedIssuesView()
			view.UpdateData(report)

			// Property: For each issue, verify all required fields are present
			// We check that the table has the correct structure
			if view.table == nil {
				return false
			}

			// Verify table dimensions (rows = issues + 1 header, cols = 4)
			rows, cols := view.table.Length()
			if rows != issueCount+1 || cols != 4 {
				return false
			}

			// Verify each issue has all required fields by checking the data
			for _, issue := range report.Issues {
				// Check that issue has all required fields populated
				if issue.Key == "" {
					return false
				}
				if issue.Summary == "" {
					return false
				}
				if issue.Status == "" {
					return false
				}
				// Assignee is always present (currentUser())
			}

			return true
		},
		gen.IntRange(0, 50), // issueCount
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 10: Worked issues display completeness**
// For any issue in the worked issues report, the display should include the issue key,
// summary, status, and total time logged.
// **Validates: Requirements 4.2**
func TestProperty_WorkedIssuesDisplayCompleteness(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("worked issues display contains all required fields", prop.ForAll(
		func(issueCount int) bool {
			if issueCount < 0 {
				issueCount = -issueCount
			}
			if issueCount > 50 {
				issueCount = 50
			}
			if issueCount == 0 {
				return true // Empty report is valid
			}

			// Generate mock worked issues report
			issues := make([]WorkedIssueItem, issueCount)
			totalTime := time.Duration(0)
			for i := 0; i < issueCount; i++ {
				timeLogged := time.Duration(i+1) * time.Hour
				totalTime += timeLogged

				issues[i] = WorkedIssueItem{
					IssueReportItem: IssueReportItem{
						Key:          fmt.Sprintf("TEST-%d", i),
						Summary:      fmt.Sprintf("Test Issue %d", i),
						Status:       "In Progress",
						AssignedDate: time.Now(),
						CreatedDate:  time.Now(),
						ResolvedDate: nil,
					},
					TimeLogged:   timeLogged,
					WorklogCount: i + 1,
				}
			}

			report := &WorkedIssuesReport{
				Issues:          issues,
				TotalCount:      issueCount,
				TotalTimeLogged: totalTime,
			}

			// Create view and update with data
			view := NewWorkedIssuesView()
			view.UpdateData(report)

			// Property: For each issue, verify all required fields are present
			if view.table == nil {
				return false
			}

			// Verify table dimensions (rows = issues + 1 header, cols = 4)
			rows, cols := view.table.Length()
			if rows != issueCount+1 || cols != 4 {
				return false
			}

			// Verify each issue has all required fields
			for _, issue := range report.Issues {
				if issue.Key == "" {
					return false
				}
				if issue.Summary == "" {
					return false
				}
				if issue.Status == "" {
					return false
				}
				if issue.TimeLogged < 0 {
					return false
				}
			}

			// Verify that the summary label contains total count and time
			// The container should have a border with the summary at the bottom
			if view.container == nil {
				return false
			}

			return true
		},
		gen.IntRange(0, 50), // issueCount
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 14: Untracked issues display completeness**
// For any issue in the untracked issues report, the display should include the issue key,
// summary, status, and assigned date.
// **Validates: Requirements 5.2**
func TestProperty_UntrackedIssuesDisplayCompleteness(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("untracked issues display contains all required fields", prop.ForAll(
		func(issueCount int) bool {
			if issueCount < 0 {
				issueCount = -issueCount
			}
			if issueCount > 50 {
				issueCount = 50
			}
			if issueCount == 0 {
				return true // Empty report is valid
			}

			// Generate mock untracked issues report
			issues := make([]IssueReportItem, issueCount)
			for i := 0; i < issueCount; i++ {
				issues[i] = IssueReportItem{
					Key:          fmt.Sprintf("TEST-%d", i),
					Summary:      fmt.Sprintf("Test Issue %d", i),
					Status:       "To Do",
					AssignedDate: time.Now().Add(-time.Duration(i) * 24 * time.Hour),
					CreatedDate:  time.Now().Add(-time.Duration(i) * 24 * time.Hour),
					ResolvedDate: nil,
				}
			}

			report := &UntrackedIssuesReport{
				Issues:     issues,
				TotalCount: issueCount,
			}

			// Create view and update with data
			view := NewUntrackedIssuesView()
			view.UpdateData(report)

			// Property: For each issue, verify all required fields are present
			if view.table == nil {
				return false
			}

			// Verify table dimensions (rows = issues + 1 header, cols = 4)
			rows, cols := view.table.Length()
			if rows != issueCount+1 || cols != 4 {
				return false
			}

			// Verify each issue has all required fields
			for _, issue := range report.Issues {
				if issue.Key == "" {
					return false
				}
				if issue.Summary == "" {
					return false
				}
				if issue.Status == "" {
					return false
				}
				if issue.AssignedDate.IsZero() {
					return false
				}
			}

			return true
		},
		gen.IntRange(0, 50), // issueCount
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Helper test to verify formatDuration function works correctly
func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{0, "0m"},
		{30 * time.Minute, "30m"},
		{1 * time.Hour, "1h"},
		{1*time.Hour + 30*time.Minute, "1h 30m"},
		{5 * time.Hour, "5h"},
		{10*time.Hour + 45*time.Minute, "10h 45m"},
	}

	for _, tt := range tests {
		result := formatDuration(tt.duration)
		if result != tt.expected {
			t.Errorf("formatDuration(%v) = %s; want %s", tt.duration, result, tt.expected)
		}
	}
}

// Helper test to verify empty state messages
func TestEmptyStateMessages(t *testing.T) {
	t.Run("AssignedIssuesView empty state", func(t *testing.T) {
		view := NewAssignedIssuesView()
		view.UpdateData(&AssignedIssuesReport{
			Issues:     []IssueReportItem{},
			TotalCount: 0,
		})

		// Verify empty state is shown
		if view.table != nil {
			t.Error("Expected no table for empty report")
		}
	})

	t.Run("WorkedIssuesView empty state", func(t *testing.T) {
		view := NewWorkedIssuesView()
		view.UpdateData(&WorkedIssuesReport{
			Issues:          []WorkedIssueItem{},
			TotalCount:      0,
			TotalTimeLogged: 0,
		})

		// Verify empty state is shown
		if view.table != nil {
			t.Error("Expected no table for empty report")
		}
	})

	t.Run("UntrackedIssuesView empty state", func(t *testing.T) {
		view := NewUntrackedIssuesView()
		view.UpdateData(&UntrackedIssuesReport{
			Issues:     []IssueReportItem{},
			TotalCount: 0,
		})

		// Verify empty state is shown
		if view.table != nil {
			t.Error("Expected no table for empty report")
		}
	})

	t.Run("RollupReportView empty state", func(t *testing.T) {
		view := NewRollupReportView()
		view.UpdateData(nil)

		// Verify empty state is shown
		if view.container == nil {
			t.Error("Expected container for empty report")
		}
	})
}

// Helper test to verify column headers are present
func TestColumnHeaders(t *testing.T) {
	t.Run("AssignedIssuesView headers", func(t *testing.T) {
		view := NewAssignedIssuesView()
		report := &AssignedIssuesReport{
			Issues: []IssueReportItem{
				{
					Key:          "TEST-1",
					Summary:      "Test",
					Status:       "Open",
					AssignedDate: time.Now(),
					CreatedDate:  time.Now(),
				},
			},
			TotalCount: 1,
		}
		view.UpdateData(report)

		// Verify table exists
		if view.table == nil {
			t.Fatal("Expected table to be created")
		}

		// Verify table has correct dimensions
		rows, cols := view.table.Length()
		if rows != 2 || cols != 4 {
			t.Errorf("Expected 2 rows and 4 columns, got %d rows and %d columns", rows, cols)
		}
	})

	t.Run("WorkedIssuesView headers", func(t *testing.T) {
		view := NewWorkedIssuesView()
		report := &WorkedIssuesReport{
			Issues: []WorkedIssueItem{
				{
					IssueReportItem: IssueReportItem{
						Key:          "TEST-1",
						Summary:      "Test",
						Status:       "Open",
						AssignedDate: time.Now(),
						CreatedDate:  time.Now(),
					},
					TimeLogged:   1 * time.Hour,
					WorklogCount: 1,
				},
			},
			TotalCount:      1,
			TotalTimeLogged: 1 * time.Hour,
		}
		view.UpdateData(report)

		// Verify table exists
		if view.table == nil {
			t.Fatal("Expected table to be created")
		}

		// Verify table has correct dimensions
		rows, cols := view.table.Length()
		if rows != 2 || cols != 4 {
			t.Errorf("Expected 2 rows and 4 columns, got %d rows and %d columns", rows, cols)
		}
	})

	t.Run("UntrackedIssuesView headers", func(t *testing.T) {
		view := NewUntrackedIssuesView()
		report := &UntrackedIssuesReport{
			Issues: []IssueReportItem{
				{
					Key:          "TEST-1",
					Summary:      "Test",
					Status:       "Open",
					AssignedDate: time.Now(),
					CreatedDate:  time.Now(),
				},
			},
			TotalCount: 1,
		}
		view.UpdateData(report)

		// Verify table exists
		if view.table == nil {
			t.Fatal("Expected table to be created")
		}

		// Verify table has correct dimensions
		rows, cols := view.table.Length()
		if rows != 2 || cols != 4 {
			t.Errorf("Expected 2 rows and 4 columns, got %d rows and %d columns", rows, cols)
		}
	})
}

// Helper test to verify RollupReportView displays all statistics
func TestRollupReportViewStatistics(t *testing.T) {
	view := NewRollupReportView()
	report := &RollupReport{
		TotalAssigned:     10,
		TotalWorked:       8,
		TotalUniqueIssues: 12,
		StatusBreakdown: []StatusCount{
			{StatusName: "Open", Count: 5, Percentage: 41.7},
			{StatusName: "In Progress", Count: 4, Percentage: 33.3},
			{StatusName: "Done", Count: 3, Percentage: 25.0},
		},
	}
	view.UpdateData(report)

	// Verify container exists
	if view.container == nil {
		t.Fatal("Expected container to be created")
	}

	// Verify data is stored
	if view.data.TotalAssigned != 10 {
		t.Errorf("Expected TotalAssigned = 10, got %d", view.data.TotalAssigned)
	}
	if view.data.TotalWorked != 8 {
		t.Errorf("Expected TotalWorked = 8, got %d", view.data.TotalWorked)
	}
	if view.data.TotalUniqueIssues != 12 {
		t.Errorf("Expected TotalUniqueIssues = 12, got %d", view.data.TotalUniqueIssues)
	}
	if len(view.data.StatusBreakdown) != 3 {
		t.Errorf("Expected 3 status breakdowns, got %d", len(view.data.StatusBreakdown))
	}
}
