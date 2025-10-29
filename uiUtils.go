package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type UIComponents struct {
	StartTime            *time.Time
	StopTime             *time.Time
	Duration             *time.Duration
	StartText            *widget.Label
	StopText             *widget.Label
	DurationText         *widget.Label
	DurationEntry        *widget.Entry
	CommentEntry         *widget.Entry
	StatusLabel          *widget.Label
	RecentSelect         *widget.Select
	SelectedIssue        string // Store the selected issue key
	StartContainer       *fyne.Container
	EndContainer         *fyne.Container
	DurationContainer    *fyne.Container
	TimeButtonsContainer *fyne.Container
	CommentContainer     *fyne.Container
	LogButton            *widget.Button
	BrowserButton        *widget.Button
	MainWindow           fyne.Window
}

func createTimeButtons(ui *UIComponents) *fyne.Container {
	startButton := widget.NewButton("Start", func() {
		*ui.StartTime = time.Now()
		ui.StartText.SetText(ui.StartTime.Format("15:04:05"))
		log.Println(ui.StartTime.Format(time.RFC3339))
		ui.StatusLabel.SetText("⏱️ Time tracking started")
		
		// Show start time, hide end time
		ui.StartContainer.Show()
		ui.EndContainer.Hide()
		
		// Clear previous values
		*ui.Duration = 0
		ui.StopText.SetText("")
		
		// Update log button state
		updateLogButtonState(ui)
	})
	
	stopButton := widget.NewButton("Stop", func() {
		if ui.StartTime.IsZero() {
			ui.StatusLabel.SetText("❌ Please start timing first")
			return
		}
		
		*ui.StopTime = time.Now()
		ui.StopText.SetText(ui.StopTime.Format("15:04:05"))
		log.Println(ui.StopTime.Format(time.RFC3339))
		*ui.Duration = ui.StopTime.Sub(*ui.StartTime)
		
		// Format duration for entry field
		durationStr := formatDurationForJira(*ui.Duration)
		ui.DurationEntry.SetText(durationStr)
		
		// Show end time
		ui.EndContainer.Show()
		
		ui.StatusLabel.SetText("✅ Time tracking stopped")
		
		// Update log button state
		updateLogButtonState(ui)
	})
	
	resetButton := widget.NewButton("Reset", func() {
		ui.StartText.SetText("")
		ui.StopText.SetText("")
		ui.DurationEntry.SetText("")
		ui.StatusLabel.SetText("🔄 Timer reset")
		*ui.Duration = 0
		*ui.StartTime = time.Time{}
		*ui.StopTime = time.Time{}
		
		// Hide start and end containers, but keep duration and buttons visible if issue is selected
		ui.StartContainer.Hide()
		ui.EndContainer.Hide()
		
		// Update log button state
		updateLogButtonState(ui)
	})
	
	// Create browser button to open issue in browser
	ui.BrowserButton = widget.NewButton("🌐", func() {
		if ui.SelectedIssue != "" {
			// Construct Jira issue URL
			baseURL := strings.Replace(jiraGraphQlBaseUri, "/gateway/api/graphql", "", 1)
			issueURL := fmt.Sprintf("%s/browse/%s", baseURL, ui.SelectedIssue)
			
			if err := openBrowser(issueURL); err != nil {
				log.Printf("Error opening browser: %v", err)
				ui.StatusLabel.SetText("❌ Failed to open browser")
			} else {
				ui.StatusLabel.SetText(fmt.Sprintf("🌐 Opened %s in browser", ui.SelectedIssue))
			}
		}
	})
	ui.BrowserButton.Hide() // Initially hidden until issue is selected
	
	// Create the container and store it for show/hide control
	ui.TimeButtonsContainer = container.NewHBox(startButton, stopButton, resetButton, ui.BrowserButton)
	ui.TimeButtonsContainer.Hide() // Initially hidden until issue is selected
	
	return ui.TimeButtonsContainer
}

