package main

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

// Global variable to track the singleton reporting window instance
var reportingWindowInstance *ReportingWindow

// ReportingWindow manages the UI for the reporting window
type ReportingWindow struct {
	Window              fyne.Window
	StartDateButton     *widget.Button
	EndDateButton       *widget.Button
	ReportTabs          *container.AppTabs
	LoadingIndicator    *widget.ProgressBarInfinite
	StatusLabel         *widget.Label
	ExportButton        *widget.Button
	RefreshButton       *widget.Button
	ReportService       *ReportService
	CurrentUser         *User
	CurrentDateRange    DateRange
	AssignedIssuesView  *AssignedIssuesView
	WorkedIssuesView    *WorkedIssuesView
	UntrackedIssuesView *UntrackedIssuesView
	RollupReportView    *RollupReportView
	RequestQueue        *RequestQueue
	NextRequestID       int
}

// NewReportingWindow creates and initializes the reporting window
// Implements singleton pattern - returns existing instance if already open
func NewReportingWindow(user *User, app fyne.App) *ReportingWindow {
	// Check if window already exists
	if reportingWindowInstance != nil && reportingWindowInstance.Window != nil {
		// Bring existing window to focus
		reportingWindowInstance.Window.RequestFocus()
		reportingWindowInstance.Window.Show()
		return reportingWindowInstance
	}

	// Create new window instance
	window := app.NewWindow("Jira Reports")

	// Initialize report service
	reportService := NewReportService(user)

	// Get default date range (current month)
	defaultDateRange := GetDefaultDateRange()

	// Create UI components
	rw := &ReportingWindow{
		Window:           window,
		ReportService:    reportService,
		CurrentUser:      user,
		CurrentDateRange: defaultDateRange,
		RequestQueue:     NewRequestQueue(),
		NextRequestID:    1,
	}

	// Initialize UI components
	rw.initializeComponents()

	// Build the UI layout
	content := rw.buildLayout()
	window.SetContent(content)

	// Set window size
	window.Resize(fyne.NewSize(900, 600))

	// Set up window close handler to clear singleton and cache
	window.SetOnClosed(func() {
		rw.ReportService.InvalidateCache()
		reportingWindowInstance = nil
	})

	// Store as singleton instance
	reportingWindowInstance = rw

	return rw
}

// initializeComponents creates and initializes all UI components
func (rw *ReportingWindow) initializeComponents() {
	// Date range picker buttons
	rw.StartDateButton = widget.NewButton(rw.CurrentDateRange.StartDate.Format("2006-01-02"), func() {
		rw.showStartDatePicker()
	})

	rw.EndDateButton = widget.NewButton(rw.CurrentDateRange.EndDate.Format("2006-01-02"), func() {
		rw.showEndDatePicker()
	})

	// Status label for messages
	rw.StatusLabel = widget.NewLabel("Select a date range and click Refresh to generate reports")

	// Loading indicator
	rw.LoadingIndicator = widget.NewProgressBarInfinite()
	rw.LoadingIndicator.Hide()

	// Refresh button
	rw.RefreshButton = widget.NewButton("Refresh Reports", func() {
		rw.refreshReports()
	})

	// Export button
	rw.ExportButton = widget.NewButton("Export", func() {
		rw.showExportDialog()
	})
	rw.ExportButton.Disable() // Initially disabled until reports are generated

	// Create report views
	rw.AssignedIssuesView = NewAssignedIssuesView()
	rw.WorkedIssuesView = NewWorkedIssuesView()
	rw.UntrackedIssuesView = NewUntrackedIssuesView()
	rw.RollupReportView = NewRollupReportView()

	// Create report view tabs
	rw.ReportTabs = container.NewAppTabs(
		container.NewTabItem("Assigned Issues", rw.AssignedIssuesView.GetContainer()),
		container.NewTabItem("Worked Issues", rw.WorkedIssuesView.GetContainer()),
		container.NewTabItem("Untracked Issues", rw.UntrackedIssuesView.GetContainer()),
		container.NewTabItem("Rollup Report", rw.RollupReportView.GetContainer()),
	)
}

// buildLayout constructs the main window layout
func (rw *ReportingWindow) buildLayout() *fyne.Container {
	// Date range controls with date picker buttons
	toLabel := widget.NewLabel("to")
	
	dateInputs := container.NewHBox(
		rw.StartDateButton,
		toLabel,
		rw.EndDateButton,
	)

	// Top row: date inputs on left, refresh button on far right
	topControls := container.NewBorder(
		nil, nil, dateInputs, rw.RefreshButton,
		container.NewHBox(), // Empty spacer
	)

	// Second row: status label on left, loading indicator and export button on right
	rightButtons := container.NewHBox(
		rw.LoadingIndicator,
		rw.ExportButton,
	)
	
	statusBar := container.NewBorder(
		nil, nil, nil, rightButtons,
		rw.StatusLabel,
	)

	// Main layout
	mainLayout := container.NewBorder(
		container.NewVBox(topControls, statusBar),
		nil, nil, nil,
		rw.ReportTabs,
	)

	return mainLayout
}

