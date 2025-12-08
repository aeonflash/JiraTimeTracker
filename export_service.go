package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// ExportService handles report data export to various formats
type ExportService struct {
	User      *User
	DateRange DateRange
}

// NewExportService creates a new ExportService instance
func NewExportService(user *User, dateRange DateRange) *ExportService {
	return &ExportService{
		User:      user,
		DateRange: dateRange,
	}
}

// ExportToCSV generates a CSV file containing all report data (deprecated - use ExportToCSVFile)
func (es *ExportService) ExportToCSV(
	assignedReport *AssignedIssuesReport,
	workedReport *WorkedIssuesReport,
	untrackedReport *UntrackedIssuesReport,
	rollupReport *RollupReport,
) (string, error) {
	// Generate filename with timestamp
	filename := es.generateFilename("csv")
	return filename, es.ExportToCSVFile(assignedReport, workedReport, untrackedReport, rollupReport, filename)
}

// ExportToCSVFile generates a CSV file containing all report data at the specified path
func (es *ExportService) ExportToCSVFile(
	assignedReport *AssignedIssuesReport,
	workedReport *WorkedIssuesReport,
	untrackedReport *UntrackedIssuesReport,
	rollupReport *RollupReport,
	filepath string,
) error {
	// Create the file
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	// Write header section
	if err := es.writeCSVHeader(writer); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}
	
	// Write blank row separator
	if err := writer.Write([]string{}); err != nil {
		return err
	}
	
	// Write assigned issues section
	if err := es.writeAssignedIssuesCSV(writer, assignedReport); err != nil {
		return fmt.Errorf("failed to write assigned issues: %w", err)
	}
	
	// Write blank row separator
	if err := writer.Write([]string{}); err != nil {
		return err
	}
	
	// Write worked issues section
	if err := es.writeWorkedIssuesCSV(writer, workedReport); err != nil {
		return fmt.Errorf("failed to write worked issues: %w", err)
	}
	
	// Write blank row separator
	if err := writer.Write([]string{}); err != nil {
		return err
	}
	
	// Write untracked issues section
	if err := es.writeUntrackedIssuesCSV(writer, untrackedReport); err != nil {
		return fmt.Errorf("failed to write untracked issues: %w", err)
	}
	
	// Write blank row separator
	if err := writer.Write([]string{}); err != nil {
		return err
	}
	
	// Write rollup report section
	if err := es.writeRollupReportCSV(writer, rollupReport); err != nil {
		return fmt.Errorf("failed to write rollup report: %w", err)
	}
	
	return nil
}

// writeCSVHeader writes the header section with user, date range, and export date
func (es *ExportService) writeCSVHeader(writer *csv.Writer) error {
	exportDate := time.Now()
	
	// Write user information
	if err := writer.Write([]string{"User", es.getUserName()}); err != nil {
		return err
	}
	
	// Write date range in ISO 8601 format
	if err := writer.Write([]string{"Start Date", es.DateRange.StartDate.Format(time.RFC3339)}); err != nil {
		return err
	}
	
	if err := writer.Write([]string{"End Date", es.DateRange.EndDate.Format(time.RFC3339)}); err != nil {
		return err
	}
	
	// Write export generation timestamp in ISO 8601 format
	if err := writer.Write([]string{"Export Date", exportDate.Format(time.RFC3339)}); err != nil {
		return err
	}
	
	return nil
}