func updateLogButtonState(ui *UIComponents) {
	if ui.LogButton != nil {
		// Enable button if issue is selected and either:
		// 1. Timer has recorded time, OR
		// 2. Manual duration has been entered
		hasTimerDuration := ui.Duration.Seconds() > 0
		hasManualDuration := ui.DurationEntry != nil && ui.DurationEntry.Text != ""
		
		canLog := ui.SelectedIssue != "" && (hasTimerDuration || hasManualDuration)
		if canLog {
			ui.LogButton.Enable()
		} else {
			ui.LogButton.Disable()
		}
	}
}

// Parse manual duration input (e.g., "1h 30m", "45m", "2h")
func parseDuration(input string) time.Duration {
	if input == "" {
		return 0
	}
	
	// Try to parse as Go duration first
	if d, err := time.ParseDuration(input); err == nil {
		return d
	}
	
	// Try common formats like "1h 30m", "45m", "2h"
	input = strings.ReplaceAll(input, " ", "")
	if d, err := time.ParseDuration(input); err == nil {
		return d
	}
	
	return 0
}

// Open URL in default browser
func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

// Dynamically resize window based on visible content
func resizeWindowToContent(ui *UIComponents) {
	if ui.MainWindow == nil {
		return
	}
	
	// Calculate height based on visible components
	baseHeight := float32(120) // Issue selector + status
	
	if !ui.TimeButtonsContainer.Hidden {
		baseHeight += 50 // Time buttons
	}
	
	if !ui.DurationContainer.Hidden {
		baseHeight += 40 // Duration field
	}
	
	if !ui.CommentContainer.Hidden {
		baseHeight += 120 // Comment section
	}
	
	if !ui.LogButton.Hidden {
		baseHeight += 50 // Log button
	}
	
	// Add some padding
	baseHeight += 40
	
	// Ensure minimum width for dropdown options (wider to accommodate full text)
	minWidth := float32(550) // Increased from 450 to better fit dropdown options
	currentSize := ui.MainWindow.Content().Size()
	width := currentSize.Width
	if width < minWidth {
		width = minWidth
	}
	
	// Resize window
	newSize := fyne.NewSize(width, baseHeight)
	ui.MainWindow.Resize(newSize)
}

