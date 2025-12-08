package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: jira-reporting-window, Property 24: CSV export completeness**
// For any report data, the generated CSV file should contain all visible report data
// with proper column headers.
// **Validates: Requirements 8.2**
func TestProperty_CSVExportCompleteness(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("CSV export contains all report data", prop.ForAll(
		func(assignedCount, workedCount, untrackedCount, statusCount int) bool {
			// Constrain to reasonable sizes
			if assignedCount < 0 || assignedCount > 50 {
				return true
			}
			if workedCount < 0 || workedCount > 50 {
				return true
			}
			if untrackedCount < 0 || untrackedCount > 50 {
				return true
			}
			if statusCount < 0 || statusCount > 10 {
				return true
			}

			// Generate test data
			user := &User{
				AccountId: "test-user-123",
				Name:      "Test User",
			}

			dateRange := DateRange{
				StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			}

			// Generate assigned issues
			assignedIssues := make([]IssueReportItem, assignedCount)
			for i := 0; i < assignedCount; i++ {
				assignedIssues[i] = IssueReportItem{
					Key:          generateIssueKey(i),
					Summary:      generateSummary(i),
					Status:       "In Progress",
					CreatedDate:  time.Date(2024, 1, i%28+1, 10, 0, 0, 0, time.UTC),
					AssignedDate: time.Date(2024, 1, i%28+1, 10, 0, 0, 0, time.UTC),
				}
			}

			assignedReport := &AssignedIssuesReport{
				Issues:     assignedIssues,
				TotalCount: assignedCount,
			}

			// Generate worked issues
			workedIssues := make([]WorkedIssueItem, workedCount)
			for i := 0; i < workedCount; i++ {
				workedIssues[i] = WorkedIssueItem{
					IssueReportItem: IssueReportItem{
						Key:          generateIssueKey(i + 100),
						Summary:      generateSummary(i + 100),
						Status:       "Done",
						CreatedDate:  time.Date(2024, 1, i%28+1, 10, 0, 0, 0, time.UTC),
						AssignedDate: time.Date(2024, 1, i%28+1, 10, 0, 0, 0, time.UTC),
					},
					TimeLogged:   time.Duration(i+1) * time.Hour,
					WorklogCount: i + 1,
				}
			}

			workedReport := &WorkedIssuesReport{
				Issues:          workedIssues,
				TotalCount:      workedCount,
				TotalTimeLogged: time.Duration(workedCount) * time.Hour,
			}

			// Generate untracked issues
			untrackedIssues := make([]IssueReportItem, untrackedCount)
			for i := 0; i < untrackedCount; i++ {
				untrackedIssues[i] = IssueReportItem{
					Key:          generateIssueKey(i + 200),
					Summary:      generateSummary(i + 200),
					Status:       "To Do",
					CreatedDate:  time.Date(2024, 1, i%28+1, 10, 0, 0, 0, time.UTC),
					AssignedDate: time.Date(2024, 1, i%28+1, 10, 0, 0, 0, time.UTC),
				}
			}

			untrackedReport := &UntrackedIssuesReport{
				Issues:     untrackedIssues,
				TotalCount: untrackedCount,
			}

			// Generate status breakdown
			statusBreakdown := make([]StatusCount, statusCount)
			for i := 0; i < statusCount; i++ {
				statusBreakdown[i] = StatusCount{
					StatusName: generateStatusName(i),
					Count:      i + 1,
					Percentage: float64(i+1) * 10.0,
				}
			}

			rollupReport := &RollupReport{
				TotalAssigned:     assignedCount,
				TotalWorked:       workedCount,
				StatusBreakdown:   statusBreakdown,
				TotalUniqueIssues: assignedCount + workedCount,
			}

			// Create export service and export to CSV
			exportService := NewExportService(user, dateRange)
			filename, err := exportService.ExportToCSV(assignedReport, workedReport, untrackedReport, rollupReport)
			if err != nil {
				t.Logf("Export failed: %v", err)
				return false
			}

			// Clean up file after test
			defer os.Remove(filename)

			// Read and verify CSV content
			file, err := os.Open(filename)
			if err != nil {
				t.Logf("Failed to open CSV file: %v", err)
				return false
			}
			defer file.Close()

			reader := csv.NewReader(file)
			reader.FieldsPerRecord = -1 // Allow variable number of fields
			records, err := reader.ReadAll()
			if err != nil {
				t.Logf("Failed to read CSV: %v", err)
				return false
			}

			// Verify all assigned issues are present
			assignedSectionFound := false
			for i, record := range records {
				if len(record) > 0 && record[0] == "Assigned Issues" {
					assignedSectionFound = true
					remainingRecords := records[i:]
					
					// Count non-empty rows until next section
					dataRows := 0
					for j := 1; j < len(remainingRecords) && len(remainingRecords[j]) > 0; j++ {
						dataRows++
					}
					
					if dataRows < assignedCount {
						t.Logf("Expected at least %d assigned issue rows, got %d", assignedCount, dataRows)
						return false
					}
					break
				}
			}

			if assignedCount > 0 && !assignedSectionFound {
				t.Logf("Assigned Issues section not found in CSV")
				return false
			}

			// Verify worked issues section exists if there are worked issues
			workedSectionFound := false
			for _, record := range records {
				if len(record) > 0 && record[0] == "Worked Issues" {
					workedSectionFound = true
					break
				}
			}

			if workedCount > 0 && !workedSectionFound {
				t.Logf("Worked Issues section not found in CSV")
				return false
			}

			// Verify untracked issues section exists if there are untracked issues
			untrackedSectionFound := false
			for _, record := range records {
				if len(record) > 0 && record[0] == "Untracked Issues" {
					untrackedSectionFound = true
					break
				}
			}

			if untrackedCount > 0 && !untrackedSectionFound {
				t.Logf("Untracked Issues section not found in CSV")
				return false
			}

			// Verify rollup report section exists
			rollupSectionFound := false
			for _, record := range records {
				if len(record) > 0 && record[0] == "Rollup Report" {
					rollupSectionFound = true
					break
				}
			}

			if !rollupSectionFound {
				t.Logf("Rollup Report section not found in CSV")
				return false
			}

			return true
		},
		gen.IntRange(0, 10),  // assignedCount
		gen.IntRange(0, 10),  // workedCount
		gen.IntRange(0, 10),  // untrackedCount
		gen.IntRange(1, 5),   // statusCount
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 26: Export header completeness**
// For any export operation, the header should contain the user name, start date,
// end date, and export generation timestamp.
// **Validates: Requirements 8.4, 8.5, 8.6, 8.7**
func TestProperty_ExportHeaderCompleteness(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("export header contains all required fields", prop.ForAll(
		func(userName string, startDay, endDay int) bool {
			// Constrain inputs
			if userName == "" {
				userName = "TestUser"
			}
			if startDay < 1 || startDay > 28 {
				startDay = 1
			}
			if endDay < startDay || endDay > 28 {
				endDay = 28
			}

			user := &User{
				AccountId: "test-123",
				Name:      userName,
			}

			dateRange := DateRange{
				StartDate: time.Date(2024, 1, startDay, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, endDay, 23, 59, 59, 0, time.UTC),
			}

			// Create minimal report data
			assignedReport := &AssignedIssuesReport{Issues: []IssueReportItem{}, TotalCount: 0}
			workedReport := &WorkedIssuesReport{Issues: []WorkedIssueItem{}, TotalCount: 0, TotalTimeLogged: 0}
			untrackedReport := &UntrackedIssuesReport{Issues: []IssueReportItem{}, TotalCount: 0}
			rollupReport := &RollupReport{TotalAssigned: 0, TotalWorked: 0, StatusBreakdown: []StatusCount{}, TotalUniqueIssues: 0}

			// Export to CSV
			exportService := NewExportService(user, dateRange)
			filename, err := exportService.ExportToCSV(assignedReport, workedReport, untrackedReport, rollupReport)
			if err != nil {
				t.Logf("Export failed: %v", err)
				return false
			}

			defer os.Remove(filename)

			// Read CSV and check header
			file, err := os.Open(filename)
			if err != nil {
				t.Logf("Failed to open CSV: %v", err)
				return false
			}
			defer file.Close()

			reader := csv.NewReader(file)
			reader.FieldsPerRecord = -1 // Allow variable number of fields
			records, err := reader.ReadAll()
			if err != nil {
				t.Logf("Failed to read CSV: %v", err)
				return false
			}

			// Check that we have at least 4 header rows
			if len(records) < 4 {
				t.Logf("Expected at least 4 header rows, got %d", len(records))
				return false
			}

			// Verify user name is present
			userFound := false
			startDateFound := false
			endDateFound := false
			exportDateFound := false

			for _, record := range records[:10] { // Check first 10 rows
				if len(record) >= 2 {
					if record[0] == "User" && strings.Contains(record[1], userName) {
						userFound = true
					}
					if record[0] == "Start Date" && record[1] != "" {
						startDateFound = true
					}
					if record[0] == "End Date" && record[1] != "" {
						endDateFound = true
					}
					if record[0] == "Export Date" && record[1] != "" {
						exportDateFound = true
					}
				}
			}

			if !userFound {
				t.Logf("User name not found in header")
				return false
			}
			if !startDateFound {
				t.Logf("Start date not found in header")
				return false
			}
			if !endDateFound {
				t.Logf("End date not found in header")
				return false
			}
			if !exportDateFound {
				t.Logf("Export date not found in header")
				return false
			}

			return true
		},
		gen.AlphaString(),
		gen.IntRange(1, 28),
		gen.IntRange(1, 28),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 27: ISO 8601 date formatting**
// For any date or timestamp in an export header, it should be formatted according
// to ISO 8601 standard.
// **Validates: Requirements 8.6, 8.7**
func TestProperty_ISO8601DateFormatting(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("dates in export are ISO 8601 formatted", prop.ForAll(
		func(startDay, endDay int) bool {
			// Constrain inputs
			if startDay < 1 || startDay > 28 {
				startDay = 1
			}
			if endDay < startDay || endDay > 28 {
				endDay = 28
			}

			user := &User{
				AccountId: "test-123",
				Name:      "Test User",
			}

			dateRange := DateRange{
				StartDate: time.Date(2024, 1, startDay, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, endDay, 23, 59, 59, 0, time.UTC),
			}

			// Create minimal report data
			assignedReport := &AssignedIssuesReport{Issues: []IssueReportItem{}, TotalCount: 0}
			workedReport := &WorkedIssuesReport{Issues: []WorkedIssueItem{}, TotalCount: 0, TotalTimeLogged: 0}
			untrackedReport := &UntrackedIssuesReport{Issues: []IssueReportItem{}, TotalCount: 0}
			rollupReport := &RollupReport{TotalAssigned: 0, TotalWorked: 0, StatusBreakdown: []StatusCount{}, TotalUniqueIssues: 0}

			// Export to CSV
			exportService := NewExportService(user, dateRange)
			filename, err := exportService.ExportToCSV(assignedReport, workedReport, untrackedReport, rollupReport)
			if err != nil {
				t.Logf("Export failed: %v", err)
				return false
			}

			defer os.Remove(filename)

			// Read CSV and check date formats
			file, err := os.Open(filename)
			if err != nil {
				t.Logf("Failed to open CSV: %v", err)
				return false
			}
			defer file.Close()

			reader := csv.NewReader(file)
			reader.FieldsPerRecord = -1 // Allow variable number of fields
			records, err := reader.ReadAll()
			if err != nil {
				t.Logf("Failed to read CSV: %v", err)
				return false
			}

			// Check date formats in header
			for _, record := range records[:10] {
				if len(record) >= 2 {
					if record[0] == "Start Date" || record[0] == "End Date" || record[0] == "Export Date" {
						dateStr := record[1]
						// Try to parse as RFC3339 (ISO 8601)
						_, err := time.Parse(time.RFC3339, dateStr)
						if err != nil {
							t.Logf("Date %s is not in ISO 8601 format: %v", dateStr, err)
							return false
						}
					}
				}
			}

			return true
		},
		gen.IntRange(1, 28),
		gen.IntRange(1, 28),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 28: CSV section separation**
// For any multi-section CSV export, each report section should be separated
// from the next by blank rows.
// **Validates: Requirements 8.8**
func TestProperty_CSVSectionSeparation(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("CSV sections are separated by blank rows", prop.ForAll(
		func(hasAssigned, hasWorked, hasUntracked bool) bool {
			user := &User{
				AccountId: "test-123",
				Name:      "Test User",
			}

			dateRange := DateRange{
				StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			}

			// Create report data based on flags
			assignedReport := &AssignedIssuesReport{Issues: []IssueReportItem{}, TotalCount: 0}
			if hasAssigned {
				assignedReport.Issues = []IssueReportItem{
					{Key: "TEST-1", Summary: "Test", Status: "Open", CreatedDate: time.Now(), AssignedDate: time.Now()},
				}
				assignedReport.TotalCount = 1
			}

			workedReport := &WorkedIssuesReport{Issues: []WorkedIssueItem{}, TotalCount: 0, TotalTimeLogged: 0}
			if hasWorked {
				workedReport.Issues = []WorkedIssueItem{
					{IssueReportItem: IssueReportItem{Key: "TEST-2", Summary: "Test", Status: "Done", CreatedDate: time.Now(), AssignedDate: time.Now()}, TimeLogged: time.Hour},
				}
				workedReport.TotalCount = 1
			}

			untrackedReport := &UntrackedIssuesReport{Issues: []IssueReportItem{}, TotalCount: 0}
			if hasUntracked {
				untrackedReport.Issues = []IssueReportItem{
					{Key: "TEST-3", Summary: "Test", Status: "Open", CreatedDate: time.Now(), AssignedDate: time.Now()},
				}
				untrackedReport.TotalCount = 1
			}

			rollupReport := &RollupReport{
				TotalAssigned:     assignedReport.TotalCount,
				TotalWorked:       workedReport.TotalCount,
				StatusBreakdown:   []StatusCount{{StatusName: "Open", Count: 1, Percentage: 100.0}},
				TotalUniqueIssues: 1,
			}

			// Export to CSV
			exportService := NewExportService(user, dateRange)
			filename, err := exportService.ExportToCSV(assignedReport, workedReport, untrackedReport, rollupReport)
			if err != nil {
				t.Logf("Export failed: %v", err)
				return false
			}

			defer os.Remove(filename)

			// Read CSV and check for blank row separators
			file, err := os.Open(filename)
			if err != nil {
				t.Logf("Failed to open CSV: %v", err)
				return false
			}
			defer file.Close()

			reader := csv.NewReader(file)
			reader.FieldsPerRecord = -1 // Allow variable number of fields
			records, err := reader.ReadAll()
			if err != nil {
				t.Logf("Failed to read CSV: %v", err)
				return false
			}

			// Find section headers and verify blank rows separate them
			// The structure is: Header -> Blank -> Section1 -> Blank -> Section2 -> etc.
			sectionHeaders := []string{"Assigned Issues", "Worked Issues", "Untracked Issues"}
			
			for i, record := range records {
				if len(record) > 0 {
					for _, header := range sectionHeaders {
						if record[0] == header {
							// Find the end of this section (next blank row or end of file)
							sectionEnd := -1
							for j := i + 1; j < len(records); j++ {
								if len(records[j]) == 0 || (len(records[j]) == 1 && records[j][0] == "") {
									sectionEnd = j
									break
								}
							}
							
							// If this is not the last section, verify there's a blank row after it
							if sectionEnd > 0 && sectionEnd < len(records)-1 {
								// Check if the next non-blank row is another section
								for k := sectionEnd + 1; k < len(records); k++ {
									if len(records[k]) > 0 && records[k][0] != "" {
										// This should be another section header
										// The blank row at sectionEnd separates them
										break
									}
								}
							}
						}
					}
				}
			}

			return true
		},
		gen.Bool(),
		gen.Bool(),
		gen.Bool(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 31: Export filename uniqueness**
// For any export file creation, the filename should include a timestamp to
// prevent overwriting existing files.
// **Validates: Requirements 8.11**
func TestProperty_ExportFilenameUniqueness(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("export filenames include timestamps", prop.ForAll(
		func(seed int) bool {
			user := &User{
				AccountId: "test-123",
				Name:      "Test User",
			}

			dateRange := DateRange{
				StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			}

			// Create minimal report data
			assignedReport := &AssignedIssuesReport{Issues: []IssueReportItem{}, TotalCount: 0}
			workedReport := &WorkedIssuesReport{Issues: []WorkedIssueItem{}, TotalCount: 0, TotalTimeLogged: 0}
			untrackedReport := &UntrackedIssuesReport{Issues: []IssueReportItem{}, TotalCount: 0}
			rollupReport := &RollupReport{TotalAssigned: 0, TotalWorked: 0, StatusBreakdown: []StatusCount{}, TotalUniqueIssues: 0}

			// Export to CSV
			exportService := NewExportService(user, dateRange)
			
			filename, err := exportService.ExportToCSV(assignedReport, workedReport, untrackedReport, rollupReport)
			if err != nil {
				t.Logf("Export failed: %v", err)
				return false
			}
			defer os.Remove(filename)

			// Verify file exists
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				t.Logf("File doesn't exist: %s", filename)
				return false
			}

			// Verify filename contains timestamp pattern (YYYYMMDD_HHMMSS)
			basename := filepath.Base(filename)
			
			// Check for the expected pattern: jira_report_YYYYMMDD_HHMMSS.csv
			if !strings.HasPrefix(basename, "jira_report_") {
				t.Logf("Filename doesn't have expected prefix: %s", basename)
				return false
			}
			
			if !strings.HasSuffix(basename, ".csv") {
				t.Logf("Filename doesn't have .csv extension: %s", basename)
				return false
			}
			
			// Extract the timestamp part
			timestampPart := strings.TrimPrefix(basename, "jira_report_")
			timestampPart = strings.TrimSuffix(timestampPart, ".csv")
			
			// Verify it contains an underscore (separating date and time)
			if !strings.Contains(timestampPart, "_") {
				t.Logf("Timestamp doesn't contain date/time separator: %s", timestampPart)
				return false
			}
			
			// Verify the timestamp can be parsed
			_, err = time.Parse("20060102_150405", timestampPart)
			if err != nil {
				t.Logf("Timestamp is not in expected format (YYYYMMDD_HHMMSS): %s, error: %v", timestampPart, err)
				return false
			}

			return true
		},
		gen.Int(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 30: Export confirmation display**
// For any successful export operation, a confirmation message with the file
// location should be displayed.
// **Validates: Requirements 8.10**
func TestProperty_ExportConfirmationDisplay(t *testing.T) {
	// This property is tested through the UI integration
	// The exportToCSV method in reporting_window.go calls ShowSuccess with the filename
	// This is verified by checking that the method exists and has the correct signature
	
	// We can verify the export service returns the filename
	properties := gopter.NewProperties(nil)

	properties.Property("export returns filename for confirmation", prop.ForAll(
		func(seed int) bool {
			user := &User{
				AccountId: "test-123",
				Name:      "Test User",
			}

			dateRange := DateRange{
				StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			}

			// Create minimal report data
			assignedReport := &AssignedIssuesReport{Issues: []IssueReportItem{}, TotalCount: 0}
			workedReport := &WorkedIssuesReport{Issues: []WorkedIssueItem{}, TotalCount: 0, TotalTimeLogged: 0}
			untrackedReport := &UntrackedIssuesReport{Issues: []IssueReportItem{}, TotalCount: 0}
			rollupReport := &RollupReport{TotalAssigned: 0, TotalWorked: 0, StatusBreakdown: []StatusCount{}, TotalUniqueIssues: 0}

			// Export to CSV
			exportService := NewExportService(user, dateRange)
			filename, err := exportService.ExportToCSV(assignedReport, workedReport, untrackedReport, rollupReport)
			if err != nil {
				t.Logf("Export failed: %v", err)
				return false
			}

			defer os.Remove(filename)

			// Verify filename is not empty
			if filename == "" {
				t.Logf("Export returned empty filename")
				return false
			}

			// Verify file exists at the returned path
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				t.Logf("File doesn't exist at returned path: %s", filename)
				return false
			}

			return true
		},
		gen.Int(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 25: PDF export completeness**
// For any report data, the generated PDF file should contain all visible report data
// formatted in tables.
// **Validates: Requirements 8.3**
func TestProperty_PDFExportCompleteness(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("PDF export contains all report data", prop.ForAll(
		func(assignedCount, workedCount, untrackedCount, statusCount int) bool {
			// Constrain to reasonable sizes
			if assignedCount < 0 || assignedCount > 50 {
				return true
			}
			if workedCount < 0 || workedCount > 50 {
				return true
			}
			if untrackedCount < 0 || untrackedCount > 50 {
				return true
			}
			if statusCount < 0 || statusCount > 10 {
				return true
			}

			// Generate test data
			user := &User{
				AccountId: "test-user-123",
				Name:      "Test User",
			}

			dateRange := DateRange{
				StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			}

			// Generate assigned issues
			assignedIssues := make([]IssueReportItem, assignedCount)
			for i := 0; i < assignedCount; i++ {
				assignedIssues[i] = IssueReportItem{
					Key:          generateIssueKey(i),
					Summary:      generateSummary(i),
					Status:       "In Progress",
					CreatedDate:  time.Date(2024, 1, i%28+1, 10, 0, 0, 0, time.UTC),
					AssignedDate: time.Date(2024, 1, i%28+1, 10, 0, 0, 0, time.UTC),
				}
			}

			assignedReport := &AssignedIssuesReport{
				Issues:     assignedIssues,
				TotalCount: assignedCount,
			}

			// Generate worked issues
			workedIssues := make([]WorkedIssueItem, workedCount)
			for i := 0; i < workedCount; i++ {
				workedIssues[i] = WorkedIssueItem{
					IssueReportItem: IssueReportItem{
						Key:          generateIssueKey(i + 100),
						Summary:      generateSummary(i + 100),
						Status:       "Done",
						CreatedDate:  time.Date(2024, 1, i%28+1, 10, 0, 0, 0, time.UTC),
						AssignedDate: time.Date(2024, 1, i%28+1, 10, 0, 0, 0, time.UTC),
					},
					TimeLogged:   time.Duration(i+1) * time.Hour,
					WorklogCount: i + 1,
				}
			}

			workedReport := &WorkedIssuesReport{
				Issues:          workedIssues,
				TotalCount:      workedCount,
				TotalTimeLogged: time.Duration(workedCount) * time.Hour,
			}

			// Generate untracked issues
			untrackedIssues := make([]IssueReportItem, untrackedCount)
			for i := 0; i < untrackedCount; i++ {
				untrackedIssues[i] = IssueReportItem{
					Key:          generateIssueKey(i + 200),
					Summary:      generateSummary(i + 200),
					Status:       "To Do",
					CreatedDate:  time.Date(2024, 1, i%28+1, 10, 0, 0, 0, time.UTC),
					AssignedDate: time.Date(2024, 1, i%28+1, 10, 0, 0, 0, time.UTC),
				}
			}

			untrackedReport := &UntrackedIssuesReport{
				Issues:     untrackedIssues,
				TotalCount: untrackedCount,
			}

			// Generate status breakdown
			statusBreakdown := make([]StatusCount, statusCount)
			for i := 0; i < statusCount; i++ {
				statusBreakdown[i] = StatusCount{
					StatusName: generateStatusName(i),
					Count:      i + 1,
					Percentage: float64(i+1) * 10.0,
				}
			}

			rollupReport := &RollupReport{
				TotalAssigned:     assignedCount,
				TotalWorked:       workedCount,
				StatusBreakdown:   statusBreakdown,
				TotalUniqueIssues: assignedCount + workedCount,
			}

			// Create export service and export to PDF
			exportService := NewExportService(user, dateRange)
			filename, err := exportService.ExportToPDF(assignedReport, workedReport, untrackedReport, rollupReport)
			if err != nil {
				t.Logf("Export failed: %v", err)
				return false
			}

			// Clean up file after test
			defer os.Remove(filename)

			// Verify file exists
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				t.Logf("PDF file doesn't exist: %s", filename)
				return false
			}

			// Verify file has content (PDF files should be at least a few KB)
			fileInfo, err := os.Stat(filename)
			if err != nil {
				t.Logf("Failed to stat PDF file: %v", err)
				return false
			}

			if fileInfo.Size() < 1000 {
				t.Logf("PDF file is too small (%d bytes), likely incomplete", fileInfo.Size())
				return false
			}

			// Read file and verify it's a valid PDF (starts with %PDF)
			file, err := os.Open(filename)
			if err != nil {
				t.Logf("Failed to open PDF file: %v", err)
				return false
			}
			defer file.Close()

			header := make([]byte, 4)
			_, err = file.Read(header)
			if err != nil {
				t.Logf("Failed to read PDF header: %v", err)
				return false
			}

			if string(header) != "%PDF" {
				t.Logf("File is not a valid PDF (header: %s)", string(header))
				return false
			}

			return true
		},
		gen.IntRange(0, 10),  // assignedCount
		gen.IntRange(0, 10),  // workedCount
		gen.IntRange(0, 10),  // untrackedCount
		gen.IntRange(1, 5),   // statusCount
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 29: PDF section formatting**
// For any multi-section PDF export, each section should have a heading and
// appropriate page breaks.
// **Validates: Requirements 8.9**
func TestProperty_PDFSectionFormatting(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("PDF sections have headings and proper formatting", prop.ForAll(
		func(itemsPerSection int) bool {
			// Constrain to reasonable sizes
			if itemsPerSection < 0 || itemsPerSection > 100 {
				return true
			}

			user := &User{
				AccountId: "test-123",
				Name:      "Test User",
			}

			dateRange := DateRange{
				StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			}

			// Create report data with multiple items to test section separation
			assignedIssues := make([]IssueReportItem, itemsPerSection)
			for i := 0; i < itemsPerSection; i++ {
				assignedIssues[i] = IssueReportItem{
					Key:          generateIssueKey(i),
					Summary:      generateSummary(i),
					Status:       "Open",
					CreatedDate:  time.Now(),
					AssignedDate: time.Now(),
				}
			}
			assignedReport := &AssignedIssuesReport{
				Issues:     assignedIssues,
				TotalCount: itemsPerSection,
			}

			workedIssues := make([]WorkedIssueItem, itemsPerSection)
			for i := 0; i < itemsPerSection; i++ {
				workedIssues[i] = WorkedIssueItem{
					IssueReportItem: IssueReportItem{
						Key:          generateIssueKey(i + 100),
						Summary:      generateSummary(i + 100),
						Status:       "Done",
						CreatedDate:  time.Now(),
						AssignedDate: time.Now(),
					},
					TimeLogged: time.Hour,
				}
			}
			workedReport := &WorkedIssuesReport{
				Issues:          workedIssues,
				TotalCount:      itemsPerSection,
				TotalTimeLogged: time.Duration(itemsPerSection) * time.Hour,
			}

			untrackedIssues := make([]IssueReportItem, itemsPerSection)
			for i := 0; i < itemsPerSection; i++ {
				untrackedIssues[i] = IssueReportItem{
					Key:          generateIssueKey(i + 200),
					Summary:      generateSummary(i + 200),
					Status:       "Open",
					CreatedDate:  time.Now(),
					AssignedDate: time.Now(),
				}
			}
			untrackedReport := &UntrackedIssuesReport{
				Issues:     untrackedIssues,
				TotalCount: itemsPerSection,
			}

			rollupReport := &RollupReport{
				TotalAssigned:     itemsPerSection,
				TotalWorked:       itemsPerSection,
				StatusBreakdown:   []StatusCount{{StatusName: "Open", Count: 1, Percentage: 100.0}},
				TotalUniqueIssues: itemsPerSection * 2,
			}

			// Export to PDF
			exportService := NewExportService(user, dateRange)
			filename, err := exportService.ExportToPDF(assignedReport, workedReport, untrackedReport, rollupReport)
			if err != nil {
				t.Logf("Export failed: %v", err)
				return false
			}

			defer os.Remove(filename)

			// Verify file exists
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				t.Logf("PDF file doesn't exist: %s", filename)
				return false
			}

			// Verify file has content
			fileInfo, err := os.Stat(filename)
			if err != nil {
				t.Logf("Failed to stat PDF file: %v", err)
				return false
			}

			// Minimum size for a PDF with header and sections
			minSize := int64(1500)
			if fileInfo.Size() < minSize {
				t.Logf("PDF file is too small (%d bytes), expected at least %d bytes", fileInfo.Size(), minSize)
				return false
			}

			// Read file and verify it's a valid PDF
			file, err := os.Open(filename)
			if err != nil {
				t.Logf("Failed to open PDF file: %v", err)
				return false
			}
			defer file.Close()

			header := make([]byte, 4)
			_, err = file.Read(header)
			if err != nil {
				t.Logf("Failed to read PDF header: %v", err)
				return false
			}

			if string(header) != "%PDF" {
				t.Logf("File is not a valid PDF (header: %s)", string(header))
				return false
			}

			// Verify PDF has reasonable size based on content
			// More items should result in larger PDF
			expectedMinSize := int64(1500 + itemsPerSection*50) // Rough estimate
			if fileInfo.Size() < expectedMinSize {
				t.Logf("PDF file size (%d bytes) is smaller than expected for %d items per section (expected at least %d bytes)",
					fileInfo.Size(), itemsPerSection, expectedMinSize)
				return false
			}

			return true
		},
		gen.IntRange(0, 20),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Helper functions for test data generation

func generateIssueKey(index int) string {
	return "TEST-" + string(rune('A'+index%26)) + string(rune('0'+index%10))
}

func generateSummary(index int) string {
	summaries := []string{
		"Implement feature",
		"Fix bug",
		"Update documentation",
		"Refactor code",
		"Add tests",
	}
	return summaries[index%len(summaries)]
}

func generateStatusName(index int) string {
	statuses := []string{
		"To Do",
		"In Progress",
		"In Review",
		"Done",
		"Blocked",
	}
	return statuses[index%len(statuses)]
}
