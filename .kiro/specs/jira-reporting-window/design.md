# Design Document: Jira Reporting Window

## Overview

The Jira Reporting Window feature extends the JiraWidgetLite application by adding a comprehensive reporting interface that allows users to analyze their Jira work activity over customizable date ranges. The feature will be implemented as a separate window that can be opened from the main application, providing various report views including assigned issues, worked issues, untracked issues, and rollup statistics.

The design leverages the existing Jira REST API integration and follows the established patterns in the codebase using Go with the Fyne UI framework. The reporting window will operate independently of the main time tracking window, allowing users to generate reports without disrupting their workflow.

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Main Application                         │
│  ┌──────────────────┐         ┌──────────────────────────┐ │
│  │  Main Window     │         │  Reporting Window        │ │
│  │  (Time Tracking) │◄────────┤  (Report Generation)     │ │
│  └──────────────────┘         └──────────────────────────┘ │
│                                         │                    │
│                                         ▼                    │
│                          ┌──────────────────────────┐       │
│                          │   Report Service         │       │
│                          │  - Data Fetching         │       │
│                          │  - Data Aggregation      │       │
│                          │  - Caching               │       │
│                          └──────────────────────────┘       │
│                                         │                    │
│                                         ▼                    │
│                          ┌──────────────────────────┐       │
│                          │   Jira API Client        │       │
│                          │  (Existing Integration)  │       │
│                          └──────────────────────────┘       │
└─────────────────────────────────────────────────────────────┘
                                          │
                                          ▼
                              ┌──────────────────────┐
                              │   Jira REST API v3   │
                              └──────────────────────┘
```

### Component Layers

1. **UI Layer**: Fyne-based reporting window with date range controls, report view tabs, and export functionality
2. **Service Layer**: Report generation logic, data aggregation, and caching
3. **API Layer**: Existing Jira API client for REST API calls
4. **Export Layer**: CSV and PDF generation utilities

## Components and Interfaces

### 1. Reporting Window Component

**Responsibility**: Manages the UI for the reporting window, including date range selection, report view tabs, and user interactions.

**Key Structures**:
```go
type ReportingWindow struct {
    Window          fyne.Window
    StartDatePicker *widget.Entry
    EndDatePicker   *widget.Entry
    ReportTabs      *container.AppTabs
    LoadingIndicator *widget.ProgressBarInfinite
    StatusLabel     *widget.Label
    ExportButton    *widget.Button
    RefreshButton   *widget.Button
    ReportService   *ReportService
    CurrentUser     *User
}

type DateRange struct {
    StartDate time.Time
    EndDate   time.Time
}
```

**Key Methods**:
- `NewReportingWindow(user *User) *ReportingWindow`: Creates and initializes the reporting window
- `Show()`: Displays the reporting window
- `UpdateDateRange(start, end time.Time)`: Updates the date range and refreshes reports
- `ShowLoading(message string)`: Displays loading indicator
- `HideLoading()`: Hides loading indicator
- `ShowError(message string)`: Displays error message to user

### 2. Report Service Component

**Responsibility**: Handles data fetching, aggregation, caching, and report generation logic.

**Key Structures**:
```go
type ReportService struct {
    APIClient     *jiraApiFunctions
    Cache         *ReportCache
    CurrentUser   *User
}

type ReportCache struct {
    Data          map[string]interface{}
    Timestamp     time.Time
    DateRange     DateRange
    mutex         sync.RWMutex
}

type AssignedIssuesReport struct {
    Issues    []IssueReportItem
    TotalCount int
}

type WorkedIssuesReport struct {
    Issues    []WorkedIssueItem
    TotalCount int
    TotalTimeLogged time.Duration
}

type UntrackedIssuesReport struct {
    Issues    []IssueReportItem
    TotalCount int
}

type RollupReport struct {
    TotalAssigned      int
    TotalWorked        int
    StatusBreakdown    []StatusCount
    TotalUniqueIssues  int
}

type IssueReportItem struct {
    Key         string
    Summary     string
    Status      string
    AssignedDate time.Time
    CreatedDate  time.Time
    ResolvedDate *time.Time
}

type WorkedIssueItem struct {
    IssueReportItem
    TimeLogged  time.Duration
    WorklogCount int
}

