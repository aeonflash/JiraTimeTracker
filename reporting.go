package main

import (
	"fmt"
	"sync"
	"time"
)

// DateRange represents a period defined by a start date and end date
type DateRange struct {
	StartDate time.Time
	EndDate   time.Time
}

// Validate checks if the date range is valid (start <= end)
func (dr DateRange) Validate() error {
	if dr.StartDate.After(dr.EndDate) {
		return fmt.Errorf("start date must be before or equal to end date")
	}
	return nil
}

// GetDefaultDateRange returns the current month as the default date range
func GetDefaultDateRange() DateRange {
	now := time.Now()
	year, month, _ := now.Date()
	location := now.Location()
	
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, location)
	endDate := startDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	
	return DateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}
}

// IsIssueOpenDuringRange checks if an issue was open during the date range
// An issue is considered open during the range if:
// - It was created on or before the range end date, AND
// - It is still open (ResolvedDate is nil) OR it was closed on or after the range start date
func IsIssueOpenDuringRange(issue IssueReportItem, dateRange DateRange) bool {
	// Issue must have been created on or before the range end date
	if issue.CreatedDate.After(dateRange.EndDate) {
		return false
	}
	
	// If issue is still open (no resolved date), it's open during the range
	if issue.ResolvedDate == nil {
		return true
	}
	
	// If issue was closed, it must have been closed on or after the range start date
	return !issue.ResolvedDate.Before(dateRange.StartDate)
}

// IssueReportItem represents a Jira issue in a report
type IssueReportItem struct {
	Key          string
	Summary      string
	Status       string
	AssignedDate time.Time
	CreatedDate  time.Time
	ResolvedDate *time.Time
}

// WorkedIssueItem represents a Jira issue with work log information
type WorkedIssueItem struct {
	IssueReportItem
	TimeLogged   time.Duration
	WorklogCount int
}

// StatusCount represents the count and percentage of issues in a specific status
type StatusCount struct {
	StatusName string
	Count      int
	Percentage float64
}

// AssignedIssuesReport contains the list of issues assigned to the user
type AssignedIssuesReport struct {
	Issues     []IssueReportItem
	TotalCount int
}

// WorkedIssuesReport contains the list of issues the user has worked on
type WorkedIssuesReport struct {
	Issues          []WorkedIssueItem
	TotalCount      int
	TotalTimeLogged time.Duration
}

// UntrackedIssuesReport contains the list of assigned issues without work logs
type UntrackedIssuesReport struct {
	Issues     []IssueReportItem
	TotalCount int
}

// RollupReport contains summary statistics across all report types
type RollupReport struct {
	TotalAssigned     int
	TotalWorked       int
	StatusBreakdown   []StatusCount
	TotalUniqueIssues int
}

// ReportCache stores cached report data with thread safety
type ReportCache struct {
	Data      map[string]interface{}
	Timestamp time.Time
	DateRange DateRange
	mutex     sync.RWMutex
}

// NewReportCache creates a new report cache instance
func NewReportCache() *ReportCache {
	return &ReportCache{
		Data: make(map[string]interface{}),
	}
}

// Store saves data to the cache with the current timestamp
func (rc *ReportCache) Store(key string, value interface{}, dateRange DateRange) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	
	rc.Data[key] = value
	rc.Timestamp = time.Now()
	rc.DateRange = dateRange
}

// Get retrieves data from the cache
func (rc *ReportCache) Get(key string) (interface{}, bool) {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()
	
	value, exists := rc.Data[key]
	return value, exists
}

// IsValid checks if the cached data is still valid
// Data is valid if:
// - It's less than 5 minutes old
// - The date range matches the requested date range
func (rc *ReportCache) IsValid(dateRange DateRange) bool {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()
	
	// Check if cache is older than 5 minutes
	if time.Since(rc.Timestamp) > 5*time.Minute {
		return false
	}
	
	// Check if date range matches
	if !rc.DateRange.StartDate.Equal(dateRange.StartDate) || !rc.DateRange.EndDate.Equal(dateRange.EndDate) {
		return false
	}
	
	return true
}

// Invalidate clears all cached data
func (rc *ReportCache) Invalidate() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	
	rc.Data = make(map[string]interface{})
	rc.Timestamp = time.Time{}
	rc.DateRange = DateRange{}
}

// Clear removes all cached data (alias for Invalidate for clarity)
func (rc *ReportCache) Clear() {
	rc.Invalidate()
}