func createIssueSelector(ui *UIComponents) *fyne.Container {
	// Function to fetch and display issue
	fetchIssue := func(issueKey string) {
		if issueKey == "" {
			ui.StatusLabel.SetText("Select an issue to start tracking time")
			return
		}
		
		log.Println("Fetching Jira issue:", issueKey)
		ui.StatusLabel.SetText("🔍 Fetching issue...")
		
		jiraItem := getJiraItem(issueKey, jiraApiKey)
		if jiraItem == "" {
			ui.StatusLabel.SetText("❌ Failed to fetch issue")
			return
		}
		
		var issueData map[string]interface{}
		if err := json.Unmarshal([]byte(jiraItem), &issueData); err != nil {
			ui.StatusLabel.SetText("❌ Error parsing response")
			log.Println("JSON parse error:", err)
			return
		}
		
		// Check for error response
		if errorMsg, exists := issueData["errorMessage"]; exists {
			ui.StatusLabel.SetText(fmt.Sprintf("❌ API Error: %v", errorMsg))
			return
		}
		
		// Extract summary from fields (we don't display it, just validate the response)
		if fields, ok := issueData["fields"].(map[string]interface{}); ok {
			if _, ok := fields["summary"].(string); ok {
				ui.StatusLabel.SetText(fmt.Sprintf("✅ Ready to track time on %s", issueKey))
				ui.SelectedIssue = issueKey
				
				// Show time buttons, duration field, comment section and browser button when issue is selected
				ui.TimeButtonsContainer.Show()
				ui.DurationContainer.Show()
				ui.CommentContainer.Show()
				if ui.BrowserButton != nil {
					ui.BrowserButton.Show()
				}
				
				// Resize window to fit new content
				resizeWindowToContent(ui)
				
				// Update log button state
				updateLogButtonState(ui)
			} else {
				ui.StatusLabel.SetText("⚠️ Summary not found")
			}
		} else {
			ui.StatusLabel.SetText("❌ Unexpected data format")
		}
	}
	
	// Load recent issues for dropdown
	recentIssues := getRecentIssues(20)
	var recentOptions []string
	recentMap := make(map[string]string)
	
	if len(recentIssues) > 0 {
		recentOptions = append(recentOptions, "Select an issue to track time...")
		for _, issue := range recentIssues {
			option := fmt.Sprintf("%s [%s] - %s", issue.Key, issue.Status, issue.Summary)
			if len(option) > 100 {
				option = option[:97] + "..."
			}
			recentOptions = append(recentOptions, option)
			recentMap[option] = issue.Key
		}
	} else {
		recentOptions = append(recentOptions, "No recent issues found")
	}
	
	ui.RecentSelect = widget.NewSelect(recentOptions, func(selected string) {
		if issueKey, exists := recentMap[selected]; exists {
			fetchIssue(issueKey)
		}
	})
	
	// Create refresh button (icon only)
	refreshButton := widget.NewButton("🔄", func() {
		ui.StatusLabel.SetText("🔄 Refreshing issues...")
		
		// Reload recent issues
		newRecentIssues := getRecentIssues(20)
		var newOptions []string
		newMap := make(map[string]string)
		
		if len(newRecentIssues) > 0 {
			newOptions = append(newOptions, "Select an issue to track time...")
			for _, issue := range newRecentIssues {
				option := fmt.Sprintf("%s [%s] - %s", issue.Key, issue.Status, issue.Summary)
				if len(option) > 100 {
					option = option[:97] + "..."
				}
				newOptions = append(newOptions, option)
				newMap[option] = issue.Key
			}
		} else {
			newOptions = append(newOptions, "No recent issues found")
		}
		
		ui.RecentSelect.Options = newOptions
		ui.RecentSelect.OnChanged = func(selected string) {
			if issueKey, exists := newMap[selected]; exists {
				fetchIssue(issueKey)
			}
		}
		ui.RecentSelect.Refresh()
		ui.StatusLabel.SetText("✅ Issues refreshed")
	})
	
	// Create horizontal container for dropdown and refresh button
	selectorRow := container.NewBorder(nil, nil, nil, refreshButton, ui.RecentSelect)
	
	return container.NewVBox(selectorRow)
}

func createJiraItemContainer(ui *UIComponents) *fyne.Container {
	// Create editable duration entry
	ui.DurationEntry = widget.NewEntry()
	ui.DurationEntry.SetPlaceHolder("e.g., 1h 30m, 45m, 2h")
	ui.DurationEntry.Resize(fyne.NewSize(200, 32)) // Make it wider
	ui.DurationEntry.OnChanged = func(text string) {
		// Update log button state when duration is manually entered
		updateLogButtonState(ui)
	}
	
	// Time tracking section - will be updated dynamically with minimal spacing
	ui.StartContainer = container.NewHBox(widget.NewLabel("Start"), ui.StartText)
	ui.EndContainer = container.NewHBox(widget.NewLabel("End"), ui.StopText)
	
	// Create duration container with more space for the entry field
	durationLabel := widget.NewLabel("Duration")
	durationLabel.Resize(fyne.NewSize(80, 32)) // Fixed width for label
	ui.DurationContainer = container.NewBorder(nil, nil, durationLabel, nil, ui.DurationEntry)
	
	// Initially hide all time containers
	ui.StartContainer.Hide()
	ui.EndContainer.Hide()
	ui.DurationContainer.Hide()
	
	// Create a compact time tracking container with minimal spacing
	timeTrackingContainer := container.NewVBox(
		ui.StartContainer,
		ui.EndContainer,
		ui.DurationContainer,
	)
	
	// Work comment section - multiline entry
	commentLabel := widget.NewLabel("Work Comment")
	
	// Create multiline entry
	ui.CommentEntry = widget.NewMultiLineEntry()
	ui.CommentEntry.Wrapping = fyne.TextWrapWord
	ui.CommentEntry.SetMinRowsVisible(3)
	
	// Create comment container that can be hidden
	ui.CommentContainer = container.NewVBox(commentLabel, ui.CommentEntry)
	ui.CommentContainer.Hide() // Initially hidden
	
	// Create main container with controlled spacing (no summary section)
	return container.NewVBox(
		timeTrackingContainer,
		ui.CommentContainer,
		ui.StatusLabel,
	)
}

