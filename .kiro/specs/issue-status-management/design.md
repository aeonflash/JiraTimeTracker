# Design Document

## Overview

This design document outlines the implementation approach for adding issue status viewing and updating capabilities to JiraWidgetLite. The feature will integrate seamlessly with the existing time tracking interface, allowing users to see and change issue statuses without leaving the application.

The implementation will leverage the existing Jira REST API integration and Fyne UI framework, adding new UI components and API interactions to support status management.

## Architecture

### High-Level Components

1. **Status Display Component**: A read-only label showing the current status of the selected issue
2. **Status Change Component**: An interactive control (button/menu) that allows users to select and apply status transitions
3. **API Integration Layer**: Functions to retrieve available transitions and execute status changes via Jira REST API
4. **UI State Management**: Logic to update the UI based on status changes and API responses

### Component Interaction Flow

```
User selects issue → Fetch issue details (including status) → Display status
User clicks status change → Fetch available transitions → Display transition options
User selects transition → Execute transition via API → Update displayed status → Show feedback
```

## Components and Interfaces

### 1. UI Components

#### Status Display Label
- **Location**: Positioned near the issue selector dropdown, integrated into the existing status label area
- **Behavior**: Updates automatically when an issue is selected or after a successful status transition
- **Styling**: Uses consistent styling with existing UI elements, with visual distinction (e.g., badge or colored text)

#### Status Change Button
- **Type**: Button that opens a popup menu or dialog with available transitions
- **Location**: Adjacent to the status display label
- **Behavior**: 
  - Disabled when no issue is selected
  - Fetches transitions on click
  - Displays loading state while fetching
  - Shows available transitions with directional icons

#### Transition Selection Menu
- **Type**: Popup menu or dialog
- **Content**: List of available status transitions with:
  - Transition name (target status)
  - Directional icon (→ for forward, ← for backward)
  - Visual grouping if needed
- **Behavior**: Executes selected transition and closes on selection

### 2. Data Structures

```go
// StatusInfo represents the current status of an issue
type StatusInfo struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Category    string `json:"statusCategory"` // e.g., "To Do", "In Progress", "Done"
}

// Transition represents an available status transition
type Transition struct {
    ID          string     `json:"id"`
    Name        string     `json:"name"`
    To          StatusInfo `json:"to"`
    IsForward   bool       // Derived field for icon display
}

// TransitionsResponse represents the API response for available transitions
type TransitionsResponse struct {
    Transitions []Transition `json:"transitions"`
}
```

### 3. API Integration Functions

#### GetIssueStatus
```go
func GetIssueStatus(issueKey string) (*StatusInfo, error)
```
- Calls existing `GetIssue` function with fields parameter set to "status"
- Parses response to extract status information
- Returns StatusInfo struct or error

#### GetAvailableTransitions
```go
func GetAvailableTransitions(issueKey string) ([]Transition, error)
```
- Calls existing `GetIssueTransitions` function from jiraIssueFunctions.go
- Parses response to extract available transitions
- Determines forward/backward direction based on status category
- Returns slice of Transition structs or error

#### ExecuteStatusTransition
```go
func ExecuteStatusTransition(issueKey string, transitionID string) error
```
- Calls existing `TransitionIssue` function with transition data
- Request body format: `{"transition": {"id": "transitionID"}}`
- Returns error if transition fails

### 4. UI State Management

#### Extended UIComponents Struct
```go
type UIComponents struct {
    // ... existing fields ...
    CurrentStatus      *StatusInfo
    StatusDisplayLabel *widget.Label
    StatusChangeButton *widget.Button
    TransitionsMenu    *widget.PopUpMenu
}
```

#### Status Update Flow
1. When issue is selected, call `GetIssueStatus` and update `StatusDisplayLabel`
2. Store current status in `CurrentStatus` field
3. Enable `StatusChangeButton`
4. On button click, fetch transitions and display menu
5. On transition selection, execute transition and refresh status display

## Data Models

### Status Categories and Direction Logic

Jira status categories follow a general workflow progression:
- **To Do** (initial state)
- **In Progress** (active work)
- **Done** (completed)