type StatusCount struct {
    StatusName  string
    Count       int
    Percentage  float64
}
```

**Key Methods**:
- `GetAssignedIssues(dateRange DateRange) (*AssignedIssuesReport, error)`: Fetches assigned issues within date range
- `GetWorkedIssues(dateRange DateRange) (*WorkedIssuesReport, error)`: Fetches issues with worklogs within date range
- `GetUntrackedIssues(dateRange DateRange) (*UntrackedIssuesReport, error)`: Identifies assigned issues without worklogs
- `GetRollupReport(dateRange DateRange) (*RollupReport, error)`: Generates summary statistics
- `InvalidateCache()`: Clears cached report data
- `IssueCachedDataValid(dateRange DateRange) bool`: Checks if cached data is still valid

### 3. Report View Components

**Responsibility**: Individual UI components for displaying different report types.

**Components**:
- `AssignedIssuesView`: Displays table of assigned issues
- `WorkedIssuesView`: Displays table of worked issues with time logged
- `UntrackedIssuesView`: Displays table of untracked issues
- `RollupReportView`: Displays summary statistics with charts/tables

Each view component will:
- Render data in a Fyne table widget
- Support sorting by columns
- Display appropriate empty state messages
- Handle loading and error states

### 4. Export Service Component

**Responsibility**: Generates CSV and PDF exports of report data.

**Key Structures**:
```go
type ExportService struct {
    ReportData interface{}
    DateRange  DateRange
    User       *User
    ExportDate time.Time
}

type ExportFormat int

