# Implementation Plan

- [x] 1. Set up core data structures and models
  - Create report data structures (AssignedIssuesReport, WorkedIssuesReport, UntrackedIssuesReport, RollupReport)
  - Create IssueReportItem, WorkedIssueItem, StatusCount types
  - Create DateRange type with validation methods
  - Create ReportCache structure with mutex for thread safety
  - _Requirements: 2.2, 2.3, 3.2, 4.2, 5.2, 6.3, 6.4_

- [x] 1.1 Write property test for date range validation
  - **Property 3: Date range validation**
  - **Validates: Requirements 2.2, 2.3**

- [x] 1.2 Write property test for issue open during range logic
  - **Property 15: Issue open during range logic**
  - **Validates: Requirements 5.4, 5.5, 5.6**

- [x] 2. Implement ReportService with data fetching logic
  - Create ReportService struct with API client integration
  - Implement GetAssignedIssues method with JQL query construction
  - Implement GetWorkedIssues method with worklog filtering
  - Implement GetUntrackedIssues method using set difference logic
  - Implement GetRollupReport method with aggregation logic
  - Add error handling for API failures
  - _Requirements: 3.1, 3.4, 3.5, 4.1, 4.4, 4.5, 5.1, 5.4, 6.1, 6.2, 6.3, 6.4, 6.5, 7.1_

- [x] 2.1 Write property test for assigned issues query correctness
  - **Property 5: Assigned issues query correctness**
  - **Validates: Requirements 3.1**

- [x] 2.2 Write property test for assigned issues date filtering
  - **Property 7: Assigned issues date filtering**
  - **Validates: Requirements 3.4**

- [x] 2.3 Write property test for assigned issues sort order
  - **Property 8: Assigned issues sort order**
  - **Validates: Requirements 3.5**

- [x] 2.4 Write property test for worked issues query correctness
  - **Property 9: Worked issues query correctness**
  - **Validates: Requirements 4.1**

- [x] 2.5 Write property test for worked issues date filtering
  - **Property 11: Worked issues date filtering**
  - **Validates: Requirements 4.4**

- [x] 2.6 Write property test for time aggregation correctness
  - **Property 12: Time aggregation correctness**
  - **Validates: Requirements 4.5**

- [x] 2.7 Write property test for untracked issues set difference
  - **Property 13: Untracked issues set difference**
  - **Validates: Requirements 5.1**

- [x] 2.8 Write property test for untracked with external worklogs
  - **Property 16: Untracked with external worklogs**
  - **Validates: Requirements 5.7**

- [x] 2.9 Write property test for rollup count consistency
  - **Property 17: Rollup assigned count consistency**
  - **Property 18: Rollup worked count consistency**
  - **Validates: Requirements 6.1, 6.2**

- [x] 2.10 Write property test for status aggregation
  - **Property 19: Status count aggregation**
  - **Property 20: Status percentage calculation**
  - **Property 21: Unique issue denominator**
  - **Validates: Requirements 6.3, 6.4, 6.5**

- [x] 3. Implement caching mechanism
  - Create ReportCache with timestamp tracking
  - Implement cache validation logic (5-minute expiration)
  - Implement cache invalidation on date range change
  - Implement cache cleanup on window close
  - Add thread-safe cache access with read-write mutex
  - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5_

- [x] 3.1 Write property test for cache storage on success
  - **Property 37: Cache storage on success**
  - **Validates: Requirements 10.1**

- [x] 3.2 Write property test for cache hit behavior
  - **Property 38: Cache hit behavior**
  - **Validates: Requirements 10.2**

- [x] 3.3 Write property test for cache invalidation
  - **Property 39: Cache invalidation on date change**
  - **Validates: Requirements 10.3**

- [x] 3.4 Write property test for cache cleanup
  - **Property 40: Cache cleanup on window close**
  - **Validates: Requirements 10.4**

- [x] 3.5 Write property test for cache expiration
  - **Property 41: Cache expiration behavior**
  - **Validates: Requirements 10.5**

- [x] 4. Create reporting window UI structure
  - Create ReportingWindow struct with Fyne window
  - Implement NewReportingWindow constructor
  - Add date range picker widgets (start and end date entries)
  - Create tab container for different report views
  - Add loading indicator widget
  - Add status label for messages
  - Add export and refresh buttons
  - Implement window singleton pattern to prevent duplicates
  - _Requirements: 1.1, 1.2, 1.4, 2.1, 2.5, 9.1_

- [x] 4.1 Write property test for window singleton behavior
  - **Property 2: Window singleton behavior**
  - **Validates: Requirements 1.4**

- [x] 5. Implement report view components
  - Create AssignedIssuesView with Fyne table widget
  - Create WorkedIssuesView with Fyne table widget
  - Create UntrackedIssuesView with Fyne table widget
  - Create RollupReportView with statistics display
  - Implement empty state messages for each view
  - Add column headers for each table view
  - _Requirements: 3.2, 3.3, 4.2, 4.3, 5.2, 5.3, 6.1, 6.2, 6.3, 6.4_