// updateReportViews updates all report views with new data
func (rw *ReportingWindow) updateReportViews(
	assignedReport *AssignedIssuesReport,
	workedReport *WorkedIssuesReport,
	untrackedReport *UntrackedIssuesReport,
	rollupReport *RollupReport,
) {
	// Update each view
	rw.AssignedIssuesView.UpdateData(assignedReport)
	rw.WorkedIssuesView.UpdateData(workedReport)
	rw.UntrackedIssuesView.UpdateData(untrackedReport)
	rw.RollupReportView.UpdateData(rollupReport)

	// Update tabs with new containers
	rw.ReportTabs.Items[0].Content = rw.AssignedIssuesView.GetContainer()
	rw.ReportTabs.Items[1].Content = rw.WorkedIssuesView.GetContainer()
	rw.ReportTabs.Items[2].Content = rw.UntrackedIssuesView.GetContainer()
	rw.ReportTabs.Items[3].Content = rw.RollupReportView.GetContainer()

	// Refresh the tabs
	rw.ReportTabs.Refresh()
}

// Show displays the reporting window
func (rw *ReportingWindow) Show() {
	rw.Window.Show()
}

// ShowLoading displays the loading indicator with a message
func (rw *ReportingWindow) ShowLoading(message string) {
	rw.StatusLabel.SetText(message)
	rw.LoadingIndicator.Show()
	rw.LoadingIndicator.Start()
	
	// Disable controls during loading
	rw.RefreshButton.Disable()
	rw.ExportButton.Disable()
	rw.StartDateButton.Disable()
	rw.EndDateButton.Disable()
}

// HideLoading hides the loading indicator
func (rw *ReportingWindow) HideLoading() {
	rw.LoadingIndicator.Stop()
	rw.LoadingIndicator.Hide()
	
	// Re-enable controls
	rw.RefreshButton.Enable()
	rw.ExportButton.Enable()
	rw.StartDateButton.Enable()
	rw.EndDateButton.Enable()
}

// ShowError displays an error message to the user
func (rw *ReportingWindow) ShowError(message string) {
	rw.HideLoading()
	rw.StatusLabel.SetText(fmt.Sprintf("❌ Error: %s", message))
}

// ShowSuccess displays a success message to the user
func (rw *ReportingWindow) ShowSuccess(message string) {
	rw.HideLoading()
	rw.StatusLabel.SetText(fmt.Sprintf("✅ %s", message))
}

// ShowWarning displays a warning message to the user
func (rw *ReportingWindow) ShowWarning(message string) {
	rw.HideLoading()
	rw.StatusLabel.SetText(fmt.Sprintf("⚠️ Warning: %s", message))
}

// showStartDatePicker displays a calendar popup for selecting the start date
func (rw *ReportingWindow) showStartDatePicker() {
	// Create popup first so we can reference it in the callback
	var popup *widget.PopUp
	
	calendar := xwidget.NewCalendar(rw.CurrentDateRange.StartDate, func(selectedDate time.Time) {
		// Update start date
		rw.CurrentDateRange.StartDate = selectedDate
		rw.StartDateButton.SetText(selectedDate.Format("2006-01-02"))
		
		// Validate and update
		if err := rw.CurrentDateRange.Validate(); err != nil {
			rw.StatusLabel.SetText(fmt.Sprintf("❌ Invalid date range: %s", err.Error()))
		} else {
			rw.ReportService.InvalidateCache()
			rw.StatusLabel.SetText("Date range updated. Click 'Refresh Reports' to reload data.")
		}
		
		// Hide the popup after selection
		if popup != nil {
			popup.Hide()
		}
	})
	
	// Create popup dialog
	popup = widget.NewModalPopUp(calendar, rw.Window.Canvas())
	popup.Show()
}

