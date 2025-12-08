# Requirements Document

## Introduction

This document specifies the requirements for a Jira reporting feature that allows users to generate and view various reports about their Jira issues within a date range. The feature will provide insights into assigned issues, work log activity, and status distributions through a dedicated reporting window in the JiraWidgetLite application.

## Glossary

- **JiraWidgetLite**: The Go-based desktop application that provides Jira time tracking functionality
- **Reporting Window**: A separate UI window that displays Jira reports and allows date range selection
- **Work Log Entry**: A time tracking record in Jira that indicates a user has worked on an issue
- **Assigned Issue**: A Jira issue where the current user is set as the assignee
- **Status**: The current state of a Jira issue (e.g., To Do, In Progress, Done)
- **Date Range**: A period defined by a start date and end date for filtering report data
- **Rollup Report**: A summary report that aggregates multiple metrics into counts and percentages
- **JQL**: Jira Query Language used to search and filter issues
- **REST API**: The Jira REST API v3 used for retrieving issue and worklog data

## Requirements

### Requirement 1

**User Story:** As a user, I want to open a separate reporting window from the main application, so that I can view reports without disrupting my time tracking workflow.

#### Acceptance Criteria

1. WHEN the user clicks a reports button in the main window THEN the system SHALL open a new reporting window
2. WHEN the reporting window is opened THEN the system SHALL display it as a separate window that can be positioned independently
3. WHEN the reporting window is closed THEN the system SHALL maintain the main application window in its current state
4. WHEN the reporting window is already open THEN the system SHALL bring the existing window to focus rather than creating a duplicate

### Requirement 2

**User Story:** As a user, I want to select a date range for my reports, so that I can analyze my work activity for specific time periods.

#### Acceptance Criteria

1. WHEN the reporting window opens THEN the system SHALL display date range selection controls with start date and end date inputs
2. WHEN the user selects a start date THEN the system SHALL validate that the start date is not after the end date
3. WHEN the user selects an end date THEN the system SHALL validate that the end date is not before the start date
4. WHEN the user modifies the date range THEN the system SHALL update all displayed reports to reflect the new date range
5. WHEN the reporting window opens THEN the system SHALL default the date range to the current month

### Requirement 3

**User Story:** As a user, I want to view a list of all issues assigned to me within a date range, so that I can see my current and past responsibilities.

#### Acceptance Criteria

1. WHEN the user requests the assigned issues report THEN the system SHALL query Jira for all issues where the current user is the assignee
2. WHEN displaying assigned issues THEN the system SHALL show the issue key, summary, status, and assignee for each issue
3. WHEN the assigned issues list is empty THEN the system SHALL display a message indicating no assigned issues were found
4. WHEN the system retrieves assigned issues THEN the system SHALL filter results to include only issues updated within the selected date range
5. WHEN displaying the assigned issues list THEN the system SHALL sort issues by issue key in ascending order

### Requirement 4

**User Story:** As a user, I want to view a list of all issues I have worked on within a date range, so that I can track which tasks I have actively contributed to.

#### Acceptance Criteria

1. WHEN the user requests the worked issues report THEN the system SHALL query Jira for all issues where the current user has work log entries
2. WHEN displaying worked issues THEN the system SHALL show the issue key, summary, status, and total time logged for each issue
3. WHEN the worked issues list is empty THEN the system SHALL display a message indicating no worked issues were found
4. WHEN the system retrieves worked issues THEN the system SHALL filter results to include only issues with work logs created within the selected date range
5. WHEN calculating total time logged THEN the system SHALL sum all work log entries for the current user within the date range

### Requirement 5

**User Story:** As a user, I want to view a list of assigned issues that do not have work logs, so that I can identify tasks I may have forgotten to track time for.

#### Acceptance Criteria

1. WHEN the user requests the untracked issues report THEN the system SHALL identify all assigned issues that have no work log entries from the current user
2. WHEN displaying untracked issues THEN the system SHALL show the issue key, summary, status, and assigned date for each issue
3. WHEN the untracked issues list is empty THEN the system SHALL display a message indicating all assigned issues have work logs
4. WHEN the system identifies untracked issues THEN the system SHALL filter to include only issues that were open during the selected date range
5. WHEN an issue was created after the date range end date THEN the system SHALL exclude that issue from the untracked list
6. WHEN an issue was closed before the date range start date THEN the system SHALL exclude that issue from the untracked list
7. WHEN an issue has work logs outside the date range but none within it THEN the system SHALL include that issue in the untracked list if the issue was open during the date range