**Forward transitions**: Moving from To Do → In Progress → Done
**Backward transitions**: Moving from Done → In Progress → To Do

Direction determination logic:
```go
func determineTransitionDirection(currentCategory, targetCategory string) bool {
    categoryOrder := map[string]int{
        "To Do": 1,
        "In Progress": 2,
        "Done": 3,
    }
    
    currentOrder := categoryOrder[currentCategory]
    targetOrder := categoryOrder[targetCategory]
    
    return targetOrder > currentOrder // true = forward, false = backward
}
```

### Icon Selection

- Forward transitions: Use right arrow icon (→ or "▶")
- Backward transitions: Use left arrow icon (← or "◀")
- Icons will be prepended to transition names in the menu

## Error Handling

### API Error Scenarios

1. **Failed to fetch status**
   - Display error message in status label
   - Keep status change button disabled
   - Log error details

2. **Failed to fetch transitions**
   - Show error dialog to user
   - Keep current status displayed
   - Log error details

3. **Failed to execute transition**
   - Show error dialog with failure reason
   - Keep current status displayed (don't update)
   - Log error details

4. **Network timeout**
   - Show timeout message
   - Provide retry option
   - Log timeout event

### Error Message Format

- User-facing: Clear, actionable messages (e.g., "Failed to update status. Please try again.")
- Logs: Detailed technical information including API response codes and messages

### Graceful Degradation

- If status fetching fails, show "Status unavailable" but allow other operations
- If transitions fetching fails, disable status change button but keep display functional
- Maintain UI responsiveness during API calls with loading indicators

## Testing Strategy

### Unit Tests

1. **Status Parsing Tests**
   - Test parsing of various status response formats
   - Test handling of missing or malformed status data
   - Test status category extraction

2. **Transition Direction Tests**
   - Test direction determination for all category combinations
   - Test edge cases (same category transitions)
   - Test handling of unknown categories

3. **API Integration Tests**
   - Mock API responses for status fetching
   - Mock API responses for transitions fetching
   - Mock API responses for transition execution
   - Test error response handling

### Integration Tests

1. **UI Integration Tests**
   - Test status display updates when issue is selected
   - Test status change button enable/disable logic
   - Test transitions menu population
   - Test status update after successful transition

2. **End-to-End Flow Tests**
   - Test complete flow: select issue → view status → change status → verify update
   - Test error recovery flows
   - Test concurrent operations (e.g., status change while timer is running)

### Manual Testing Checklist

- [ ] Status displays correctly for various issue types
- [ ] Transitions menu shows correct options for different workflows
- [ ] Forward/backward icons display correctly
- [ ] Status updates immediately after successful transition
- [ ] Error messages display appropriately for various failure scenarios
- [ ] UI remains responsive during API calls
- [ ] Status feature integrates seamlessly with time tracking
- [ ] Window resizing accommodates new status components
- [ ] Status persists correctly when switching between issues

## UI Layout Integration

### Positioning Strategy

The status components will be integrated into the existing UI with minimal disruption:

1. **Status Display**: Add to the issue selector row, showing as "[Status: In Progress]" next to the issue key
2. **Status Change Button**: Small button with "⚙️" icon positioned next to the status display
3. **Compact Layout**: Maintain the existing compact window size (550x200 base)

### Layout Mockup

```
┌─────────────────────────────────────────────────────┐
│ [Issue Dropdown ▼] [Status: In Progress] [⚙️] [🔄] │
│                                                     │
│ [Start] [Stop] [Reset] [🌐]                        │
│                                                     │
│ Duration: [1h 30m                              ]   │
│                                                     │
│ Work Comment:                                       │
│ [                                              ]   │
│ [                                              ]   │
│                                                     │
│ ✅ Ready to track time on PROJ-123                 │
│                                          [Log Time] │
└─────────────────────────────────────────────────────┘
```

### Responsive Behavior

- Status components will be part of the dynamic resizing logic
- When status change menu is open, it will overlay the main window
- Status display will truncate long status names with ellipsis if needed