// showEndDatePicker displays a calendar popup for selecting the end date
func (rw *ReportingWindow) showEndDatePicker() {
	// Create popup first so we can reference it in the callback
	var popup *widget.PopUp
	
	calendar := xwidget.NewCalendar(rw.CurrentDateRange.EndDate, func(selectedDate time.Time) {
		// Update end date
		rw.CurrentDateRange.EndDate = selectedDate
		rw.EndDateButton.SetText(selectedDate.Format("2006-01-02"))
		
		// Validate and update
		if err := rw.CurrentDateRange.Validate(); err != nil {
			rw.StatusLabel.SetText(fmt.Sprintf("❌ Invalid date range: %s", err.Error()))
		} else {
			rw.ReportService.InvalidateCache()
			rw.StatusLabel.SetText("Date range updated. Click 'Refresh Reports' to reload data.")
		}
		
		// Hide the popup after selection
		if popup != nil {
			popup.Hide()
		}
	})
	
	// Create popup dialog
	popup = widget.NewModalPopUp(calendar, rw.Window.Canvas())
	popup.Show()
}

// refreshReports fetches and displays report data
func (rw *ReportingWindow) refreshReports() {
	// Validate current date range
	if err := rw.CurrentDateRange.Validate(); err != nil {
		reportErr := CategorizeError(err)
		rw.ShowError(reportErr.GetUserFriendlyMessage())
		return
	}

	// Create a new request
	request := ReportRequest{
		ID:        rw.NextRequestID,
		DateRange: rw.CurrentDateRange,
		Timestamp: time.Now(),
	}
	rw.NextRequestID++

	// Enqueue the request
	rw.RequestQueue.Enqueue(request)

	// Process the queue if not already processing
	if !rw.RequestQueue.IsProcessing() {
		go rw.processRequestQueue()
	}
}

// processRequestQueue processes requests from the queue sequentially
func (rw *ReportingWindow) processRequestQueue() {
	// Mark as processing
	rw.RequestQueue.SetProcessing(true)
	defer rw.RequestQueue.SetProcessing(false)

	// Process all requests in the queue
	for !rw.RequestQueue.IsEmpty() {
		request := rw.RequestQueue.Dequeue()
		if request == nil {
			break
		}

		// Process this request
		rw.processReportRequest(*request)
	}
}

// processReportRequest processes a single report request
// Note: This is called from processRequestQueue which runs in a goroutine
func (rw *ReportingWindow) processReportRequest(request ReportRequest) {
	// Show loading on UI thread
	fyne.Do(func() {
		rw.ShowLoading("Fetching report data from Jira...")
	})

	// Track partial data errors
	var partialErrors []string
	hasPartialData := false

	// Fetch all reports with error tracking
	assignedReport, err := rw.ReportService.GetAssignedIssues(request.DateRange)
	if err != nil {
		reportErr := CategorizeError(err)
		// For critical errors, stop and show error
		if reportErr.Type == ErrorTypeAuth || reportErr.Type == ErrorTypeNetwork || reportErr.Type == ErrorTypeRateLimit {
			fyne.Do(func() {
				rw.ShowError(reportErr.GetUserFriendlyMessage())
			})
			return
		}
		// For other errors, track as partial data
		hasPartialData = true
		partialErrors = append(partialErrors, "assigned issues")
		assignedReport = &AssignedIssuesReport{Issues: []IssueReportItem{}, TotalCount: 0}
	}

	workedReport, err := rw.ReportService.GetWorkedIssues(request.DateRange)
	if err != nil {
		reportErr := CategorizeError(err)
		// For critical errors, stop and show error
		if reportErr.Type == ErrorTypeAuth || reportErr.Type == ErrorTypeNetwork || reportErr.Type == ErrorTypeRateLimit {
			fyne.Do(func() {
				rw.ShowError(reportErr.GetUserFriendlyMessage())
			})
			return
		}
		// For other errors, track as partial data
		hasPartialData = true
		partialErrors = append(partialErrors, "worked issues")
		workedReport = &WorkedIssuesReport{Issues: []WorkedIssueItem{}, TotalCount: 0, TotalTimeLogged: 0}
	}

	untrackedReport, err := rw.ReportService.GetUntrackedIssues(request.DateRange)
	if err != nil {
		reportErr := CategorizeError(err)
		// For critical errors, stop and show error
		if reportErr.Type == ErrorTypeAuth || reportErr.Type == ErrorTypeNetwork || reportErr.Type == ErrorTypeRateLimit {
			fyne.Do(func() {
				rw.ShowError(reportErr.GetUserFriendlyMessage())
			})
			return
		}
		// For other errors, track as partial data
		hasPartialData = true
		partialErrors = append(partialErrors, "untracked issues")
		untrackedReport = &UntrackedIssuesReport{Issues: []IssueReportItem{}, TotalCount: 0}
	}

	rollupReport, err := rw.ReportService.GetRollupReport(request.DateRange)
	if err != nil {
		reportErr := CategorizeError(err)
		// For critical errors, stop and show error
		if reportErr.Type == ErrorTypeAuth || reportErr.Type == ErrorTypeNetwork || reportErr.Type == ErrorTypeRateLimit {
			fyne.Do(func() {
				rw.ShowError(reportErr.GetUserFriendlyMessage())
			})
			return
		}
		// For other errors, track as partial data
		hasPartialData = true
		partialErrors = append(partialErrors, "rollup report")
		rollupReport = &RollupReport{TotalAssigned: 0, TotalWorked: 0, StatusBreakdown: []StatusCount{}, TotalUniqueIssues: 0}
	}

	// Update views on the UI thread
	fyne.Do(func() {
		rw.updateReportViews(assignedReport, workedReport, untrackedReport, rollupReport)

		// Show appropriate message based on success/partial data
		if hasPartialData {
			errorList := strings.Join(partialErrors, ", ")
			rw.ShowWarning(fmt.Sprintf("Some data could not be retrieved (%s). Showing partial results.", errorList))
		} else {
			rw.ShowSuccess("Reports loaded successfully")
		}
		rw.ExportButton.Enable()
	})
}

