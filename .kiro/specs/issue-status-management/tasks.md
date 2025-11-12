# Implementation Plan

- [x] 1. Create data structures for status and transitions
  - Define StatusInfo struct to represent issue status with ID, Name, Description, and Category fields
  - Define Transition struct to represent available transitions with ID, Name, To (StatusInfo), and IsForward fields
  - Define TransitionsResponse struct for parsing API responses
  - Add these types to types.go file
  - _Requirements: 1.1, 1.2, 2.2, 2.3_

- [x] 2. Implement API integration functions for status management
  - [x] 2.1 Create GetIssueStatus function
    - Implement function that calls existing GetIssue API with fields parameter set to "status"
    - Parse JSON response to extract status information
    - Return StatusInfo struct or error
    - Handle API errors and malformed responses
    - _Requirements: 1.2, 1.4_
  
  - [x] 2.2 Create GetAvailableTransitions function
    - Implement function that calls existing GetIssueTransitions API
    - Parse JSON response to extract available transitions
    - Implement logic to determine forward/backward direction based on status categories
    - Return slice of Transition structs with IsForward field populated
    - Handle API errors and empty transition lists
    - _Requirements: 2.2, 2.3, 2.6, 2.7_
  
  - [x] 2.3 Create ExecuteStatusTransition function
    - Implement function that calls existing TransitionIssue API with transition ID
    - Format request body as {"transition": {"id": "transitionID"}}
    - Return error if transition fails
    - Handle API errors with descriptive messages
    - _Requirements: 2.4, 2.5, 3.2_

- [x] 3. Extend UI components structure
  - Add CurrentStatus field (*StatusInfo) to UIComponents struct
  - Add StatusDisplayLabel field (*widget.Label) to UIComponents struct
  - Add StatusChangeButton field (*widget.Button) to UIComponents struct
  - Update UIComponents struct in uiUtils.go
  - _Requirements: 1.3, 2.1, 4.1, 4.2_

- [x] 4. Implement status display UI component
  - [x] 4.1 Create status display label
    - Add StatusDisplayLabel to UI initialization
    - Position label in the issue selector row near the issue dropdown
    - Format status text as "[Status: {status_name}]"
    - Apply consistent styling with existing UI elements
    - Initially hide until issue is selected
    - _Requirements: 1.1, 1.3, 4.1, 4.2, 4.5_
  
  - [x] 4.2 Integrate status fetching into issue selection flow
    - Modify fetchIssue function in createIssueSelector to call GetIssueStatus
    - Update StatusDisplayLabel text when status is retrieved
    - Show StatusDisplayLabel when issue is selected
    - Handle status fetch errors by displaying error in status label
    - Ensure status updates within 2 seconds of issue selection
    - _Requirements: 1.1, 1.2, 1.4, 1.5_

- [x] 5. Implement status change UI component
  - [x] 5.1 Create status change button
    - Add StatusChangeButton with gear icon (⚙️) next to status display
    - Initially disable button until issue is selected
    - Enable button when issue is selected and status is loaded
    - Add click handler to fetch and display transitions
    - _Requirements: 2.1, 4.2, 4.5_
  
  - [x] 5.2 Create transitions menu with directional icons
    - Implement popup menu that displays available transitions
    - Fetch transitions using GetAvailableTransitions when button is clicked
    - Display loading indicator while fetching transitions
    - Format menu items with directional icons: "→ {transition_name}" for forward, "← {transition_name}" for backward
    - Handle empty transitions list with appropriate message
    - Handle API errors with error dialog
    - _Requirements: 2.2, 2.3, 2.6, 2.7, 3.5_
  
  - [x] 5.3 Implement transition execution and status update
    - Add click handler for each menu item to execute selected transition
    - Call ExecuteStatusTransition with selected transition ID
    - Show loading state during transition execution
    - On success, refresh status display by calling GetIssueStatus
    - Display success message in status label for 3+ seconds
    - On failure, display error message with failure details
    - Clear previous feedback messages before showing new ones
    - Close transitions menu after execution
    - _Requirements: 2.4, 2.5, 3.1, 3.2, 3.3, 3.4, 3.5_

- [x] 6. Update window resizing logic
  - Modify resizeWindowToContent function to account for status display components
  - Ensure status components are included in height calculations
  - Test that window resizes appropriately when status components are shown/hidden
  - Verify compact window size is maintained (550x200 base)
  - _Requirements: 4.3, 4.4_

- [x] 7. Integrate status components into main form
  - Update createIssueSelector to include status display and change button in the selector row
  - Arrange components horizontally: [Issue Dropdown] [Status Display] [Status Change Button] [Refresh Button]
  - Ensure proper spacing and alignment with existing components
  - Test that all components are visible and functional
  - Verify layout matches design mockup
  - _Requirements: 4.1, 4.2, 4.5_
