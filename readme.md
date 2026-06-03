# JiraTimeWidget

A desktop time-tracking application for Jira Cloud built with Go and the Fyne GUI framework. Track time on assigned issues, log worklogs directly to Jira, manage issue status transitions, and generate reports.

## Running

```
go run .
```

CLI mode (view today's logged time):

```
go run . logs
```

## Configuration

On first launch, a setup wizard walks you through connecting to your Jira instance. It prompts for your Jira URL, email, and API token, then writes `~/.jirarc` automatically.

To generate an API token, visit: https://id.atlassian.com/manage-profile/security/api-tokens

The resulting config file (`~/.jirarc`):

```json
{
  "jira": "<your-jira-api-token>",
  "email": "<your-atlassian-email>",
  "graphqlUri": "https://your-org.atlassian.net/gateway/api/graphql",
  "cloudId": "<your-cloud-id>"
}
```

Authentication uses Basic Auth (email + API token) for Jira Cloud, or Bearer token if no email is provided. If the config file is deleted, the setup wizard will reappear on next launch.

## Features

### Setup Wizard

On first run (when `~/.jirarc` is missing), a guided setup wizard collects the Jira URL, email, and API token. Validates inputs, writes the config with `0600` permissions, and transitions directly to the main app.

**Key functions:**
- `jiraConfigExists()` — Checks if `~/.jirarc` exists and has content
- `saveJiraConfig()` — Writes config as formatted JSON
- `showSetupWizard()` — Renders the wizard UI and handles validation/save

### Time Tracking

Track time on Jira issues with start/stop controls or manual duration entry (e.g., `1h 30m`, `45m`, `2h`). Worklogs are submitted to the Jira REST API with ADF-formatted comments.

**Key functions:**
- `createTimeButtons()` — Renders Start, Stop, Reset, and Browser buttons
- `createIssueSelector()` — Dropdown populated with recent issues, refresh capability
- `logWorkToJira()` — Submits worklog via `jiraApiFunctions.AddWorklog()`
- `parseDuration()` — Parses human-readable duration strings
- `formatDurationForJira()` — Formats `time.Duration` into Jira-compatible strings

### Recent Issues

Fetches up to 20 recently updated issues assigned to the current user. Excludes issues in Deferred or On Hold status. Issues closed more than 7 days ago are also excluded.

**Key functions:**
- `getRecentIssues()` — Queries Jira via `jiraApiFunctions.MakeJiraAPICall()` with JQL filtering

### Issue Status Management

Displays the current status of the selected issue and allows transitioning to available statuses. Transitions display directional icons (→ forward, ← backward) based on status category progression.

**Key functions:**
- `GetIssueStatus()` — Fetches status via `jiraApiFunctions.GetIssue()`
- `GetAvailableTransitions()` — Fetches transitions via `jiraApiFunctions.GetIssueTransitions()`
- `ExecuteStatusTransition()` — Executes via `jiraApiFunctions.TransitionIssue()`
- `determineTransitionDirection()` — Categorizes transitions as forward/backward using status category ordering
- `showTransitionsMenu()` — Popup menu UI for transition selection

### Open in Browser

Opens the selected Jira issue in the system's default browser.

**Key functions:**
- `openBrowser()` — Cross-platform browser launch (Windows, macOS, Linux)

### Local Time Log

All logged worklogs are persisted locally to `~/.jira_time_log.json` for offline reference.

**Key functions:**
- `saveTimeLogEntry()` — Appends entry to JSON log file
- `getTodaysTimeLog()` — Reads and filters entries for the current date
- `viewLogs()` — CLI output of today's entries with totals

### Reporting Window

A dedicated window (singleton pattern) for generating reports over a configurable date range. Defaults to the current month. Includes a date picker (calendar widget), loading indicators, and tabbed report views.

**Report types:**

| Tab | Description | Service Function |
|-----|-------------|-----------------|
| Assigned Issues | Issues assigned to user updated within the date range | `ReportService.GetAssignedIssues()` |
| Worked Issues | Issues with worklogs by user in the date range, with time aggregation | `ReportService.GetWorkedIssues()` |
| Untracked Issues | Assigned issues that were open during the range but have no worklogs | `ReportService.GetUntrackedIssues()` |
| Rollup Report | Summary statistics: total assigned, total worked, unique issues, status breakdown with percentages | `ReportService.GetRollupReport()` |

**Supporting infrastructure:**
- `ReportCache` — Thread-safe (RWMutex) cache with 5-minute TTL and date range validation
- `RequestQueue` — Mutex-based sequential processing of report generation requests
- `CategorizeError()` — Classifies API errors into Auth, Network, RateLimit, PartialData, InvalidDateRange with user-friendly messages
- `IsIssueOpenDuringRange()` — Determines if an issue was open during a date range based on created/resolved dates

### Report Export

Export generated reports to CSV or PDF via a file save dialog.

**Key functions:**
- `ExportService.ExportToCSVFile()` — Writes all report sections to a structured CSV
- `ExportService.ExportToPDFFile()` — Generates a formatted PDF with tables using `gofpdf`

### Current User

Fetches the authenticated user's profile via the Jira GraphQL API to identify the user for report filtering.

**Key functions:**
- `getCurrentUser()` — GraphQL query against `me { user { ... } }`
- `loadJiraConfig()` — Reads `~/.jirarc` and sets API credentials

## Project Structure

| File | Responsibility |
|------|---------------|
| `main.go` | Entry point, CLI mode, Fyne app initialization, config check |
| `setup_wizard.go` | First-run setup wizard UI and config file creation |
| `uiUtils.go` | Main form layout, time buttons, issue selector, log button |
| `jiraUtils.go` | Jira REST API interactions (get issue, log work, status, transitions) |
| `userUtils.go` | Config loading, current user GraphQL query |
| `recentIssues.go` | Recent issues JQL query and parsing |
| `timeLog.go` | Local time log persistence and CLI display |
| `queries.go` | GraphQL query constants |
| `types.go` | All shared types and data structures |
| `reporting.go` | DateRange, report types, ReportCache |
| `report_service.go` | Report data fetching and generation logic |
| `report_views.go` | Fyne table views for each report type |
| `report_errors.go` | Error categorization and user-friendly messages |
| `reporting_window.go` | Reporting window UI, date picker, request processing |
| `request_queue.go` | Thread-safe request queue implementation |
| `export_service.go` | CSV and PDF export logic |
| `jiraApiFunctions/` | Jira REST API wrapper (search, issues, worklogs, transitions, etc.) |

## Dependencies

| Library | Purpose |
|---------|---------|
| `fyne.io/fyne/v2` | Desktop GUI framework |
| `fyne.io/x/fyne` | Extended widgets (Calendar date picker) |
| `github.com/jung-kurt/gofpdf` | PDF generation |
| `github.com/stretchr/testify` | Test assertions |