// showExportDialog displays the export format selection dialog
func (rw *ReportingWindow) showExportDialog() {
	// Check if we have data to export
	if rw.AssignedIssuesView.data == nil || rw.WorkedIssuesView.data == nil ||
		rw.UntrackedIssuesView.data == nil || rw.RollupReportView.data == nil {
		rw.ShowError("No report data available to export. Please generate reports first.")
		return
	}

	// Create format selection dialog
	csvButton := widget.NewButton("Export as CSV", func() {
		rw.exportWithFormatSelection(ExportFormatCSV)
	})

	pdfButton := widget.NewButton("Export as PDF", func() {
		rw.exportWithFormatSelection(ExportFormatPDF)
	})

	cancelButton := widget.NewButton("Cancel", func() {
		// Dialog will be closed automatically
	})

	// Create dialog content
	content := container.NewVBox(
		widget.NewLabel("Select export format:"),
		csvButton,
		pdfButton,
		widget.NewSeparator(),
		cancelButton,
	)

	// Create and show dialog
	dialog := widget.NewModalPopUp(content, rw.Window.Canvas())
	
	// Update button callbacks to close dialog
	csvButton.OnTapped = func() {
		dialog.Hide()
		rw.exportWithFormatSelection(ExportFormatCSV)
	}
	
	pdfButton.OnTapped = func() {
		dialog.Hide()
		rw.exportWithFormatSelection(ExportFormatPDF)
	}
	
	cancelButton.OnTapped = func() {
		dialog.Hide()
	}

	dialog.Show()
}

// OpenReportingWindow is a helper function to open the reporting window from the main application
func OpenReportingWindow(user *User, app fyne.App) {
	rw := NewReportingWindow(user, app)
	rw.Show()
}

// exportWithFormatSelection exports the current report data with the selected format
func (rw *ReportingWindow) exportWithFormatSelection(format ExportFormat) {
	// Get current report data from views
	assignedReport := rw.AssignedIssuesView.data
	workedReport := rw.WorkedIssuesView.data
	untrackedReport := rw.UntrackedIssuesView.data
	rollupReport := rw.RollupReportView.data

	// Create export service
	exportService := NewExportService(rw.CurrentUser, rw.CurrentDateRange)

	// Determine file extension and default filename
	var extension string
	var formatName string
	if format == ExportFormatCSV {
		extension = ".csv"
		formatName = "CSV"
	} else {
		extension = ".pdf"
		formatName = "PDF"
	}

	// Generate default filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	defaultFilename := fmt.Sprintf("jira_report_%s%s", timestamp, extension)

	// Show file save dialog
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			rw.ShowError(fmt.Sprintf("Failed to select file location: %s", err.Error()))
			return
		}
		if writer == nil {
			// User cancelled the dialog
			return
		}
		defer writer.Close()

		// Get the selected file path
		filepath := writer.URI().Path()

		// Show loading indicator
		rw.ShowLoading(fmt.Sprintf("Exporting report to %s...", formatName))

		// Export in a goroutine to keep UI responsive
		go func() {
			var err error

			if format == ExportFormatCSV {
				err = exportService.ExportToCSVFile(assignedReport, workedReport, untrackedReport, rollupReport, filepath)
			} else {
				err = exportService.ExportToPDFFile(assignedReport, workedReport, untrackedReport, rollupReport, filepath)
			}

			if err != nil {
				rw.ShowError(fmt.Sprintf("Failed to export report: %s", err.Error()))
				return
			}

			// Show success message with file location
			rw.ShowSuccess(fmt.Sprintf("Report exported successfully to: %s", filepath))
		}()
	}, rw.Window)

	// Set default filename
	saveDialog.SetFileName(defaultFilename)

	// Show the dialog
	saveDialog.Show()
}
