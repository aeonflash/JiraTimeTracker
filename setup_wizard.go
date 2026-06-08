package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
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

// validateJiraURL checks that the URL is a valid HTTPS Atlassian URL
func validateJiraURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("Jira URL is required")
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}

	if parsed.Scheme != "https" {
		return fmt.Errorf("URL must start with https://")
	}

	if parsed.Host == "" {
		return fmt.Errorf("URL must include a hostname")
	}

	return nil
}

// validateEmail checks for a basic valid email format
func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Basic email regex: something@something.something
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// validateAPIToken checks that the token is non-empty and reasonable length
func validateAPIToken(token string) error {
	if token == "" {
		return fmt.Errorf("API token is required")
	}

	if len(token) < 10 {
		return fmt.Errorf("API token seems too short — check that you pasted the full token")
	}

	return nil
}

// testJiraConnection makes a test API call to verify credentials work
func testJiraConnection(baseURL, email, apiToken string) error {
	testURL := baseURL + "/rest/api/3/myself"

	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(email + ":" + apiToken))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed — check your URL and internet connection")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case 401:
		return fmt.Errorf("authentication failed — check your email and API token")
	case 403:
		return fmt.Errorf("access forbidden — your token may lack permissions")
	case 404:
		return fmt.Errorf("Jira API not found — check your URL is correct")
	default:
		return fmt.Errorf("unexpected response (%d): %s", resp.StatusCode, string(body[:min(200, len(body))]))
	}
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

	// Save button (declared first so it can be referenced in the callback)
	var saveButton *widget.Button
	saveButton = widget.NewButton("Connect", func() {
		// Validate inputs
		jiraURL := strings.TrimSpace(jiraURLEntry.Text)
		email := strings.TrimSpace(emailEntry.Text)
		apiToken := strings.TrimSpace(apiTokenEntry.Text)

		// URL validation
		if err := validateJiraURL(jiraURL); err != nil {
			statusLabel.SetText(fmt.Sprintf("❌ %s", err.Error()))
			return
		}

		// Email validation
		if err := validateEmail(email); err != nil {
			statusLabel.SetText(fmt.Sprintf("❌ %s", err.Error()))
			return
		}

		// API token validation
		if err := validateAPIToken(apiToken); err != nil {
			statusLabel.SetText(fmt.Sprintf("❌ %s", err.Error()))
			return
		}

		// Normalize URL: strip trailing slash
		jiraURL = strings.TrimRight(jiraURL, "/")
		graphqlUri := jiraURL + "/gateway/api/graphql"

		// Test connection
		statusLabel.SetText("⏳ Testing connection to Jira...")
		saveButton.Disable()

		go func() {
			if err := testJiraConnection(jiraURL, email, apiToken); err != nil {
				fyne.Do(func() {
					statusLabel.SetText(fmt.Sprintf("❌ %s", err.Error()))
					saveButton.Enable()
				})
				return
			}

			// Build config
			config := map[string]string{
				"jira":       apiToken,
				"email":      email,
				"graphqlUri": graphqlUri,
			}

			// Save to disk
			if err := saveJiraConfig(config); err != nil {
				fyne.Do(func() {
					statusLabel.SetText(fmt.Sprintf("❌ Failed to save: %v", err))
					saveButton.Enable()
				})
				log.Printf("Setup wizard save error: %v", err)
				return
			}

			// Load the config into the running app
			loadJiraConfig()

			fyne.Do(func() {
				statusLabel.SetText("✅ Connected successfully!")
				// Close wizard and launch main app
				w.Close()
				onComplete()
			})
		}()
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
