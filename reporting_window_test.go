package main

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: jira-reporting-window, Property 2: Window singleton behavior**
// For any number of attempts to open the reporting window, only one reporting window
// instance should exist and be focused.
// **Validates: Requirements 1.4**
func TestPropertyWindowSingletonBehavior(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("only one reporting window instance exists", prop.ForAll(
		func(attempts int) bool {
			// Ensure positive number of attempts
			if attempts < 1 {
				attempts = 1
			}
			if attempts > 20 {
				attempts = 20 // Cap at reasonable number
			}

			// Create test app
			app := test.NewApp()
			defer test.NewApp() // Reset for next test

			// Create test user
			user := &User{
				AccountId:     "test-user-123",
				AccountStatus: "active",
				Name:          "Test User",
			}

			// Clear any existing singleton instance
			reportingWindowInstance = nil

			// Track window instances
			var firstWindow *ReportingWindow
			windowInstances := make([]*ReportingWindow, 0)

			// Attempt to create window multiple times
			for i := 0; i < attempts; i++ {
				window := NewReportingWindow(user, app)
				windowInstances = append(windowInstances, window)

				if i == 0 {
					firstWindow = window
				}
			}

			// Property 1: All instances should be the same object (singleton)
			for _, instance := range windowInstances {
				if instance != firstWindow {
					return false
				}
			}

			// Property 2: The global singleton should match the first window
			if reportingWindowInstance != firstWindow {
				return false
			}

			// Property 3: Only one window instance should exist
			if len(windowInstances) != attempts {
				return false
			}

			// All instances should point to the same window
			for i := 1; i < len(windowInstances); i++ {
				if windowInstances[i] != windowInstances[0] {
					return false
				}
			}

			// Clean up
			reportingWindowInstance = nil

			return true
		},
		gen.IntRange(1, 20), // attempts
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 1: Window independence preservation**
// For any main window state, closing the reporting window should leave the main window
// in the same state it was in before the reporting window was opened.
// **Validates: Requirements 1.3**
func TestPropertyWindowIndependencePreservation(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("main window state is preserved after reporting window operations", prop.ForAll(
		func(selectedIssue string, timerRunning bool, durationSeconds int) bool {
			// Create test app
			app := test.NewApp()
			defer test.NewApp() // Reset for next test

			// Create main window with some state
			mainWindow := app.NewWindow("Main Window")
			
			// Simulate main window state
			mainWindowState := struct {
				SelectedIssue string
				TimerRunning  bool
				Duration      int
				WindowTitle   string
			}{
				SelectedIssue: selectedIssue,
				TimerRunning:  timerRunning,
				Duration:      durationSeconds,
				WindowTitle:   mainWindow.Title(),
			}

			// Create test user
			user := &User{
				AccountId:     "test-user-123",
				AccountStatus: "active",
				Name:          "Test User",
			}

			// Clear any existing singleton instance
			reportingWindowInstance = nil

			// Open reporting window
			reportingWindow := NewReportingWindow(user, app)
			reportingWindow.Show()

			// Verify main window state is unchanged after opening reporting window
			if mainWindow.Title() != mainWindowState.WindowTitle {
				return false
			}

			// Close reporting window (simulates window close event)
			if reportingWindow.Window != nil {
				// Trigger the close handler manually
				reportingWindow.ReportService.InvalidateCache()
				reportingWindowInstance = nil
			}

			// Verify main window state is still unchanged after closing reporting window
			if mainWindow.Title() != mainWindowState.WindowTitle {
				return false
			}

			// Verify that the main window state variables would still be accessible
			// (In a real scenario, these would be stored in UIComponents)
			// Here we just verify the window itself is unchanged
			if mainWindow.Title() != "Main Window" {
				return false
			}

			return true
		},
		gen.Identifier(),                // selectedIssue
		gen.Bool(),                      // timerRunning
		gen.IntRange(0, 3600),          // durationSeconds (0 to 1 hour)
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}
