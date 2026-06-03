package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// jiraConfigExists checks whether ~/.jirarc exists and has content
func jiraConfigExists() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	info, err := os.Stat(homeDir + "/.jirarc")
	if err != nil {
		return false
	}
	return info.Size() > 0
}

// saveJiraConfig writes the config map to ~/.jirarc as JSON
func saveJiraConfig(config map[string]string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to determine home directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	filePath := homeDir + "/.jirarc"
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write %s: %w", filePath, err)
	}

	return nil
}

// showSetupWizard displays the configuration wizard and calls onComplete when finished
func showSetupWizard(app fyne.App, onComplete func()) {
	w := app.NewWindow("JiraWidgetLite - Setup")
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(500, 400))

	// Header
	title := widget.NewLabel("Welcome to JiraWidgetLite")
	title.TextStyle = fyne.TextStyle{Bold: true}
	subtitle := widget.NewLabel("Let's connect to your Jira instance.")

	// Form fields
	jiraURLEntry := widget.NewEntry()
	jiraURLEntry.SetPlaceHolder("https://your-org.atlassian.net")

	emailEntry := widget.NewEntry()
	emailEntry.SetPlaceHolder("you@company.com")

	apiTokenEntry := widget.NewPasswordEntry()
	apiTokenEntry.SetPlaceHolder("Your Jira API token")

	// Status label for validation messages
	statusLabel := widget.NewLabel("")

	// Help text
	helpText := widget.NewRichTextFromMarkdown(
		"Generate an API token at [id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens)",
	)

	// Save button
	saveButton := widget.NewButton("Connect", func() {
		// Validate inputs
		jiraURL := strings.TrimSpace(jiraURLEntry.Text)
		email := strings.TrimSpace(emailEntry.Text)
		apiToken := strings.TrimSpace(apiTokenEntry.Text)

		if jiraURL == "" {
			statusLabel.SetText("❌ Jira URL is required")
			return
		}
		if email == "" {
			statusLabel.SetText("❌ Email is required")
			return
		}
		if apiToken == "" {
			statusLabel.SetText("❌ API token is required")
			return
		}

		// Normalize URL: strip trailing slash, build graphqlUri
		jiraURL = strings.TrimRight(jiraURL, "/")
		graphqlUri := jiraURL + "/gateway/api/graphql"

		// Build config
		config := map[string]string{
			"jira":       apiToken,
			"email":      email,
			"graphqlUri": graphqlUri,
		}

		statusLabel.SetText("⏳ Saving configuration...")

		// Save to disk
		if err := saveJiraConfig(config); err != nil {
			statusLabel.SetText(fmt.Sprintf("❌ Failed to save: %v", err))
			log.Printf("Setup wizard save error: %v", err)
			return
		}

		// Load the config into the running app
		loadJiraConfig()

		statusLabel.SetText("✅ Configuration saved!")

		// Close wizard and launch main app
		w.Close()
		onComplete()
	})

	// Layout
	form := container.New(layout.NewFormLayout(),
		widget.NewLabel("Jira URL"), jiraURLEntry,
		widget.NewLabel("Email"), emailEntry,
		widget.NewLabel("API Token"), apiTokenEntry,
	)

	content := container.NewVBox(
		title,
		subtitle,
		widget.NewSeparator(),
		form,
		helpText,
		widget.NewSeparator(),
		statusLabel,
		container.NewHBox(layout.NewSpacer(), saveButton),
	)

	w.SetContent(container.NewPadded(content))
	w.SetOnClosed(func() {
		// If user closes wizard without completing, exit app
		if !jiraConfigExists() {
			os.Exit(0)
		}
	})
	w.Show()
}