### Requirement 6

**User Story:** As a user, I want to view a rollup report with summary statistics, so that I can quickly understand my work distribution and progress.

#### Acceptance Criteria

1. WHEN the user requests the rollup report THEN the system SHALL display the total number of assigned issues within the date range
2. WHEN the rollup report is displayed THEN the system SHALL show the total number of issues worked on within the date range
3. WHEN the rollup report is displayed THEN the system SHALL show the count of issues grouped by status
4. WHEN the rollup report is displayed THEN the system SHALL calculate and display the percentage of total issues for each status
5. WHEN calculating status percentages THEN the system SHALL use the total number of unique issues across all report categories as the denominator

### Requirement 7

**User Story:** As a user, I want the reporting window to handle API errors gracefully, so that I can understand when data cannot be retrieved and why.

#### Acceptance Criteria

1. WHEN a Jira API call fails THEN the system SHALL display an error message indicating the failure
2. WHEN authentication fails THEN the system SHALL display a message prompting the user to check their API credentials
3. WHEN a network error occurs THEN the system SHALL display a message indicating connectivity issues
4. WHEN an API rate limit is exceeded THEN the system SHALL display a message indicating the rate limit and suggest retrying later
5. WHEN partial data is retrieved successfully THEN the system SHALL display the available data with a warning about incomplete results

### Requirement 8

**User Story:** As a user, I want to export report data to CSV or PDF format, so that I can share or archive my work activity information in different formats.

#### Acceptance Criteria

1. WHEN the user clicks an export button THEN the system SHALL display format options for CSV and PDF
2. WHEN the user selects CSV export THEN the system SHALL generate a CSV file containing all visible report data with proper column headers
3. WHEN the user selects PDF export THEN the system SHALL generate a PDF file containing all visible report data with formatted tables
4. WHEN exporting report data THEN the system SHALL include a header section containing the user name, date range criteria, and the export generation date
5. WHEN the export header is generated THEN the system SHALL display the current user's name or account identifier
6. WHEN the export header is generated THEN the system SHALL display the start date and end date of the report in ISO 8601 format
7. WHEN the export header is generated THEN the system SHALL display the export generation timestamp in ISO 8601 format
8. WHEN exporting to CSV THEN the system SHALL format each report section as a separate table with blank rows as separators
9. WHEN exporting to PDF THEN the system SHALL format the data with proper page breaks and section headings
10. WHEN the export is complete THEN the system SHALL display a confirmation message with the file location
11. WHEN the export file is created THEN the system SHALL name the file with a timestamp to prevent overwrites

### Requirement 9

**User Story:** As a user, I want the reporting window to display loading indicators, so that I understand when data is being fetched from Jira.

#### Acceptance Criteria

1. WHEN a report is being generated THEN the system SHALL display a loading indicator in the reporting window
2. WHEN data retrieval is in progress THEN the system SHALL disable report controls to prevent duplicate requests
3. WHEN data retrieval completes successfully THEN the system SHALL hide the loading indicator and display the results
4. WHEN data retrieval fails THEN the system SHALL hide the loading indicator and display an error message
5. WHEN multiple reports are requested sequentially THEN the system SHALL queue requests and process them in order

### Requirement 10

**User Story:** As a user, I want the reporting window to cache report data temporarily, so that I can switch between different report views without re-fetching data unnecessarily.

#### Acceptance Criteria

1. WHEN report data is successfully retrieved THEN the system SHALL cache the results in memory
2. WHEN the user switches between report views with the same date range THEN the system SHALL use cached data instead of making new API calls
3. WHEN the user changes the date range THEN the system SHALL invalidate the cache and fetch new data
4. WHEN the reporting window is closed THEN the system SHALL clear all cached report data
5. WHEN cached data is older than 5 minutes THEN the system SHALL automatically refresh the data on the next report view request
