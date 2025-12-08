package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// AssignedIssuesView creates a view for displaying assigned issues
type AssignedIssuesView struct {
	container *fyne.Container
	table     *widget.Table
	data      *AssignedIssuesReport
}

// NewAssignedIssuesView creates a new assigned issues view
func NewAssignedIssuesView() *AssignedIssuesView {
	view := &AssignedIssuesView{
		data: &AssignedIssuesReport{
			Issues:     []IssueReportItem{},
			TotalCount: 0,
		},
	}

	// Create empty state message
	emptyLabel := widget.NewLabel("No assigned issues found. Click 'Refresh Reports' to load data.")
	view.container = container.NewCenter(emptyLabel)

	return view
}

// UpdateData updates the view with new report data
func (v *AssignedIssuesView) UpdateData(report *AssignedIssuesReport) {
	v.data = report

	if report == nil || report.TotalCount == 0 {
		// Show empty state
		emptyLabel := widget.NewLabel("No assigned issues found for the selected date range.")
		v.container = container.NewCenter(emptyLabel)
		return
	}

	// Create table with data
	v.table = widget.NewTable(
		func() (int, int) {
			return len(report.Issues) + 1, 4 // +1 for header row, 4 columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			// Header row
			if id.Row == 0 {
				switch id.Col {
				case 0:
					label.SetText("Issue Key")
				case 1:
					label.SetText("Summary")
				case 2:
					label.SetText("Status")
				case 3:
					label.SetText("Assignee")
				}
				label.TextStyle = fyne.TextStyle{Bold: true}
				return
			}

			// Data rows
			issue := report.Issues[id.Row-1]
			label.TextStyle = fyne.TextStyle{Bold: false}

			switch id.Col {
			case 0:
				label.SetText(issue.Key)
			case 1:
				label.SetText(issue.Summary)
			case 2:
				label.SetText(issue.Status)
			case 3:
				label.SetText("currentUser()") // Placeholder for assignee
			}
		},
	)

	// Set column widths
	v.table.SetColumnWidth(0, 130)  // Issue Key
	v.table.SetColumnWidth(1, 550)  // Summary - much wider
	v.table.SetColumnWidth(2, 100)  // Status
	v.table.SetColumnWidth(3, 100)  // Assignee

	v.container = container.NewBorder(nil, nil, nil, nil, v.table)
}

// GetContainer returns the container for this view
func (v *AssignedIssuesView) GetContainer() *fyne.Container {
	return v.container
}

// WorkedIssuesView creates a view for displaying worked issues
type WorkedIssuesView struct {
	container *fyne.Container
	table     *widget.Table
	data      *WorkedIssuesReport
}

// NewWorkedIssuesView creates a new worked issues view
func NewWorkedIssuesView() *WorkedIssuesView {
	view := &WorkedIssuesView{
		data: &WorkedIssuesReport{
			Issues:          []WorkedIssueItem{},
			TotalCount:      0,
			TotalTimeLogged: 0,
		},
	}

	// Create empty state message
	emptyLabel := widget.NewLabel("No worked issues found. Click 'Refresh Reports' to load data.")
	view.container = container.NewCenter(emptyLabel)

	return view
}

// UpdateData updates the view with new report data
func (v *WorkedIssuesView) UpdateData(report *WorkedIssuesReport) {
	v.data = report

	if report == nil || report.TotalCount == 0 {
		// Show empty state
		emptyLabel := widget.NewLabel("No worked issues found for the selected date range.")
		v.container = container.NewCenter(emptyLabel)
		return
	}

	// Create table with data
	v.table = widget.NewTable(
		func() (int, int) {
			return len(report.Issues) + 1, 4 // +1 for header row, 4 columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			// Header row
			if id.Row == 0 {
				switch id.Col {
				case 0:
					label.SetText("Issue Key")
				case 1:
					label.SetText("Summary")
				case 2:
					label.SetText("Status")
				case 3:
					label.SetText("Time Logged")
				}
				label.TextStyle = fyne.TextStyle{Bold: true}
				return
			}

			// Data rows
			issue := report.Issues[id.Row-1]
			label.TextStyle = fyne.TextStyle{Bold: false}

			switch id.Col {
			case 0:
				label.SetText(issue.Key)
			case 1:
				label.SetText(issue.Summary)
			case 2:
				label.SetText(issue.Status)
			case 3:
				label.SetText(formatDuration(issue.TimeLogged))
			}
		},
	)

	// Set column widths
	v.table.SetColumnWidth(0, 130)  // Issue Key
	v.table.SetColumnWidth(1, 550)  // Summary - much wider
	v.table.SetColumnWidth(2, 100)  // Status
	v.table.SetColumnWidth(3, 100)  // Time Logged

	// Create summary label
	summaryText := fmt.Sprintf("Total: %d issues, %s logged",
		report.TotalCount,
		formatDuration(report.TotalTimeLogged))
	summaryLabel := widget.NewLabel(summaryText)

	v.container = container.NewBorder(nil, summaryLabel, nil, nil, v.table)
}

// GetContainer returns the container for this view
func (v *WorkedIssuesView) GetContainer() *fyne.Container {
	return v.container
}

// UntrackedIssuesView creates a view for displaying untracked issues
type UntrackedIssuesView struct {
	container *fyne.Container
	table     *widget.Table
	data      *UntrackedIssuesReport
}