// writeAssignedIssuesCSV writes the assigned issues section
func (es *ExportService) writeAssignedIssuesCSV(writer *csv.Writer, report *AssignedIssuesReport) error {
	// Section header
	if err := writer.Write([]string{"Assigned Issues"}); err != nil {
		return err
	}
	
	// Column headers
	if err := writer.Write([]string{"Issue Key", "Summary", "Status", "Created Date"}); err != nil {
		return err
	}
	
	// Data rows
	for _, issue := range report.Issues {
		row := []string{
			issue.Key,
			issue.Summary,
			issue.Status,
			issue.CreatedDate.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	
	// Total count
	if err := writer.Write([]string{"Total", fmt.Sprintf("%d", report.TotalCount)}); err != nil {
		return err
	}
	
	return nil
}

// writeWorkedIssuesCSV writes the worked issues section
func (es *ExportService) writeWorkedIssuesCSV(writer *csv.Writer, report *WorkedIssuesReport) error {
	// Section header
	if err := writer.Write([]string{"Worked Issues"}); err != nil {
		return err
	}
	
	// Column headers
	if err := writer.Write([]string{"Issue Key", "Summary", "Status", "Time Logged"}); err != nil {
		return err
	}
	
	// Data rows
	for _, issue := range report.Issues {
		row := []string{
			issue.Key,
			issue.Summary,
			issue.Status,
			formatDuration(issue.TimeLogged),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	
	// Total count and time
	if err := writer.Write([]string{
		"Total",
		fmt.Sprintf("%d issues", report.TotalCount),
		"",
		formatDuration(report.TotalTimeLogged),
	}); err != nil {
		return err
	}
	
	return nil
}

// writeUntrackedIssuesCSV writes the untracked issues section
func (es *ExportService) writeUntrackedIssuesCSV(writer *csv.Writer, report *UntrackedIssuesReport) error {
	// Section header
	if err := writer.Write([]string{"Untracked Issues"}); err != nil {
		return err
	}
	
	// Column headers
	if err := writer.Write([]string{"Issue Key", "Summary", "Status", "Created Date"}); err != nil {
		return err
	}
	
	// Data rows
	for _, issue := range report.Issues {
		row := []string{
			issue.Key,
			issue.Summary,
			issue.Status,
			issue.CreatedDate.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	
	// Total count
	if err := writer.Write([]string{"Total", fmt.Sprintf("%d", report.TotalCount)}); err != nil {
		return err
	}
	
	return nil
}

// writeRollupReportCSV writes the rollup report section
func (es *ExportService) writeRollupReportCSV(writer *csv.Writer, report *RollupReport) error {
	// Section header
	if err := writer.Write([]string{"Rollup Report"}); err != nil {
		return err
	}
	
	// Summary statistics
	if err := writer.Write([]string{"Total Assigned Issues", fmt.Sprintf("%d", report.TotalAssigned)}); err != nil {
		return err
	}
	
	if err := writer.Write([]string{"Total Worked Issues", fmt.Sprintf("%d", report.TotalWorked)}); err != nil {
		return err
	}
	
	if err := writer.Write([]string{"Total Unique Issues", fmt.Sprintf("%d", report.TotalUniqueIssues)}); err != nil {
		return err
	}
	
	// Blank row before status breakdown
	if err := writer.Write([]string{}); err != nil {
		return err
	}
	
	// Status breakdown header
	if err := writer.Write([]string{"Status Breakdown"}); err != nil {
		return err
	}
	
	if err := writer.Write([]string{"Status", "Count", "Percentage"}); err != nil {
		return err
	}
	
	// Status breakdown data
	for _, status := range report.StatusBreakdown {
		row := []string{
			status.StatusName,
			fmt.Sprintf("%d", status.Count),
			fmt.Sprintf("%.2f%%", status.Percentage),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	
	return nil
}

// generateFilename creates a filename with timestamp to prevent overwrites
func (es *ExportService) generateFilename(extension string) string {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("jira_report_%s.%s", timestamp, extension)
	
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return filename
	}
	
	// Create reports directory if it doesn't exist
	reportsDir := filepath.Join(homeDir, "JiraReports")
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		// Fallback to current directory
		return filename
	}
	
	return filepath.Join(reportsDir, filename)
}

// getUserName returns the user's name or account ID
func (es *ExportService) getUserName() string {
	if es.User.Name != "" {
		return es.User.Name
	}
	return es.User.AccountId
}

// ExportToPDF generates a PDF file containing all report data (deprecated - use ExportToPDFFile)
func (es *ExportService) ExportToPDF(
	assignedReport *AssignedIssuesReport,
	workedReport *WorkedIssuesReport,
	untrackedReport *UntrackedIssuesReport,
	rollupReport *RollupReport,
) (string, error) {
	// Generate filename with timestamp
	filename := es.generateFilename("pdf")
	return filename, es.ExportToPDFFile(assignedReport, workedReport, untrackedReport, rollupReport, filename)
}

// ExportToPDFFile generates a PDF file containing all report data at the specified path
func (es *ExportService) ExportToPDFFile(
	assignedReport *AssignedIssuesReport,
	workedReport *WorkedIssuesReport,
	untrackedReport *UntrackedIssuesReport,
	rollupReport *RollupReport,
	filepath string,
) error {
	// Create new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 10)

	// Add first page
	pdf.AddPage()

	// Write header section
	es.writePDFHeader(pdf)

	// Add page break after header
	pdf.Ln(10)

	// Write assigned issues section
	es.writeAssignedIssuesPDF(pdf, assignedReport)

	// Write worked issues section
	es.writeWorkedIssuesPDF(pdf, workedReport)

	// Write untracked issues section
	es.writeUntrackedIssuesPDF(pdf, untrackedReport)

	// Write rollup report section
	es.writeRollupReportPDF(pdf, rollupReport)

	// Save PDF to file
	err := pdf.OutputFileAndClose(filepath)
	if err != nil {
		return fmt.Errorf("failed to create PDF file: %w", err)
	}

	return nil
}

// writePDFHeader writes the header section with user, date range, and export date
func (es *ExportService) writePDFHeader(pdf *gofpdf.Fpdf) {
	exportDate := time.Now()

	// Set font for header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Jira Report")
	pdf.Ln(12)

	// Set font for header details
	pdf.SetFont("Arial", "", 10)

	// Write user information
	pdf.Cell(40, 6, "User:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, es.getUserName())
	pdf.Ln(6)

	// Write date range in ISO 8601 format
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 6, "Start Date:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, es.DateRange.StartDate.Format(time.RFC3339))
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 6, "End Date:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, es.DateRange.EndDate.Format(time.RFC3339))
	pdf.Ln(6)

	// Write export generation timestamp in ISO 8601 format
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 6, "Export Date:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, exportDate.Format(time.RFC3339))
	pdf.Ln(6)
}

// writeAssignedIssuesPDF writes the assigned issues section
func (es *ExportService) writeAssignedIssuesPDF(pdf *gofpdf.Fpdf, report *AssignedIssuesReport) {
	// Check if we need a new page
	if pdf.GetY() > 250 {
		pdf.AddPage()
	}

	// Section heading
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Assigned Issues")
	pdf.Ln(10)

	if len(report.Issues) == 0 {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 6, "No assigned issues found")
		pdf.Ln(10)
		return
	}

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(30, 7, "Issue Key", "1", 0, "L", true, 0, "")
	pdf.CellFormat(80, 7, "Summary", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 7, "Status", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 7, "Created Date", "1", 1, "L", true, 0, "")

	// Table data
	pdf.SetFont("Arial", "", 8)
	for _, issue := range report.Issues {
		// Check if we need a new page
		if pdf.GetY() > 270 {
			pdf.AddPage()
			// Repeat header on new page
			pdf.SetFont("Arial", "B", 9)
			pdf.SetFillColor(200, 200, 200)
			pdf.CellFormat(30, 7, "Issue Key", "1", 0, "L", true, 0, "")
			pdf.CellFormat(80, 7, "Summary", "1", 0, "L", true, 0, "")
			pdf.CellFormat(30, 7, "Status", "1", 0, "L", true, 0, "")
			pdf.CellFormat(40, 7, "Created Date", "1", 1, "L", true, 0, "")
			pdf.SetFont("Arial", "", 8)
		}

		pdf.CellFormat(30, 6, issue.Key, "1", 0, "L", false, 0, "")
		pdf.CellFormat(80, 6, truncateString(issue.Summary, 50), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 6, issue.Status, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 6, issue.CreatedDate.Format("2006-01-02"), "1", 1, "L", false, 0, "")
	}

	// Total count
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 7, fmt.Sprintf("Total: %d", report.TotalCount))
	pdf.Ln(10)
}

// writeWorkedIssuesPDF writes the worked issues section
func (es *ExportService) writeWorkedIssuesPDF(pdf *gofpdf.Fpdf, report *WorkedIssuesReport) {
	// Check if we need a new page
	if pdf.GetY() > 250 {
		pdf.AddPage()
	}

	// Section heading
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Worked Issues")
	pdf.Ln(10)

	if len(report.Issues) == 0 {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 6, "No worked issues found")
		pdf.Ln(10)
		return
	}

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(30, 7, "Issue Key", "1", 0, "L", true, 0, "")
	pdf.CellFormat(70, 7, "Summary", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 7, "Status", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 7, "Time Logged", "1", 1, "L", true, 0, "")

	// Table data
	pdf.SetFont("Arial", "", 8)
	for _, issue := range report.Issues {
		// Check if we need a new page
		if pdf.GetY() > 270 {
			pdf.AddPage()
			// Repeat header on new page
			pdf.SetFont("Arial", "B", 9)
			pdf.SetFillColor(200, 200, 200)
			pdf.CellFormat(30, 7, "Issue Key", "1", 0, "L", true, 0, "")
			pdf.CellFormat(70, 7, "Summary", "1", 0, "L", true, 0, "")
			pdf.CellFormat(30, 7, "Status", "1", 0, "L", true, 0, "")
			pdf.CellFormat(50, 7, "Time Logged", "1", 1, "L", true, 0, "")
			pdf.SetFont("Arial", "", 8)
		}

		pdf.CellFormat(30, 6, issue.Key, "1", 0, "L", false, 0, "")
		pdf.CellFormat(70, 6, truncateString(issue.Summary, 45), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 6, issue.Status, "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 6, formatDuration(issue.TimeLogged), "1", 1, "L", false, 0, "")
	}

	// Total count and time
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 7, fmt.Sprintf("Total: %d issues, %s", report.TotalCount, formatDuration(report.TotalTimeLogged)))
	pdf.Ln(10)
}

// writeUntrackedIssuesPDF writes the untracked issues section
func (es *ExportService) writeUntrackedIssuesPDF(pdf *gofpdf.Fpdf, report *UntrackedIssuesReport) {
	// Check if we need a new page
	if pdf.GetY() > 250 {
		pdf.AddPage()
	}

	// Section heading
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Untracked Issues")
	pdf.Ln(10)

	if len(report.Issues) == 0 {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 6, "No untracked issues found")
		pdf.Ln(10)
		return
	}

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(30, 7, "Issue Key", "1", 0, "L", true, 0, "")
	pdf.CellFormat(80, 7, "Summary", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 7, "Status", "1", 0, "L", true, 0, "")
	pdf.CellFormat(40, 7, "Created Date", "1", 1, "L", true, 0, "")

	// Table data
	pdf.SetFont("Arial", "", 8)
	for _, issue := range report.Issues {
		// Check if we need a new page
		if pdf.GetY() > 270 {
			pdf.AddPage()
			// Repeat header on new page
			pdf.SetFont("Arial", "B", 9)
			pdf.SetFillColor(200, 200, 200)
			pdf.CellFormat(30, 7, "Issue Key", "1", 0, "L", true, 0, "")
			pdf.CellFormat(80, 7, "Summary", "1", 0, "L", true, 0, "")
			pdf.CellFormat(30, 7, "Status", "1", 0, "L", true, 0, "")
			pdf.CellFormat(40, 7, "Created Date", "1", 1, "L", true, 0, "")
			pdf.SetFont("Arial", "", 8)
		}

		pdf.CellFormat(30, 6, issue.Key, "1", 0, "L", false, 0, "")
		pdf.CellFormat(80, 6, truncateString(issue.Summary, 50), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 6, issue.Status, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 6, issue.CreatedDate.Format("2006-01-02"), "1", 1, "L", false, 0, "")
	}

	// Total count
	pdf.SetFont("Arial", "B", 9)
	pdf.Cell(0, 7, fmt.Sprintf("Total: %d", report.TotalCount))
	pdf.Ln(10)
}

// writeRollupReportPDF writes the rollup report section
func (es *ExportService) writeRollupReportPDF(pdf *gofpdf.Fpdf, report *RollupReport) {
	// Check if we need a new page
	if pdf.GetY() > 250 {
		pdf.AddPage()
	}

	// Section heading
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Rollup Report")
	pdf.Ln(10)

	// Summary statistics
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(60, 6, "Total Assigned Issues:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, fmt.Sprintf("%d", report.TotalAssigned))
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(60, 6, "Total Worked Issues:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, fmt.Sprintf("%d", report.TotalWorked))
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.Cell(60, 6, "Total Unique Issues:")
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 6, fmt.Sprintf("%d", report.TotalUniqueIssues))
	pdf.Ln(12)

	// Status breakdown
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "Status Breakdown")
	pdf.Ln(10)

	if len(report.StatusBreakdown) == 0 {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 6, "No status data available")
		pdf.Ln(10)
		return
	}

	// Table header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(80, 7, "Status", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 7, "Count", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 7, "Percentage", "1", 1, "L", true, 0, "")

	// Table data
	pdf.SetFont("Arial", "", 10)
	for _, status := range report.StatusBreakdown {
		// Check if we need a new page
		if pdf.GetY() > 270 {
			pdf.AddPage()
			// Repeat header on new page
			pdf.SetFont("Arial", "B", 10)
			pdf.SetFillColor(200, 200, 200)
			pdf.CellFormat(80, 7, "Status", "1", 0, "L", true, 0, "")
			pdf.CellFormat(50, 7, "Count", "1", 0, "L", true, 0, "")
			pdf.CellFormat(50, 7, "Percentage", "1", 1, "L", true, 0, "")
			pdf.SetFont("Arial", "", 10)
		}

		pdf.CellFormat(80, 6, status.StatusName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 6, fmt.Sprintf("%d", status.Count), "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 6, fmt.Sprintf("%.2f%%", status.Percentage), "1", 1, "L", false, 0, "")
	}
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