- [x] 5.1 Write property test for display completeness
  - **Property 6: Assigned issues display completeness**
  - **Property 10: Worked issues display completeness**
  - **Property 14: Untracked issues display completeness**
  - **Validates: Requirements 3.2, 4.2, 5.2**

- [x] 6. Wire up date range controls and report generation
  - Implement date range validation on user input
  - Connect date range changes to report refresh
  - Implement default date range (current month) on window open
  - Add date range change handler that invalidates cache
  - Connect refresh button to report regeneration
  - _Requirements: 2.2, 2.3, 2.4, 2.5, 10.3_

- [x] 6.1 Write property test for date range change triggers refresh
  - **Property 4: Date range change triggers refresh**
  - **Validates: Requirements 2.4**

- [x] 7. Implement loading states and error handling
  - Add loading indicator show/hide logic
  - Implement control disabling during data fetch
  - Add error message display for API failures
  - Implement specific error messages for auth, network, and rate limit errors
  - Add partial data warning display
  - Implement success/failure state transitions
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 9.1, 9.2, 9.3, 9.4_

- [x] 7.1 Write property test for loading indicator visibility
  - **Property 32: Loading indicator visibility during generation**
  - **Validates: Requirements 9.1**

- [x] 7.2 Write property test for control disabling
  - **Property 33: Control disabling during loading**
  - **Validates: Requirements 9.2**

- [x] 7.3 Write property test for state transitions
  - **Property 34: Success state transition**
  - **Property 35: Failure state transition**
  - **Validates: Requirements 9.3, 9.4**

- [x] 7.4 Write property test for API error message display
  - **Property 22: API error message display**
  - **Validates: Requirements 7.1**

- [x] 7.5 Write property test for partial data graceful degradation
  - **Property 23: Partial data graceful degradation**
  - **Validates: Requirements 7.5**

- [x] 8. Implement request queue for sequential processing
  - Create request queue structure
  - Implement queue processing logic
  - Add request ordering preservation
  - Ensure thread-safe queue operations
  - _Requirements: 9.5_

- [x] 8.1 Write property test for request queue ordering
  - **Property 36: Request queue ordering**
  - **Validates: Requirements 9.5**

- [x] 9. Implement CSV export functionality
  - Create ExportService struct
  - Implement CSV header generation with user, date range, and export date
  - Implement CSV data formatting for each report type
  - Add section separators (blank rows) between report sections
  - Implement file writing with timestamp in filename
  - Add export confirmation message display
  - _Requirements: 8.1, 8.2, 8.4, 8.5, 8.6, 8.7, 8.8, 8.10, 8.11_

- [x] 9.1 Write property test for CSV export completeness
  - **Property 24: CSV export completeness**
  - **Validates: Requirements 8.2**

- [x] 9.2 Write property test for export header completeness
  - **Property 26: Export header completeness**
  - **Validates: Requirements 8.4, 8.5, 8.6, 8.7**

- [x] 9.3 Write property test for ISO 8601 date formatting
  - **Property 27: ISO 8601 date formatting**
  - **Validates: Requirements 8.6, 8.7**

- [x] 9.4 Write property test for CSV section separation
  - **Property 28: CSV section separation**
  - **Validates: Requirements 8.8**

- [x] 9.5 Write property test for export filename uniqueness
  - **Property 31: Export filename uniqueness**
  - **Validates: Requirements 8.11**

- [x] 9.6 Write property test for export confirmation display
  - **Property 30: Export confirmation display**
  - **Validates: Requirements 8.10**

- [x] 10. Implement PDF export functionality
  - Add PDF library dependency (e.g., gofpdf or similar)
  - Implement PDF header generation with user, date range, and export date
  - Implement PDF table formatting for each report type
  - Add section headings and page breaks
  - Implement file writing with timestamp in filename
  - Add export confirmation message display
  - _Requirements: 8.1, 8.3, 8.4, 8.5, 8.6, 8.7, 8.9, 8.10, 8.11_

- [x] 10.1 Write property test for PDF export completeness
  - **Property 25: PDF export completeness**
  - **Validates: Requirements 8.3**

- [x] 10.2 Write property test for PDF section formatting
  - **Property 29: PDF section formatting**
  - **Validates: Requirements 8.9**

- [x] 11. Wire export functionality to UI
  - Add export button click handler
  - Implement format selection dialog (CSV or PDF)
  - Connect export service to current report data
  - Add file save dialog for user to choose location
  - Display export confirmation with file path
  - _Requirements: 8.1, 8.10_

- [x] 12. Integrate reporting window with main application
  - Add "Reports" button to main window UI
  - Implement button click handler to open reporting window
  - Ensure main window state is preserved when reporting window opens/closes
  - Pass current user information to reporting window
  - _Requirements: 1.1, 1.3_

- [x] 12.1 Write property test for window independence preservation
  - **Property 1: Window independence preservation**
  - **Validates: Requirements 1.3**

- [x] 13. Final checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.