// NewUntrackedIssuesView creates a new untracked issues view
func NewUntrackedIssuesView() *UntrackedIssuesView {
	view := &UntrackedIssuesView{
		data: &UntrackedIssuesReport{
			Issues:     []IssueReportItem{},
			TotalCount: 0,
		},
	}

	// Create empty state message
	emptyLabel := widget.NewLabel("No untracked issues found. Click 'Refresh Reports' to load data.")
	view.container = container.NewCenter(emptyLabel)

	return view
}

// UpdateData updates the view with new report data
func (v *UntrackedIssuesView) UpdateData(report *UntrackedIssuesReport) {
	v.data = report

	if report == nil || report.TotalCount == 0 {
		// Show empty state
		emptyLabel := widget.NewLabel("All assigned issues have work logs. Great job!")
		v.container = container.NewCenter(emptyLabel)
		return
	}

	// Create table with data
	v.table = widget.NewTable(
		func() (int, int) {
			return len(report.Issues) + 1, 4 // +1 for header row, 4 columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			// Header row
			if id.Row == 0 {
				switch id.Col {
				case 0:
					label.SetText("Issue Key")
				case 1:
					label.SetText("Summary")
				case 2:
					label.SetText("Status")
				case 3:
					label.SetText("Assigned Date")
				}
				label.TextStyle = fyne.TextStyle{Bold: true}
				return
			}

			// Data rows
			issue := report.Issues[id.Row-1]
			label.TextStyle = fyne.TextStyle{Bold: false}

			switch id.Col {
			case 0:
				label.SetText(issue.Key)
			case 1:
				label.SetText(issue.Summary)
			case 2:
				label.SetText(issue.Status)
			case 3:
				label.SetText(issue.AssignedDate.Format("2006-01-02"))
			}
		},
	)

	// Set column widths
	v.table.SetColumnWidth(0, 130)  // Issue Key
	v.table.SetColumnWidth(1, 550)  // Summary - much wider
	v.table.SetColumnWidth(2, 100)  // Status
	v.table.SetColumnWidth(3, 100)  // Assigned Date

	v.container = container.NewBorder(nil, nil, nil, nil, v.table)
}

// GetContainer returns the container for this view
func (v *UntrackedIssuesView) GetContainer() *fyne.Container {
	return v.container
}

// RollupReportView creates a view for displaying rollup statistics
type RollupReportView struct {
	container *fyne.Container
	data      *RollupReport
}

// NewRollupReportView creates a new rollup report view
func NewRollupReportView() *RollupReportView {
	view := &RollupReportView{
		data: &RollupReport{
			TotalAssigned:     0,
			TotalWorked:       0,
			StatusBreakdown:   []StatusCount{},
			TotalUniqueIssues: 0,
		},
	}

	// Create empty state message
	emptyLabel := widget.NewLabel("No rollup data available. Click 'Refresh Reports' to load data.")
	view.container = container.NewCenter(emptyLabel)

	return view
}

// UpdateData updates the view with new report data
func (v *RollupReportView) UpdateData(report *RollupReport) {
	v.data = report

	if report == nil {
		// Show empty state
		emptyLabel := widget.NewLabel("No rollup data available for the selected date range.")
		v.container = container.NewCenter(emptyLabel)
		return
	}

	// Create summary statistics
	summaryText := fmt.Sprintf(
		"Total Assigned Issues: %d\nTotal Worked Issues: %d\nTotal Unique Issues: %d",
		report.TotalAssigned,
		report.TotalWorked,
		report.TotalUniqueIssues,
	)
	summaryLabel := widget.NewLabel(summaryText)
	summaryLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Create status breakdown table
	statusTable := widget.NewTable(
		func() (int, int) {
			return len(report.StatusBreakdown) + 1, 3 // +1 for header row, 3 columns
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			// Header row
			if id.Row == 0 {
				switch id.Col {
				case 0:
					label.SetText("Status")
				case 1:
					label.SetText("Count")
				case 2:
					label.SetText("Percentage")
				}
				label.TextStyle = fyne.TextStyle{Bold: true}
				return
			}

			// Data rows
			status := report.StatusBreakdown[id.Row-1]
			label.TextStyle = fyne.TextStyle{Bold: false}

			switch id.Col {
			case 0:
				label.SetText(status.StatusName)
			case 1:
				label.SetText(fmt.Sprintf("%d", status.Count))
			case 2:
				label.SetText(fmt.Sprintf("%.1f%%", status.Percentage))
			}
		},
	)

	// Set column widths
	statusTable.SetColumnWidth(0, 200) // Status
	statusTable.SetColumnWidth(1, 100) // Count
	statusTable.SetColumnWidth(2, 120) // Percentage

	// Create status breakdown section
	statusLabel := widget.NewLabel("Status Breakdown:")
	statusLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Create header section
	headerSection := container.NewVBox(
		summaryLabel,
		widget.NewSeparator(),
		statusLabel,
	)

	// Use Border layout to give table maximum available space
	v.container = container.NewBorder(
		headerSection, // top
		nil,          // bottom
		nil,          // left
		nil,          // right
		statusTable,  // center - takes remaining space
	)
}

// GetContainer returns the container for this view
func (v *RollupReportView) GetContainer() *fyne.Container {
	return v.container
}

// formatDuration formats a time.Duration into a human-readable string
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours == 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	if minutes == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