const (
    ExportFormatCSV ExportFormat = iota
    ExportFormatPDF
)
```

**Key Methods**:
- `ExportToCSV(reports interface{}, filepath string) error`: Generates CSV export
- `ExportToPDF(reports interface{}, filepath string) error`: Generates PDF export
- `GenerateExportHeader() string`: Creates header with user, date range, and export date
- `FormatReportDataForExport(report interface{}) [][]string`: Converts report data to exportable format

### 5. Date Range Validator

**Responsibility**: Validates date range inputs and provides default date ranges.

**Key Methods**:
- `ValidateDateRange(start, end time.Time) error`: Validates that start <= end
- `GetDefaultDateRange() DateRange`: Returns current month as default
- `IsIssueOpenDuringRange(issue IssueReportItem, dateRange DateRange) bool`: Checks if issue was open during date range

## Data Models

### Issue Data Model

Issues are retrieved from Jira with the following fields:
- `key`: Issue identifier (e.g., "PROJ-123")
- `summary`: Issue title/description
- `status`: Current status object with name and category
- `created`: Issue creation timestamp
- `updated`: Last update timestamp
- `resolutiondate`: Resolution/close timestamp (null if open)
- `assignee`: Assigned user object

### Worklog Data Model

Worklogs are retrieved from Jira with the following fields:
- `id`: Worklog identifier
- `author`: User who logged the work
- `timeSpentSeconds`: Duration in seconds
- `started`: When the work was performed
- `created`: When the worklog was created

### Report Data Flow

1. **User selects date range** → Validate dates
2. **User requests report** → Check cache validity
3. **If cache invalid** → Fetch data from Jira API
4. **Aggregate and process data** → Generate report structures
5. **Cache results** → Store with timestamp
6. **Render in UI** → Display in appropriate view

## Data Models (Continued)

### JQL Query Patterns

The service will use JQL (Jira Query Language) to filter issues efficiently:

**Assigned Issues Query**:
```jql
assignee = currentUser() 
AND updated >= "YYYY-MM-DD" 
AND updated <= "YYYY-MM-DD"
ORDER BY key ASC
```

**Worked Issues Query**:
```jql
worklogAuthor = currentUser() 
AND worklogDate >= "YYYY-MM-DD" 
AND worklogDate <= "YYYY-MM-DD"
ORDER BY key ASC
```

**Issue Open During Range Logic**:
For each issue, check:
- `created <= dateRange.EndDate` (issue existed during or before range)
- `resolutiondate == null OR resolutiondate >= dateRange.StartDate` (issue was open during or after range start)

### API Endpoints Used

The implementation will use these existing Jira REST API v3 endpoints:

1. **Search Issues**: `GET /rest/api/3/search`
   - Parameters: `jql`, `fields`, `maxResults`, `startAt`
   - Returns paginated issue results

2. **Get Issue Worklogs**: `GET /rest/api/3/issue/{issueKey}/worklog`
   - Parameters: `startAt`, `maxResults`
   - Returns paginated worklog entries

3. **Get Current User**: Already implemented via GraphQL
   - Returns user account information

## Error Handling

### Error Categories

1. **Authentication Errors** (401, 403)
   - Display: "Authentication failed. Please check your API credentials in ~/.jirarc"
   - Action: Disable report generation until credentials are fixed

2. **Network Errors**
   - Display: "Network error: Unable to connect to Jira. Please check your internet connection."
   - Action: Allow retry with refresh button

3. **Rate Limit Errors** (429)
   - Display: "Jira API rate limit exceeded. Please wait a few minutes and try again."
   - Action: Disable refresh for 60 seconds

4. **Invalid Date Range**
   - Display: "Invalid date range: Start date must be before or equal to end date."
   - Action: Highlight invalid fields, prevent report generation

5. **Partial Data Errors**
   - Display: "Warning: Some data could not be retrieved. Showing partial results."
   - Action: Display available data with warning indicator

6. **Export Errors**
   - Display: "Failed to export report: [specific error message]"
   - Action: Allow retry, suggest checking file permissions

### Error Handling Strategy

- All API calls will be wrapped in error handling with appropriate user feedback
- Errors will be logged to the application log for debugging
- UI will remain responsive during errors (no blocking operations)
- Partial failures will show available data rather than failing completely
- Network operations will have reasonable timeouts (30 seconds)

## Testing Strategy

### Unit Testing

The implementation will include unit tests for:

1. **Date Range Validation**
   - Test valid date ranges
   - Test invalid date ranges (start > end)
   - Test default date range generation

2. **Issue Filtering Logic**
   - Test `IsIssueOpenDuringRange` with various issue states
   - Test issues created before, during, and after range
   - Test issues closed before, during, and after range

3. **Data Aggregation**
   - Test rollup calculations with sample data
   - Test status percentage calculations
   - Test time duration summation

4. **Export Formatting**
   - Test CSV generation with sample data
   - Test header generation with various inputs
   - Test data formatting for export

### Property-Based Testing

The implementation will use property-based testing with the `gopter` library (Go property testing library) to verify correctness properties. Each property-based test will run a minimum of 100 iterations to ensure thorough coverage.

Property-based tests will be tagged with comments explicitly referencing the correctness property from this design document using the format: `**Feature: jira-reporting-window, Property {number}: {property_text}**`

Each correctness property defined in the next section will be implemented by a single property-based test.

### Integration Testing

Integration tests will verify:
- End-to-end report generation with mock Jira API responses
- Cache invalidation behavior
- Window lifecycle (open, close, reopen)

### Manual Testing

Manual testing will cover:
- UI responsiveness and layout
- Export file generation and formatting
- Error message clarity and helpfulness
- Loading indicator behavior

## Performance Considerations

### Caching Strategy

- Cache report data for 5 minutes to reduce API calls
- Invalidate cache when date range changes
- Clear cache when window closes
- Use read-write mutex for thread-safe cache access

### Pagination

- Fetch issues in batches of 100 (Jira API default)
- Implement pagination for large result sets
- Show progress indicator for multi-page fetches

### Concurrent Operations

- Fetch different report types concurrently when possible
- Use goroutines for API calls to keep UI responsive
- Implement proper synchronization for shared data structures

### Memory Management

- Limit maximum issues per report to 1000
- Stream export data to files rather than building in memory
- Release cached data when window closes


## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Window independence preservation
*For any* main window state, closing the reporting window should leave the main window in the same state it was in before the reporting window was opened.
**Validates: Requirements 1.3**

### Property 2: Window singleton behavior
*For any* number of attempts to open the reporting window, only one reporting window instance should exist and be focused.
**Validates: Requirements 1.4**

### Property 3: Date range validation
*For any* start date and end date pair, the system should reject the date range if and only if the start date is after the end date.
**Validates: Requirements 2.2, 2.3**

### Property 4: Date range change triggers refresh
*For any* date range modification, all displayed reports should reflect the new date range after the change.
**Validates: Requirements 2.4**

### Property 5: Assigned issues query correctness
*For any* user, the assigned issues query should return only issues where that user is the assignee.
**Validates: Requirements 3.1**

### Property 6: Assigned issues display completeness
*For any* issue in the assigned issues report, the display should include the issue key, summary, status, and assignee.
**Validates: Requirements 3.2**

### Property 7: Assigned issues date filtering
*For any* date range and set of issues, the assigned issues report should include only issues updated within that date range.
**Validates: Requirements 3.4**

### Property 8: Assigned issues sort order
*For any* list of assigned issues, the issues should be sorted by issue key in ascending lexicographic order.
**Validates: Requirements 3.5**

### Property 9: Worked issues query correctness
*For any* user, the worked issues query should return only issues where that user has worklog entries.
**Validates: Requirements 4.1**

### Property 10: Worked issues display completeness
*For any* issue in the worked issues report, the display should include the issue key, summary, status, and total time logged.
**Validates: Requirements 4.2**

### Property 11: Worked issues date filtering
*For any* date range and set of issues, the worked issues report should include only issues with worklogs created within that date range.
**Validates: Requirements 4.4**

### Property 12: Time aggregation correctness
*For any* set of worklog entries for a user within a date range, the total time logged should equal the sum of all individual worklog durations.
**Validates: Requirements 4.5**

### Property 13: Untracked issues set difference
*For any* set of assigned issues and worked issues, the untracked issues should be exactly those issues that are assigned but not worked.
**Validates: Requirements 5.1**

### Property 14: Untracked issues display completeness
*For any* issue in the untracked issues report, the display should include the issue key, summary, status, and assigned date.
**Validates: Requirements 5.2**

### Property 15: Issue open during range logic
*For any* issue and date range, the issue should be included in reports if and only if the issue was created on or before the range end date AND (the issue is still open OR the issue was closed on or after the range start date).
**Validates: Requirements 5.4, 5.5, 5.6**

### Property 16: Untracked with external worklogs
*For any* issue with worklogs outside the date range but none within it, if the issue was open during the date range and is assigned to the user, it should appear in the untracked issues list.
**Validates: Requirements 5.7**

### Property 17: Rollup assigned count consistency
*For any* date range, the total assigned issues count in the rollup report should equal the number of issues in the assigned issues report.
**Validates: Requirements 6.1**

### Property 18: Rollup worked count consistency
*For any* date range, the total worked issues count in the rollup report should equal the number of issues in the worked issues report.
**Validates: Requirements 6.2**

### Property 19: Status count aggregation
*For any* set of issues grouped by status, the sum of all status counts should equal the total number of unique issues.
**Validates: Requirements 6.3**

### Property 20: Status percentage calculation
*For any* status count and total unique issues count, the percentage should equal (status count / total unique issues) × 100.
**Validates: Requirements 6.4**

### Property 21: Unique issue denominator
*For any* set of issues across all report categories, the denominator for percentage calculations should be the count of unique issue keys.
**Validates: Requirements 6.5**

### Property 22: API error message display
*For any* Jira API call failure, an error message should be displayed to the user.
**Validates: Requirements 7.1**

### Property 23: Partial data graceful degradation
*For any* partial data retrieval scenario, the system should display the available data along with a warning about incomplete results.
**Validates: Requirements 7.5**

### Property 24: CSV export completeness
*For any* report data, the generated CSV file should contain all visible report data with proper column headers.
**Validates: Requirements 8.2**

### Property 25: PDF export completeness
*For any* report data, the generated PDF file should contain all visible report data formatted in tables.
**Validates: Requirements 8.3**

### Property 26: Export header completeness
*For any* export operation, the header should contain the user name, start date, end date, and export generation timestamp.
**Validates: Requirements 8.4, 8.5, 8.6, 8.7**

### Property 27: ISO 8601 date formatting
*For any* date or timestamp in an export header, it should be formatted according to ISO 8601 standard.
**Validates: Requirements 8.6, 8.7**

### Property 28: CSV section separation
*For any* multi-section CSV export, each report section should be separated from the next by blank rows.
**Validates: Requirements 8.8**

### Property 29: PDF section formatting
*For any* multi-section PDF export, each section should have a heading and appropriate page breaks.
**Validates: Requirements 8.9**

### Property 30: Export confirmation display
*For any* successful export operation, a confirmation message with the file location should be displayed.
**Validates: Requirements 8.10**

### Property 31: Export filename uniqueness
*For any* export file creation, the filename should include a timestamp to prevent overwriting existing files.
**Validates: Requirements 8.11**

### Property 32: Loading indicator visibility during generation
*For any* report generation operation in progress, a loading indicator should be visible in the reporting window.
**Validates: Requirements 9.1**

### Property 33: Control disabling during loading
*For any* data retrieval operation in progress, report controls should be disabled to prevent duplicate requests.
**Validates: Requirements 9.2**

### Property 34: Success state transition
*For any* successful data retrieval, the loading indicator should be hidden and the results should be displayed.
**Validates: Requirements 9.3**

### Property 35: Failure state transition
*For any* failed data retrieval, the loading indicator should be hidden and an error message should be displayed.
**Validates: Requirements 9.4**

### Property 36: Request queue ordering
*For any* sequence of report requests, they should be processed in the order they were received.
**Validates: Requirements 9.5**

### Property 37: Cache storage on success
*For any* successful data retrieval, the results should be stored in the cache with the current timestamp.
**Validates: Requirements 10.1**

### Property 38: Cache hit behavior
*For any* report view request with a date range matching cached data, the cached data should be used instead of making new API calls.
**Validates: Requirements 10.2**

### Property 39: Cache invalidation on date change
*For any* date range modification, the cache should be invalidated and new data should be fetched.
**Validates: Requirements 10.3**

### Property 40: Cache cleanup on window close
*For any* reporting window close event, all cached report data should be cleared from memory.
**Validates: Requirements 10.4**

### Property 41: Cache expiration behavior
*For any* cached data with a timestamp older than 5 minutes, the data should be refreshed on the next report view request.
**Validates: Requirements 10.5**
