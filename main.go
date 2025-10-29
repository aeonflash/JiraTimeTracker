package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

var jiraGraphQlBaseUri = "https://dish-enterprise.atlassian.net/gateway/api/graphql"
var jiraCloudId ="c724e043-0c6e-49f7-96a0-ba63500910c4"

var jiraApiKey string



func main() {
	// Check if running in CLI mode
	if len(os.Args) > 1 && os.Args[1] == "logs" {
		viewLogs()
		return
	}
	
	loadJiraConfig()
	a := app.New()
	w := a.NewWindow("JiraWidgetLite")
	w.SetTitle("JiraWidgetLite")

	startTime := time.Now()
	stopTime := time.Now()
	duration := time.Duration(0)

	ui := &UIComponents{
		StartTime:     &startTime,
		StopTime:      &stopTime,
		Duration:      &duration,
		StartText:     widget.NewLabel(""),
		StopText:      widget.NewLabel(""),
		DurationText:  widget.NewLabel(""),
		StatusLabel:   widget.NewLabel("Select an issue to start tracking time"),
		SelectedIssue: "",
		MainWindow:    w, // Pass window reference for dynamic resizing
	}

	content := createMainForm(ui)
	w.SetContent(content)
	
	// Set initial compact window size - wider to accommodate dropdown options
	w.Resize(fyne.NewSize(550, 200)) // Wider to fit full dropdown text
	w.SetFixedSize(false) // Allow resizing
	
	w.ShowAndRun()
}

func viewLogs() {
	entries, err := getTodaysTimeLog()
	if err != nil {
		fmt.Printf("Error reading time logs: %v\n", err)
		return
	}
	
	if len(entries) == 0 {
		fmt.Println("No time entries logged today.")
		return
	}
	
	fmt.Printf("Time entries for %s:\n", time.Now().Format("2006-01-02"))
	fmt.Println(strings.Repeat("=", 60))
	
	totalDuration := time.Duration(0)
	
	for i, entry := range entries {
		duration := entry.EndTime.Sub(entry.StartTime)
		totalDuration += duration
		
		fmt.Printf("%d. %s - %s\n", i+1, entry.JiraID, entry.Summary)
		fmt.Printf("   Time: %s (%s - %s)\n", 
			entry.Duration,
			entry.StartTime.Format("15:04"),
			entry.EndTime.Format("15:04"))
		if entry.Comment != "" && entry.Comment != "Time tracked via JiraTimeWidget" {
			fmt.Printf("   Comment: %s\n", entry.Comment)
		}
		fmt.Printf("   Logged: %s\n", entry.LoggedAt.Format("15:04:05"))
		fmt.Println()
	}
	
	fmt.Printf("Total time logged today: %s\n", formatDurationForJira(totalDuration))
}
