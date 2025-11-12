# Requirements Document

## Introduction

This feature enables users to view and update the status of Jira issues directly within the JiraWidgetLite application. Users will be able to see the current status of a selected issue and transition it to different statuses based on the available workflow transitions for that issue.

## Glossary

- **JiraWidgetLite**: The desktop application built with Fyne that provides time tracking and issue management capabilities for Jira
- **Issue Status**: The current state of a Jira issue within its workflow (e.g., "To Do", "In Progress", "Done")
- **Status Transition**: The action of moving an issue from one status to another within the allowed workflow
- **Workflow**: The set of statuses and transitions defined for a particular issue type in Jira
- **Issue Key**: The unique identifier for a Jira issue (e.g., "PROJ-123")
- **UI Component**: A visual element in the JiraWidgetLite interface
- **REST API**: The Jira REST API v3 used for retrieving and updating issue data
- **Forward Progress**: A status transition that moves an issue toward completion in the workflow
- **Backward Movement**: A status transition that moves an issue toward an earlier state in the workflow

## Requirements

### Requirement 1

**User Story:** As a user, I want to see the current status of a selected issue, so that I know what state the issue is in without opening Jira in my browser

#### Acceptance Criteria

1. WHEN a user selects an issue from the issue dropdown, THE JiraWidgetLite SHALL display the current status of that issue in the UI
2. THE JiraWidgetLite SHALL retrieve the current status using the Jira REST API v3 issue endpoint
3. THE JiraWidgetLite SHALL display the status label in a clearly visible location within the main window
4. IF the API call to retrieve the issue status fails, THEN THE JiraWidgetLite SHALL display an error message to the user
5. THE JiraWidgetLite SHALL update the displayed status within 2 seconds of issue selection

### Requirement 2

**User Story:** As a user, I want to change the status of an issue to a different valid status, so that I can update my workflow without leaving the application

#### Acceptance Criteria

1. THE JiraWidgetLite SHALL provide a UI control that allows users to initiate a status change for the selected issue
2. WHEN a user initiates a status change, THE JiraWidgetLite SHALL retrieve all available transitions for the selected issue using the Jira REST API v3 transitions endpoint
3. THE JiraWidgetLite SHALL display only the valid status transitions available for the current issue
4. WHEN a user selects a new status from the available transitions, THE JiraWidgetLite SHALL execute the transition using the Jira REST API v3
5. IF the status transition succeeds, THEN THE JiraWidgetLite SHALL update the displayed status to reflect the new status
6. THE JiraWidgetLite SHALL display a right arrow icon next to status transitions that represent forward progress in the workflow
7. THE JiraWidgetLite SHALL display a left arrow icon next to status transitions that represent backward movement in the workflow

### Requirement 3

**User Story:** As a user, I want to receive feedback when a status change succeeds or fails, so that I know whether my action was completed

#### Acceptance Criteria

1. WHEN a status transition completes successfully, THE JiraWidgetLite SHALL display a success message to the user
2. IF a status transition fails, THEN THE JiraWidgetLite SHALL display an error message with details about the failure
3. THE JiraWidgetLite SHALL display feedback messages for a duration of 3 seconds minimum
4. THE JiraWidgetLite SHALL clear any previous feedback messages before displaying new ones
5. THE JiraWidgetLite SHALL maintain the UI in a responsive state during status transition operations

### Requirement 4

**User Story:** As a user, I want the status display to integrate seamlessly with the existing time tracking interface, so that I have a unified experience

#### Acceptance Criteria

1. THE JiraWidgetLite SHALL display the status information in the same window as the time tracking controls
2. THE JiraWidgetLite SHALL position the status controls in a logical location relative to the issue selection dropdown
3. THE JiraWidgetLite SHALL maintain the compact window size while accommodating the status display
4. WHEN the window is resized, THE JiraWidgetLite SHALL adjust the status controls proportionally
5. THE JiraWidgetLite SHALL use consistent styling for status controls that matches the existing UI design
