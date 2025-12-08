package main

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: jira-reporting-window, Property 32: Loading indicator visibility during generation**
// For any report generation operation in progress, a loading indicator should be visible
// in the reporting window.
// **Validates: Requirements 9.1**
func TestPropertyLoadingIndicatorVisibility(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("loading indicator is visible during report generation", prop.ForAll(
		func(message string) bool {
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

			// Create reporting window
			rw := NewReportingWindow(user, app)

			// Initially, loading indicator should be hidden
			if rw.LoadingIndicator.Visible() {
				return false
			}

			// Call ShowLoading
			rw.ShowLoading(message)

			// Loading indicator should now be visible
			if !rw.LoadingIndicator.Visible() {
				return false
			}

			// Status label should show the message
			if rw.StatusLabel.Text != message {
				return false
			}

			// Call HideLoading
			rw.HideLoading()

			// Loading indicator should now be hidden
			if rw.LoadingIndicator.Visible() {
				return false
			}

			// Clean up
			reportingWindowInstance = nil

			return true
		},
		gen.AnyString(), // message
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 33: Control disabling during loading**
// For any data retrieval operation in progress, report controls should be disabled
// to prevent duplicate requests.
// **Validates: Requirements 9.2**
func TestPropertyControlDisablingDuringLoading(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("controls are disabled during loading", prop.ForAll(
		func(message string) bool {
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

			// Create reporting window
			rw := NewReportingWindow(user, app)

			// Initially, controls should be enabled
			if rw.RefreshButton.Disabled() {
				return false
			}
			if rw.ExportButton.Disabled() {
				// Export button is initially disabled until reports are generated
				// This is expected behavior
			}
			if rw.StartDatePicker.Disabled() {
				return false
			}
			if rw.EndDatePicker.Disabled() {
				return false
			}

			// Call ShowLoading
			rw.ShowLoading(message)

			// All controls should now be disabled
			if !rw.RefreshButton.Disabled() {
				return false
			}
			if !rw.ExportButton.Disabled() {
				return false
			}
			if !rw.StartDatePicker.Disabled() {
				return false
			}
			if !rw.EndDatePicker.Disabled() {
				return false
			}

			// Call HideLoading
			rw.HideLoading()

			// Controls should now be enabled again
			if rw.RefreshButton.Disabled() {
				return false
			}
			if rw.ExportButton.Disabled() {
				// Export button remains disabled until reports are generated
				// This is expected behavior
			}
			if rw.StartDatePicker.Disabled() {
				return false
			}
			if rw.EndDatePicker.Disabled() {
				return false
			}

			// Clean up
			reportingWindowInstance = nil

			return true
		},
		gen.AnyString(), // message
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 34: Success state transition**
// For any successful data retrieval, the loading indicator should be hidden and
// the results should be displayed.
// **Validates: Requirements 9.3**
func TestPropertySuccessStateTransition(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("success state transition hides loading and shows success", prop.ForAll(
		func(message string) bool {
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

			// Create reporting window
			rw := NewReportingWindow(user, app)

			// Start loading state
			rw.ShowLoading("Loading...")

			// Verify loading state
			if !rw.LoadingIndicator.Visible() {
				return false
			}
			if !rw.RefreshButton.Disabled() {
				return false
			}

			// Transition to success state
			rw.ShowSuccess(message)

			// Verify success state
			// Loading indicator should be hidden
			if rw.LoadingIndicator.Visible() {
				return false
			}

			// Controls should be enabled
			if rw.RefreshButton.Disabled() {
				return false
			}
			if rw.StartDatePicker.Disabled() {
				return false
			}
			if rw.EndDatePicker.Disabled() {
				return false
			}

			// Status label should show success message with checkmark
			expectedMessage := "✅ " + message
			if rw.StatusLabel.Text != expectedMessage {
				return false
			}

			// Clean up
			reportingWindowInstance = nil

			return true
		},
		gen.AnyString(), // message
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: jira-reporting-window, Property 35: Failure state transition**
// For any failed data retrieval, the loading indicator should be hidden and
// an error message should be displayed.
// **Validates: Requirements 9.4**
func TestPropertyFailureStateTransition(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("failure state transition hides loading and shows error", prop.ForAll(
		func(message string) bool {
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

			// Create reporting window
			rw := NewReportingWindow(user, app)

			// Start loading state
			rw.ShowLoading("Loading...")

			// Verify loading state
			if !rw.LoadingIndicator.Visible() {
				return false
			}
			if !rw.RefreshButton.Disabled() {
				return false
			}

			// Transition to failure state
			rw.ShowError(message)

			// Verify failure state
			// Loading indicator should be hidden
			if rw.LoadingIndicator.Visible() {
				return false
			}

			// Controls should be enabled
			if rw.RefreshButton.Disabled() {
				return false
			}
			if rw.StartDatePicker.Disabled() {
				return false
			}
			if rw.EndDatePicker.Disabled() {
				return false
			}

			// Status label should show error message with X mark
			expectedMessage := "❌ Error: " + message
			if rw.StatusLabel.Text != expectedMessage {
				return false
			}

			// Clean up
			reportingWindowInstance = nil

			return true
		},
		gen.AnyString(), // message
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}