func createMainForm(ui *UIComponents) *fyne.Container {
	buttonContainer := createTimeButtons(ui)
	issueSelectorContainer := createIssueSelector(ui)
	jiraItemContainer := createJiraItemContainer(ui)
	
	// Create log button
	ui.LogButton = widget.NewButton("Log Time", func() {
		if ui.SelectedIssue == "" {
			ui.StatusLabel.SetText("❌ Please select an issue first")
			return
		}
		
		// Get duration from either timer or manual entry
		var finalDuration time.Duration
		var timeSpent string
		
		if ui.DurationEntry.Text != "" {
			// Use manual duration entry
			manualDuration := parseDuration(ui.DurationEntry.Text)
			if manualDuration <= 0 {
				ui.StatusLabel.SetText("❌ Invalid duration format")
				return
			}
			finalDuration = manualDuration
			timeSpent = ui.DurationEntry.Text
		} else if ui.Duration.Seconds() > 0 {
			// Use timer duration
			finalDuration = *ui.Duration
			timeSpent = formatDurationForJira(*ui.Duration)
		} else {
			ui.StatusLabel.SetText("❌ Please enter a duration")
			return
		}
		
		ui.StatusLabel.SetText("⏳ Logging work to Jira...")
		
		comment := ui.CommentEntry.Text
		if comment == "" {
			comment = "Time tracked via JiraTimeWidget"
		}
		
		// Use current time if no start time was recorded
		startTime := *ui.StartTime
		endTime := *ui.StopTime
		if startTime.IsZero() {
			startTime = time.Now().Add(-finalDuration)
			endTime = time.Now()
		}
		
		// Save to local log first
		logEntry := TimeLogEntry{
			JiraID:    ui.SelectedIssue,
			Summary:   "", // Summary removed from UI
			StartTime: startTime,
			EndTime:   endTime,
			Duration:  timeSpent,
			Comment:   comment,
			LoggedAt:  time.Now(),
		}
		
		if err := saveTimeLogEntry(logEntry); err != nil {
			log.Printf("Warning: Failed to save local log: %v", err)
		}
		
		// Log to Jira
		err := logWorkToJira(ui.SelectedIssue, timeSpent, comment, startTime)
		if err != nil {
			ui.StatusLabel.SetText(fmt.Sprintf("❌ Failed to log work: %v", err))
			log.Printf("Error logging work: %v", err)
		} else {
			ui.StatusLabel.SetText(fmt.Sprintf("✅ Logged %s to %s", timeSpent, ui.SelectedIssue))
			log.Printf("Successfully logged %s to %s", timeSpent, ui.SelectedIssue)
			
			// Reset the timer but keep the issue selected and duration field visible
			ui.StartText.SetText("")
			ui.StopText.SetText("")
			ui.DurationEntry.SetText("")
			ui.CommentEntry.SetText("")
			*ui.Duration = 0
			*ui.StartTime = time.Time{}
			*ui.StopTime = time.Time{}
			
			// Hide start/end containers but keep duration visible
			ui.StartContainer.Hide()
			ui.EndContainer.Hide()
			
			// Resize window after hiding containers
			resizeWindowToContent(ui)
			
			// Update button state
			updateLogButtonState(ui)
		}
	})
	
	// Initially disable the log button
	ui.LogButton.Disable()
	
	// Create bottom container with log button aligned right
	bottomContainer := container.NewBorder(nil, nil, nil, ui.LogButton)
	
	// Create main container with proper margins
	content := container.NewVBox(
		issueSelectorContainer,
		buttonContainer,
		jiraItemContainer,
		bottomContainer,
	)
	
	// Add margins - reduced padding (this will apply to all content including the button)
	return container.NewPadded(content)
}